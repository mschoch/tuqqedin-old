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

type Filter struct {
	source         Operator
	outputChannel  OutputChannel
	booleanFactors []ast.BooleanExpression
}

func NewFilter(source Operator, booleanFactors []ast.BooleanExpression) *Filter {
	return &Filter{
		source:         source,
		outputChannel:  make(OutputChannel),
		booleanFactors: booleanFactors,
	}
}

func (this *Filter) GetOutputChannel() OutputChannel {
	return this.outputChannel
}

func (this *Filter) Run() {
	defer close(this.outputChannel)

	// FIXME should ensure booleanFactors are sorted
	// from cheapest to most expensive to compute

	// start the source
	go this.source.Run()
DOCUMENT:
	for row := range this.source.GetOutputChannel() {
		var context ast.Context
		switch row := row.(type) {
		case datasource.Document:
			context = ast.NewContext(row)
		default:
			panic(fmt.Sprintf("Non-map rows not currently supported (saw %T)", row))
		}

		for _, booleanFactor := range this.booleanFactors {
			result, err := booleanFactor.EvaluateBoolean(context)
			if err != nil {
				log.Printf("Error evaluating boolean factor: %v", err)
				continue DOCUMENT
			}
			if !result {
				continue DOCUMENT
			}
		}
		this.outputChannel <- row
	}
}

func (this *Filter) Explain() map[string]interface{} {
	rv := map[string]interface{}{
		"type":           "filter",
		"booleanFactors": this.booleanFactors,
		"estimated_rows": this.EstimatedRows(),
		"cost":           this.Cost(),
	}
	if this.Source() != nil {
		rv["source"] = this.Source().Explain()
	}
	return rv
}
func (this *Filter) Cancel() {}

func (this *Filter) Cost() float64 {
	sourceRows := this.source.EstimatedRows()
	return float64(sourceRows) * CPU_COST
}

func (this *Filter) EstimatedRows() int {
	// FIXME this is wrong in most cases
	// the filter is further reducing the number of rows
	return this.source.EstimatedRows()
}

func (this *Filter) TotalCost() float64 {
	return this.Cost() + this.source.TotalCost()
}

func (this *Filter) String() string {
	return OperatorToString(this)
}

func (this *Filter) Source() Operator {
	return this.source
}
