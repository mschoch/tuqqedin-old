//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package ast

import (
	//	"log"
	"reflect"
	"testing"

	"github.com/couchbaselabs/tuqqedin/stats"
)

func TestBoolean(t *testing.T) {

	booleanTrue := NewLiteralBool(true)
	booleanFalse := NewLiteralBool(false)

	tests := []struct {
		input  Expression
		output interface{}
	}{
		{NewAndOperator([]BooleanExpression{booleanTrue, booleanTrue}), true},
		{NewAndOperator([]BooleanExpression{booleanTrue, booleanFalse}), false},
		{NewAndOperator([]BooleanExpression{booleanFalse, booleanTrue}), false},
		{NewAndOperator([]BooleanExpression{booleanFalse, booleanFalse}), false},

		{NewAndOperator([]BooleanExpression{booleanTrue, booleanTrue, booleanTrue, booleanTrue, booleanTrue}), true},
		{NewAndOperator([]BooleanExpression{booleanTrue, booleanTrue, booleanTrue, booleanTrue, booleanFalse}), false},

		{NewOrOperator([]BooleanExpression{booleanTrue, booleanTrue}), true},
		{NewOrOperator([]BooleanExpression{booleanTrue, booleanFalse}), true},
		{NewOrOperator([]BooleanExpression{booleanFalse, booleanTrue}), true},
		{NewOrOperator([]BooleanExpression{booleanFalse, booleanFalse}), false},

		{NewOrOperator([]BooleanExpression{booleanFalse, booleanFalse, booleanFalse, booleanFalse, booleanTrue}), true},
		{NewOrOperator([]BooleanExpression{booleanFalse, booleanFalse, booleanFalse, booleanFalse, booleanFalse}), false},

		{NewNotOperator(booleanTrue), false},
		{NewNotOperator(booleanFalse), true},
	}

	for _, x := range tests {
		result, err := x.input.Evaluate(nil)
		if err != nil {
			t.Fatalf("Error evaluating expression: %v", err)
		}
		if !reflect.DeepEqual(result, x.output) {
			t.Errorf("Expected %t %v, got %t %v", x.output, x.output, result, result)
		}
	}

}

func TestNNF(t *testing.T) {

	booleanTrue := NewLiteralBool(true)
	booleanFalse := NewLiteralBool(false)

	tests := []struct {
		input  BooleanExpression
		output BooleanExpression
	}{
		// distribute NOT over AND and OR
		{NewNotOperator(NewAndOperator([]BooleanExpression{booleanTrue, booleanTrue})), NewOrOperator([]BooleanExpression{NewNotOperator(booleanTrue), NewNotOperator(booleanTrue)})},
		{NewNotOperator(NewOrOperator([]BooleanExpression{booleanTrue, booleanFalse})), NewAndOperator([]BooleanExpression{NewNotOperator(booleanTrue), NewNotOperator(booleanFalse)})},

		// eliminate double NOT
		{NewNotOperator(NewNotOperator(booleanTrue)), booleanTrue},

		// reduce tripple NOT into single NOT
		{NewNotOperator(NewNotOperator(NewNotOperator(booleanTrue))), NewNotOperator(booleanTrue)},
	}

	for _, x := range tests {
		nnf := x.input.NegationNormalForm()
		if !reflect.DeepEqual(nnf, x.output) {
			t.Errorf("Expected %t %v, got %t %v", x.output, x.output, nnf, nnf)
		}

		result, err := x.input.Evaluate(nil)
		if err != nil {
			t.Fatalf("Error evaluating expression: %v", err)
		}
		result2, err := nnf.Evaluate(nil)
		if !reflect.DeepEqual(result, result2) {
			t.Errorf("Expected same result %t %v, got %t %v", result, result, result2, result2)
		}
	}
}

