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
	"fmt"
	"log"
	"net/http"

	"github.com/couchbaselabs/tuqqedin/plan"
)

type Executor interface {
	ExecutePlan(w http.ResponseWriter, plan plan.Operator)
}

type CouchbaseExecutor struct {
}

func NewCouchbaseExecutor() *CouchbaseExecutor {
	return &CouchbaseExecutor{}
}

func (this *CouchbaseExecutor) ExecutePlan(w http.ResponseWriter, plan plan.Operator) {

	// get reference to the output channel
	output := plan.GetOutputChannel()

	// start the plan
	go plan.Run()

	// start the output stream
	fmt.Fprint(w, "{\n")
	fmt.Fprint(w, "    \"resultset\": [\n")

	first := true
	count := 0
	// read the rows returned
	for row := range output {
		if !first {
			fmt.Fprint(w, ",\n")
		}
		body, err := json.MarshalIndent(row, "        ", "    ")
		if err != nil {
			log.Printf("Unable to format result to display %#v, %v", row, err)
		} else {
			fmt.Fprintf(w, "        %v", string(body))
		}
		first = false
		count++
	}
	fmt.Fprint(w, "\n    ],\n")
	fmt.Fprintf(w, "    \"total_rows\": %d\n", count)
	fmt.Fprint(w, "}\n")

}
