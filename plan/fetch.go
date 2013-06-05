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
	"log"

	"github.com/couchbaselabs/tuqqedin/datasource"
)

type Fetch struct {
	source        Operator
	outputChannel OutputChannel
	dataSource    datasource.DataSource
}

func NewFetch(source Operator, dataSource datasource.DataSource) *Fetch {
	return &Fetch{
		source:        source,
		outputChannel: make(OutputChannel),
		dataSource:    dataSource,
	}
}

func (this *Fetch) GetOutputChannel() OutputChannel {
	return this.outputChannel
}

func (this *Fetch) Run() {
	defer close(this.outputChannel)

	// start the source
	go this.source.Run()
DOCUMENT:
	for row := range this.source.GetOutputChannel() {
		switch row := row.(type) {
		case datasource.Document:
			docMeta := row["meta"]
			switch docMeta := docMeta.(type) {
			case map[string]interface{}:
				docID := docMeta["id"]
				switch docID := docID.(type) {
				case string:
					// lookup document
					doc, err := this.dataSource.Fetch(docID)
					if err != nil {
						log.Printf("error retrieving doc")
						continue DOCUMENT
					}
					row["doc"] = doc
					// now clear out the additional meta information
					// we got from the index, it may not be correct
					delete(docMeta, "rev")
					delete(docMeta, "flags")
					delete(docMeta, "expiration")
					this.outputChannel <- row
				}
			}
		default:
			panic("Non-map rows not currently supported")
		}
	}
}

func (this *Fetch) Explain() map[string]interface{} {
	rv := map[string]interface{}{
		"type":           "fetch",
		"estimated_rows": this.EstimatedRows(),
		"cost":           this.Cost(),
	}
	if this.Source() != nil {
		rv["source"] = this.Source().Explain()
	}
	return rv
}
func (this *Fetch) Cancel() {}

func (this *Fetch) Cost() float64 {
	// esimate this differently than views, connection is maintained
	// and considerably faster than view query
	return float64(this.EstimatedRows()) + NETWORK_COST
}

func (this *Fetch) EstimatedRows() int {
	// does not usually reduce rows, will return same as source
	// (can reduce rows if docs in index have been deleted)
	return this.source.EstimatedRows()
}

func (this *Fetch) TotalCost() float64 {
	return this.Cost() + this.source.TotalCost()
}

func (this *Fetch) String() string {
	return OperatorToString(this)
}

func (this *Fetch) Source() Operator {
	return this.source
}
