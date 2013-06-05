//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package datasource

import (
	"fmt"
	"log"

	"github.com/couchbaselabs/go-couchbase"
	"github.com/couchbaselabs/tuqqedin/ast"
)

type CouchbaseAllDocsAccessPath struct {
	dataSource *CouchbaseDataSource
	ddoc       string
	view       string
}

func NewCouchbaseAllDocsAccessPath(dataSource *CouchbaseDataSource) *CouchbaseAllDocsAccessPath {
	return &CouchbaseAllDocsAccessPath{
		dataSource: dataSource,
		ddoc:       "",
		view:       "_all_docs",
	}
}

func (this *CouchbaseAllDocsAccessPath) Name() string {
	return this.view
}

func (this *CouchbaseAllDocsAccessPath) Keys() []string {
	return []string{}
}

func (this *CouchbaseAllDocsAccessPath) DataSource() DataSource {
	return this.dataSource
}

func (this *CouchbaseAllDocsAccessPath) ReturnsAll() bool {
	return true
}

func (this *CouchbaseAllDocsAccessPath) Matches([]ast.BooleanExpression) bool {
	return false
}

func (this *CouchbaseAllDocsAccessPath) Scan(output DocumentChannel, cancel CancelChannel, options map[string]interface{}) {

	defer close(output)

	options["include_docs"] = true
	viewRowsChannel := make(chan couchbase.ViewRow)
	go WalkViewInBatches(viewRowsChannel, this.dataSource.bucket, this.ddoc, this.view, options, BATCH_SIZE)
	for row := range viewRowsChannel {
		rowdoc := (*row.Doc).(map[string]interface{})
		rowdoc["doc"] = rowdoc["json"]
		delete(rowdoc, "json")
		output <- rowdoc
	}

}

func (this *CouchbaseAllDocsAccessPath) UpdateStats() {

	vres, err := this.dataSource.bucket.View(this.ddoc, this.view, map[string]interface{}{"limit": 0})
	if err != nil {
		log.Printf("Unable to determine cardinality of view, defaulting to MAX")
	} else {
		this.dataSource.rows = vres.TotalRows
	}

}

func (this *CouchbaseAllDocsAccessPath) String() string {
	return fmt.Sprintf("%v", this.Name())
}
