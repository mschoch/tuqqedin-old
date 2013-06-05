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
	"reflect"
	"testing"

	"github.com/couchbaselabs/tuqqedin/stats"
)

func TestEvaluateLiteral(t *testing.T) {
	tests := []struct {
		input  Expression
		output interface{}
	}{
		{NewLiteralBool(true), true},
		{NewLiteralBool(false), false},
		{NewLiteralNumber(1.0), 1.0},
		{NewLiteralNumber(3.14), 3.14},
		{NewLiteralString("couchbase"), "couchbase"},
		{NewLiteralArray([]Expression{NewLiteralNumber(1.0)}), []interface{}{1.0}},
		{NewLiteralArray([]Expression{NewLiteralNumber(1.0), NewLiteralBool(false)}), []interface{}{1.0, false}},
		{NewLiteralArray([]Expression{NewLiteralNumber(1.0), NewLiteralBool(false), NewLiteralString("bob")}), []interface{}{1.0, false, "bob"}},
		{NewLiteralObject(map[string]Expression{"name": NewLiteralString("bob")}), map[string]interface{}{"name": "bob"}},
		{NewLiteralObject(map[string]Expression{"user": NewLiteralString("test"), "age": NewLiteralNumber(27.0)}), map[string]interface{}{"age": 27.0, "user": "test"}},
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

func TestLiteralBoolSelectivityNoStats(t *testing.T) {
	tests := []struct {
		input  BooleanExpression
		output interface{}
	}{
		{NewLiteralBool(true), 1.0},
		{NewLiteralBool(false), 0.0},
	}

	for _, x := range tests {
		result := x.input.GetSelectivity(map[string]stats.PathStatistics{})
		if !reflect.DeepEqual(result, x.output) {
			t.Errorf("Expected %v, got %v", x.output, result)
		}
	}
}
