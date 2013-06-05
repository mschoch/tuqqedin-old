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
	"math"
	"sort"

	"github.com/couchbaselabs/tuqqedin/ast"
	"github.com/couchbaselabs/tuqqedin/datasource"
)

type Order struct {
	source        Operator
	outputChannel OutputChannel
	orderBy       []ast.OrderedExpression
	output        []Output
}

func NewOrder(source Operator, orderBy []ast.OrderedExpression) *Order {
	return &Order{
		source:        source,
		outputChannel: make(OutputChannel),
		orderBy:       orderBy,
		output:        make([]Output, 0),
	}
}

func (this *Order) GetOutputChannel() OutputChannel {
	return this.outputChannel
}

func (this *Order) Run() {
	defer close(this.outputChannel)

	// start the source
	go this.source.Run()

	// store all the rows
	for row := range this.source.GetOutputChannel() {
		this.output = append(this.output, row)
	}

	// sort
	sort.Sort(this)

	// write the output
	for _, row := range this.output {
		this.outputChannel <- row
	}
}

func (this *Order) Explain() map[string]interface{} {
	rv := map[string]interface{}{
		"type":           "order",
		"by":             this.orderBy,
		"estimated_rows": this.EstimatedRows(),
		"cost":           this.Cost(),
	}
	if this.Source() != nil {
		rv["source"] = this.Source().Explain()
	}
	return rv
}
func (this *Order) Cancel() {}

func (this *Order) Cost() float64 {
	sourceRows := this.source.EstimatedRows()
	return float64(sourceRows) * math.Log10(float64(sourceRows))
}

func (this *Order) EstimatedRows() int {
	return this.source.EstimatedRows()
}

func (this *Order) TotalCost() float64 {
	return this.Cost() + this.source.TotalCost()
}

func (this *Order) String() string {
	return OperatorToString(this)
}

func (this *Order) Source() Operator {
	return this.source
}

// sort.Interface interface

func (this *Order) Len() int      { return len(this.output) }
func (this *Order) Swap(i, j int) { this.output[i], this.output[j] = this.output[j], this.output[i] }
func (this *Order) Less(i, j int) bool {
	left := this.output[i]
	right := this.output[j]

	var leftContext, rightContext ast.Context
	switch left := left.(type) {
	case datasource.Document:
		leftContext = ast.NewContext(left)
		switch right := right.(type) {
		case datasource.Document:
			rightContext = ast.NewContext(right)
		default:
			panic(fmt.Sprintf("Non-map rows not currently supported (saw %T)", right))
		}
	default:
		panic(fmt.Sprintf("Non-map rows not currently supported (saw %T)", left))
	}

	for _, oe := range this.orderBy {
		leftVal, err := oe.Expression().Evaluate(leftContext)
		if err != nil {
			log.Printf("Error evaluating expression: %v", err)
			return false
		}
		rightVal, err := oe.Expression().Evaluate(rightContext)
		if err != nil {
			log.Printf("Error evaluating expression: %v", err)
			return false
		}

		result := ast.CollateJSON(leftVal, rightVal)
		if result != 0 {
			if oe.Order() && result < 0 {
				return true
			} else if !oe.Order() && result > 0 {
				return true
			} else {
				return false
			}
		}
		// at this level they are the same, keep going
	}

	// if we go to this point the order expressions could not differentiate between the elements
	return false
}
