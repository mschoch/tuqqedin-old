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
	"fmt"
	"log"

	"github.com/couchbaselabs/tuqqedin/ast"
	"github.com/couchbaselabs/tuqqedin/datasource"
)

type Project struct {
	source        Operator
	outputChannel OutputChannel
	projection    ast.Expression
}

func NewProject(source Operator, projection ast.Expression) *Project {
	return &Project{
		source:        source,
		outputChannel: make(OutputChannel),
		projection:    projection,
	}
}

func (this *Project) GetOutputChannel() OutputChannel {
	return this.outputChannel
}

func (this *Project) Run() {
	defer close(this.outputChannel)

	// start the source
	go this.source.Run()
DOCUMENT:
	for row := range this.source.GetOutputChannel() {

		if this.projection != nil {
			var context ast.Context
			switch row := row.(type) {
			case datasource.Document:
				//log.Printf("creating context with %v", row)
				context = ast.NewContext(row)
			default:
				panic(fmt.Sprintf("Non-map rows not currently supported (saw %T)", row))
			}

			result, err := this.projection.Evaluate(context)
			if err != nil {
				log.Printf("Error evaluating projection: %v", err)
				continue DOCUMENT
			}
			this.outputChannel <- result
		} else {
			this.outputChannel <- row
		}
	}
}

func (this *Project) Explain() map[string]interface{} {
	rv := map[string]interface{}{
		"type":           "projection",
		"expression":     this.projection,
		"estimated_rows": this.EstimatedRows(),
		"cost":           this.Cost(),
	}
	if this.Source() != nil {
		rv["source"] = this.Source().Explain()
	}
	return rv
}
func (this *Project) Cancel() {}

func (this *Project) Cost() float64 {
	return 0.0
}

func (this *Project) EstimatedRows() int {
	return this.source.EstimatedRows()
}

func (this *Project) TotalCost() float64 {
	return this.Cost() + this.source.TotalCost()
}

func (this *Project) String() string {
	return OperatorToString(this)
}

func (this *Project) Source() Operator {
	return this.source
}
