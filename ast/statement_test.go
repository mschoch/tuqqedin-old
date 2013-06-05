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
)

func TestParseLiteral(t *testing.T) {

	tests := []struct {
		input       map[string]interface{}
		output      Expression
		outputError error
	}{
		{map[string]interface{}{"value": false}, &LiteralBool{false}, nil},
		{map[string]interface{}{"value": 2.5}, &LiteralNumber{2.5}, nil},
		{map[string]interface{}{"value": "bob"}, &LiteralString{"bob"}, nil},
		{map[string]interface{}{"value": []interface{}{map[string]interface{}{"type": "literal", "value": false}}},
			&LiteralArray{[]Expression{&LiteralBool{false}}}, nil},
		{map[string]interface{}{"value": map[string]interface{}{"wow": map[string]interface{}{"type": "literal", "value": false}}},
			&LiteralObject{map[string]Expression{"wow": &LiteralBool{false}}}, nil},
	}

	for _, test := range tests {
		expr, err := parseLiteral(test.input)

		if err != test.outputError {
			t.Errorf("Expected error %v, got %v", test.outputError, err)
		}
		if !reflect.DeepEqual(test.output, expr) {
			t.Errorf("Expected expression %v, got %v", test.output, expr)
		}
	}

}

func TestParseBoolean(t *testing.T) {

	tests := []struct {
		input       map[string]interface{}
		output      Expression
		outputError error
	}{
		{
			map[string]interface{}{
				"type": "and",
				"left": map[string]interface{}{
					"type":     "compare",
					"operator": "eq",
					"left":     map[string]interface{}{"type": "literal", "value": false},
					"right":    map[string]interface{}{"type": "literal", "value": false},
				},
				"right": map[string]interface{}{
					"type":     "compare",
					"operator": "eq",
					"left":     map[string]interface{}{"type": "literal", "value": false},
					"right":    map[string]interface{}{"type": "literal", "value": false},
				},
			},
			&AndOperator{
				BooleanOperator{
					[]BooleanExpression{
						&EqualToOperator{
							BinaryOperator{
								&LiteralBool{false},
								&LiteralBool{false},
							},
						},
						&EqualToOperator{
							BinaryOperator{
								&LiteralBool{false},
								&LiteralBool{false},
							},
						},
					},
				},
			},
			nil,
		},
		{
			map[string]interface{}{
				"type": "not",
				"operand": map[string]interface{}{
					"type":     "compare",
					"operator": "eq",
					"left":     map[string]interface{}{"type": "literal", "value": false},
					"right":    map[string]interface{}{"type": "literal", "value": false},
				},
			},
			&NotOperator{
				&EqualToOperator{
					BinaryOperator{
						&LiteralBool{false},
						&LiteralBool{false},
					},
				},
			},
			nil,
		},
	}

	for _, test := range tests {
		expr, err := parseBooleanExpression(test.input)

		if err != test.outputError {
			t.Errorf("Expected error %v, got %v", test.outputError, err)
		}
		if !reflect.DeepEqual(test.output, expr) {
			t.Errorf("Expected expression %v, got %v", test.output, expr)
		}
	}

}

func TestParseArithmetic(t *testing.T) {

	tests := []struct {
		input       map[string]interface{}
		output      Expression
		outputError error
	}{
		{
			map[string]interface{}{
				"operator": "plus",
				"left":     map[string]interface{}{"type": "literal", "value": 7.2},
				"right":    map[string]interface{}{"type": "literal", "value": 1.5},
			},
			&PlusOperator{
				&LiteralNumber{7.2},
				&LiteralNumber{1.5},
			},
			nil,
		},
		{
			map[string]interface{}{
				"operator": "minus",
				"left":     map[string]interface{}{"type": "literal", "value": 7.2},
				"right":    map[string]interface{}{"type": "literal", "value": 1.5},
			},
			&SubtractOperator{
				&LiteralNumber{7.2},
				&LiteralNumber{1.5},
			},
			nil,
		},
		{
			map[string]interface{}{
				"operator": "mult",
				"left":     map[string]interface{}{"type": "literal", "value": 7.2},
				"right":    map[string]interface{}{"type": "literal", "value": 1.5},
			},
			&MultiplyOperator{
				&LiteralNumber{7.2},
				&LiteralNumber{1.5},
			},
			nil,
		},
		{
			map[string]interface{}{
				"operator": "div",
				"left":     map[string]interface{}{"type": "literal", "value": 7.2},
				"right":    map[string]interface{}{"type": "literal", "value": 1.5},
			},
			&DivideOperator{
				&LiteralNumber{7.2},
				&LiteralNumber{1.5},
			},
			nil,
		},
	}

	for _, test := range tests {
		expr, err := parseArithmetic(test.input)

		if err != test.outputError {
			t.Errorf("Expected error %v, got %v", test.outputError, err)
		}
		if !reflect.DeepEqual(test.output, expr) {
			t.Errorf("Expected expression %v, got %v", test.output, expr)
		}
	}

}

