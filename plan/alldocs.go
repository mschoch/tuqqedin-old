//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package plan

import (
	"github.com/couchbaselabs/tuqqedin/datasource"
)

type AllDocsScanner struct {
	accessPath    *datasource.CouchbaseAllDocsAccessPath
	outputChannel OutputChannel
	cancelChannel datasource.CancelChannel
}

func NewAllDocsScanner(accessPath *datasource.CouchbaseAllDocsAccessPath) *AllDocsScanner {
	return &AllDocsScanner{
		accessPath:    accessPath,
		outputChannel: make(OutputChannel),
		cancelChannel: make(datasource.CancelChannel),
	}
}

func (this *AllDocsScanner) GetOutputChannel() OutputChannel {
	return this.outputChannel
}

func (this *AllDocsScanner) Run() {
	defer close(this.outputChannel)

	docChannel := make(datasource.DocumentChannel)

	go this.accessPath.Scan(docChannel, this.cancelChannel, map[string]interface{}{})
	for doc := range docChannel {
		this.outputChannel <- doc
	}

}

func (this *AllDocsScanner) Explain() map[string]interface{} {
	return map[string]interface{}{
		"type":           "scan",
		"index":          this.accessPath.Name(),
		"estimated_rows": this.EstimatedRows(),
		"cost":           this.Cost(),
	}
}

func (this *AllDocsScanner) Cancel() {
	close(this.cancelChannel)
}

func (this *AllDocsScanner) Cost() float64 {
	return (float64(this.EstimatedRows()) / datasource.BATCH_SIZE) * NETWORK_COST
}

func (this *AllDocsScanner) EstimatedRows() int {
	return this.accessPath.DataSource().Rows()
}

func (this *AllDocsScanner) TotalCost() float64 {
	return this.Cost()
}

func (this *AllDocsScanner) String() string {
	return OperatorToString(this)
}

func (this *AllDocsScanner) Source() Operator {
	return nil
}
