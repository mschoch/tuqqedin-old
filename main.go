//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/couchbaselabs/go-couchbase"
	"github.com/couchbaselabs/tuqqedin/ast"
	"github.com/couchbaselabs/tuqqedin/datasource"
	"github.com/couchbaselabs/tuqqedin/parser"
	"github.com/couchbaselabs/tuqqedin/plan"
	"github.com/gorilla/mux"
)

const STATIC_PREFIX = "/_static/"

var addr = flag.String("addr", ":8093", "HTTP listen address")
var cbServer = flag.String("couchbase", "http://localhost:8091/", "URL to couchbase")
var debugParsing = flag.Bool("debugParsing", false, "output parsing debug information")
var staticPath = flag.String("static-path", "static", "path to static web UI content")

var dataSourceManager datasource.DataSourceManager
var planner plan.Planner
var optimizer Optimizer
var executor Executor
var unqlParser parser.Parser

func main() {

	flag.Parse()

	client, err := couchbase.Connect(*cbServer)
	if err != nil {
		log.Fatalf("Error connecting to couchbase: %v", err)
	}

	dataSourceManager, err = datasource.NewCouchbaseDataSourceManager(client)
	unqlParser = parser.NewUnqlParser()
	planner = plan.NewCouchbasePlanner(dataSourceManager)
	optimizer = NewCouchbaseOptimizer()
	executor = NewCouchbaseExecutor()

	if *debugParsing {
		parser.DebugTokens = true
		parser.DebugGrammar = true
	}

	r := mux.NewRouter()

	// static
	r.PathPrefix(STATIC_PREFIX).Handler(
		http.StripPrefix(STATIC_PREFIX,
			http.FileServer(http.Dir(*staticPath))))

	// api
	r.HandleFunc("/api", welcome).Methods("GET")
	r.Handle("/api/{bucket}/_query_ast", http.HandlerFunc(bucketQueryAST)).Methods("POST")
	r.Handle("/api/{bucket}/_query", http.HandlerFunc(bucketQuery)).Methods("GET", "POST")
	r.Handle("/", http.RedirectHandler("/_static/index.html", 302))
	log.Printf("listening rest on: %v", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))
}

func mustEncode(w io.Writer, i interface{}) {
	if headered, ok := w.(http.ResponseWriter); ok {
		headered.Header().Set("Cache-Control", "no-cache")
		headered.Header().Set("Content-type", "application/json")
	}

	e := json.NewEncoder(w)
	if err := e.Encode(i); err != nil {
		panic(err)
	}
}

func showError(w http.ResponseWriter, r *http.Request,
	msg string, code int) {
	log.Printf("Reporting error %v/%v", code, msg)
	http.Error(w, msg, code)
}

func bucketQuery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	_, err := dataSourceManager.GetDataSource(bucket)
	if err != nil {
		showError(w, r, fmt.Sprintf("%v does not exist", bucket), 404)
		return
	}

	queryString := r.FormValue("q")
	if queryString == "" && r.Method == "POST" {
		queryStringBytes, err := ioutil.ReadAll(r.Body)
		if err == nil {
			queryString = string(queryStringBytes)
		}
	}

	if queryString == "" {
		showError(w, r, "Missing required query string", 500)
		return
	} else {
		log.Printf("Query String: %v", queryString)
	}

	statement, err := unqlParser.Parse(queryString)
	if err != nil {
		showError(w, r, err.Error(), 500)
		return
	}

	// add the from
	statement.SetFrom([]ast.DataSource{ast.NewNamedDataSource(bucket)})

	doExecuteStatement(w, r, statement)
}

func bucketQueryAST(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	bucket := vars["bucket"]
	_, err := dataSourceManager.GetDataSource(bucket)
	if err != nil {
		showError(w, r, fmt.Sprintf("%v does not exist", bucket), 404)
		return
	}

	d := json.NewDecoder(r.Body)
	request := map[string]interface{}{}

	err = d.Decode(&request)
	if err != nil {
		showError(w, r, "Error parsing request JSON", 500)
		return
	}

	statement, err := ast.NewStatementFromJSONRequestToBucket(bucket, request)
	if err != nil {
		showError(w, r, err.Error(), 500)
		return
	}

	doExecuteStatement(w, r, statement)
}

func doExecuteStatement(w http.ResponseWriter, r *http.Request, s ast.Statement) {
	log.Printf("Built statement %v", s)

	plans, err := planner.Plan(s)

	log.Printf("Plans for statement:")
	for i, plan := range plans {
		log.Printf("Plan %d:", i)
		log.Printf("%v", plan)
	}

	if err != nil {
		log.Printf("error plannging statement")
	}

	if len(plans) > 0 {
		optimalPlan := optimizer.ChooseOptimalPlan(plans)
		executor.ExecutePlan(w, optimalPlan)
	} else {
		showError(w, r, "Unable to determine an appropriate plan", 500)
	}
	log.Printf("done handling request")
}

func welcome(w http.ResponseWriter, r *http.Request) {
	mustEncode(w, map[string]interface{}{
		"tuqqedin": "relax i'm all tuqqedin",
		"version":  "tuqqedin 0.0",
	})
}