func TestParseCompareExpression(t *testing.T) {

	tests := []struct {
		input       map[string]interface{}
		output      Expression
		outputError error
	}{
		{
			map[string]interface{}{
				"operator": "gt",
				"left":     map[string]interface{}{"type": "literal", "value": 7.2},
				"right":    map[string]interface{}{"type": "literal", "value": 1.5},
			},
			&GreaterThanOperator{
				BinaryOperator{
					&LiteralNumber{7.2},
					&LiteralNumber{1.5},
				},
			},
			nil,
		},
		{
			map[string]interface{}{
				"operator": "gte",
				"left":     map[string]interface{}{"type": "literal", "value": 7.2},
				"right":    map[string]interface{}{"type": "literal", "value": 1.5},
			},
			&GreaterThanOrEqualOperator{
				BinaryOperator{
					&LiteralNumber{7.2},
					&LiteralNumber{1.5},
				},
			},
			nil,
		},
		{
			map[string]interface{}{
				"operator": "lt",
				"left":     map[string]interface{}{"type": "literal", "value": 7.2},
				"right":    map[string]interface{}{"type": "literal", "value": 1.5},
			},
			&LessThanOperator{
				BinaryOperator{
					&LiteralNumber{7.2},
					&LiteralNumber{1.5},
				},
			},
			nil,
		},
		{
			map[string]interface{}{
				"operator": "lte",
				"left":     map[string]interface{}{"type": "literal", "value": 7.2},
				"right":    map[string]interface{}{"type": "literal", "value": 1.5},
			},
			&LessThanOrEqualOperator{
				BinaryOperator{
					&LiteralNumber{7.2},
					&LiteralNumber{1.5},
				},
			},
			nil,
		},
		{
			map[string]interface{}{
				"operator": "eq",
				"left":     map[string]interface{}{"type": "literal", "value": 7.2},
				"right":    map[string]interface{}{"type": "literal", "value": 1.5},
			},
			&EqualToOperator{
				BinaryOperator{
					&LiteralNumber{7.2},
					&LiteralNumber{1.5},
				},
			},
			nil,
		},
		{
			map[string]interface{}{
				"operator": "neq",
				"left":     map[string]interface{}{"type": "literal", "value": 7.2},
				"right":    map[string]interface{}{"type": "literal", "value": 1.5},
			},
			&NotEqualToOperator{
				BinaryOperator{
					&LiteralNumber{7.2},
					&LiteralNumber{1.5},
				},
			},
			nil,
		},
	}

	for _, test := range tests {
		left, right, err := parseBinaryOperatorArguments(test.input)
		if err != nil {
			t.Errorf("Unexpected error")
		}
		expr, err := parseCompareExpression(left, right, test.input)

		if err != test.outputError {
			t.Errorf("Expected error %v, got %v", test.outputError, err)
		}
		if !reflect.DeepEqual(test.output, expr) {
			t.Errorf("Expected expression %v, got %v", test.output, expr)
		}
	}

}

func TestValidateTopLevel(t *testing.T) {

	tests := []struct {
		input       map[string]interface{}
		output      map[string]interface{}
		outputError error
	}{
		{
			map[string]interface{}{"type": "cbqast", "version": "1", "statement": map[string]interface{}{}},
			map[string]interface{}{},
			nil,
		},
	}

	for _, test := range tests {
		expr, err := validateTopLevel(test.input)

		if err != test.outputError {
			t.Errorf("Expected error %v, got %v", test.outputError, err)
		}
		if !reflect.DeepEqual(test.output, expr) {
			t.Errorf("Expected expression %v, got %v", test.output, expr)
		}
	}

}

func TestParseStatement(t *testing.T) {

	tests := []struct {
		input       map[string]interface{}
		output      Statement
		outputError error
	}{
		{
			map[string]interface{}{
				"type": "select",
				"where": map[string]interface{}{
					"type":     "compare",
					"operator": "gt",
					"left":     map[string]interface{}{"type": "literal", "value": 7.2},
					"right":    map[string]interface{}{"type": "literal", "value": 1.5},
				},
			},
			&SelectStatement{
				From: make([]DataSource, 0),
				Where: &GreaterThanOperator{
					BinaryOperator{
						&LiteralNumber{7.2},
						&LiteralNumber{1.5},
					},
				},
			},
			nil,
		},
	}

	for _, test := range tests {
		expr, err := parseStatementJSON(test.input)

		if err != test.outputError {
			t.Errorf("Expected error %v, got %v", test.outputError, err)
		}
		if !reflect.DeepEqual(test.output, expr) {
			t.Errorf("Expected expression %v, got %v", test.output, expr)
		}
	}

}