func TestCNF(t *testing.T) {

	booleanTrue := NewLiteralBool(true)
	booleanFalse := NewLiteralBool(false)

	tests := []struct {
		input  BooleanExpression
		output BooleanExpression
	}{

		// combine complex AND into simple AND
		{
			NewAndOperator([]BooleanExpression{NewAndOperator([]BooleanExpression{booleanTrue, booleanTrue}), booleanTrue}),
			NewAndOperator([]BooleanExpression{booleanTrue, booleanTrue, booleanTrue}),
		},

		// combine complex OR into simple OR
		{
			NewOrOperator([]BooleanExpression{NewOrOperator([]BooleanExpression{booleanTrue, booleanFalse}), booleanTrue}),
			NewOrOperator([]BooleanExpression{booleanTrue, booleanFalse, booleanTrue}),
		},

		//
		{
			NewOrOperator([]BooleanExpression{NewAndOperator([]BooleanExpression{booleanTrue, booleanTrue}), booleanFalse}),
			NewAndOperator([]BooleanExpression{
				NewOrOperator([]BooleanExpression{booleanFalse, booleanTrue}),
				NewOrOperator([]BooleanExpression{booleanFalse, booleanTrue}),
			}),
		},
	}

	for _, x := range tests {
		nnf := x.input.NegationNormalForm()
		cnf := nnf.ConjunctiveNormalForm()
		if !reflect.DeepEqual(cnf, x.output) {
			t.Errorf("Expected %v, got %v", x.output, cnf)
		}

		result, err := x.input.Evaluate(nil)
		if err != nil {
			t.Fatalf("Error evaluating expression: %v", err)
		}
		result2, err := cnf.Evaluate(nil)
		if !reflect.DeepEqual(result, result2) {
			t.Errorf("Expected same result %v, got %v", result, result, result2, result2)
		}
	}
}

func TestCNFDocument(t *testing.T) {

	sampleDocument := map[string]interface{}{
		"name": "will",
		"age":  39.0,
		"address": map[string]interface{}{
			"city": "New York",
		},
		"children": []interface{}{
			map[string]interface{}{
				"name": "bob",
			},
			map[string]interface{}{
				"name": "jane",
			},
		},
	}

	context := NewContext(sampleDocument)

	nameProperty := NewProperty("name")
	ageProperty := NewProperty("age")
	literal18 := NewLiteralNumber(18.1)
	literalWill := NewLiteralString("will")

	tests := []struct {
		input  BooleanExpression
		output BooleanExpression
	}{
		{
			NewOrOperator([]BooleanExpression{
				NewAndOperator([]BooleanExpression{
					NewGreaterThanOperator(ageProperty, literal18),
					NewEqualToOperator(nameProperty, literalWill)}),
				NewLessThanOperator(ageProperty, literal18)}),
			NewAndOperator([]BooleanExpression{
				NewOrOperator([]BooleanExpression{NewLessThanOperator(ageProperty, literal18), NewGreaterThanOperator(ageProperty, literal18)}),
				NewOrOperator([]BooleanExpression{NewLessThanOperator(ageProperty, literal18), NewEqualToOperator(nameProperty, literalWill)}),
			}),
		},
	}

	for _, x := range tests {
		nnf := x.input.NegationNormalForm()
		cnf := nnf.ConjunctiveNormalForm()
		if !reflect.DeepEqual(cnf, x.output) {
			t.Errorf("Expected %v, got %v", x.output, cnf)
		}

		result, err := x.input.Evaluate(context)
		if err != nil {
			t.Fatalf("Error evaluating expression: %v", err)
		}
		result2, err := cnf.Evaluate(context)
		if !reflect.DeepEqual(result, result2) {
			t.Errorf("Expected same result %v, got %v", result, result, result2, result2)
		}
	}

}

func TestBooleanSelectivityNoStats(t *testing.T) {

	numberSixty := NewLiteralNumber(60.0)

	tests := []struct {
		input  BooleanExpression
		output interface{}
	}{
		{
			NewAndOperator([]BooleanExpression{
				NewGreaterThanOperator(numberSixty, numberSixty),
				NewGreaterThanOperator(numberSixty, numberSixty),
			}),
			1.0 / 9.0,
		},
		{
			NewOrOperator([]BooleanExpression{
				NewGreaterThanOperator(numberSixty, numberSixty),
				NewGreaterThanOperator(numberSixty, numberSixty),
			}),
			((2.0 / 3.0) - (1.0 / 9.0)),
		},
		{
			NewNotOperator(NewGreaterThanOperator(numberSixty, numberSixty)),
			0.6666666666666667,
		},
	}

	for _, x := range tests {
		result := x.input.GetSelectivity(map[string]stats.PathStatistics{})
		if !reflect.DeepEqual(result, x.output) {
			t.Errorf("Expected %v, got %v", x.output, result)
		}
	}
}
