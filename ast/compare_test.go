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

func TestCompare(t *testing.T) {

	numberSixty := NewLiteralNumber(60.0)
	numberNine := NewLiteralNumber(9.0)

	tests := []struct {
		input  Expression
		output interface{}
	}{
		{NewGreaterThanOperator(numberSixty, numberSixty), false},
		{NewGreaterThanOperator(numberSixty, numberNine), true},
		{NewGreaterThanOperator(numberNine, numberSixty), false},

		{NewGreaterThanOrEqualOperator(numberSixty, numberSixty), true},
		{NewGreaterThanOrEqualOperator(numberSixty, numberNine), true},
		{NewGreaterThanOrEqualOperator(numberNine, numberSixty), false},

		{NewLessThanOperator(numberSixty, numberSixty), false},
		{NewLessThanOperator(numberSixty, numberNine), false},
		{NewLessThanOperator(numberNine, numberSixty), true},

		{NewLessThanOrEqualOperator(numberSixty, numberSixty), true},
		{NewLessThanOrEqualOperator(numberSixty, numberNine), false},
		{NewLessThanOrEqualOperator(numberNine, numberSixty), true},

		{NewEqualToOperator(numberSixty, numberSixty), true},
		{NewEqualToOperator(numberSixty, numberNine), false},
		{NewEqualToOperator(numberNine, numberSixty), false},

		{NewNotEqualToOperator(numberSixty, numberSixty), false},
		{NewNotEqualToOperator(numberSixty, numberNine), true},
		{NewNotEqualToOperator(numberNine, numberSixty), true},
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

func TestCompareSelectivityNoStats(t *testing.T) {

	numberSixty := NewLiteralNumber(60.0)

	tests := []struct {
		input  BooleanExpression
		output interface{}
	}{
		{NewGreaterThanOperator(numberSixty, numberSixty), 1.0 / 3.0},
		{NewGreaterThanOrEqualOperator(numberSixty, numberSixty), 1.0 / 3.0},
		{NewLessThanOperator(numberSixty, numberSixty), 1.0 / 3.0},
		{NewLessThanOrEqualOperator(numberSixty, numberSixty), 1.0 / 3.0},
		{NewEqualToOperator(numberSixty, numberSixty), 1.0 / 10.0},
		{NewNotEqualToOperator(numberSixty, numberSixty), 9.0 / 10.0},
	}

	for _, x := range tests {
		result := x.input.GetSelectivity(map[string]stats.PathStatistics{})
		if !reflect.DeepEqual(result, x.output) {
			t.Errorf("Expected %v, got %v", x.output, result)
		}
	}
}

func TestCompareSelectivityWithStatsNoQuantiles(t *testing.T) {

	sixtyValue := 60.0
	numberSixty := NewLiteralNumber(sixtyValue)
	fortyValue := 40.0
	numberForty := NewLiteralNumber(fortyValue)

	abvStats := stats.DefaultPathStats(0, 100)
	abvStats.Rows = 50
	abvStats.DistinctValues = 10
	abvStats.MostFrequentValues.Consider(sixtyValue, 6)

	pathStatistics := map[string]stats.PathStatistics{
		"doc.abv": abvStats,
	}

	tests := []struct {
		input  BooleanExpression
		output interface{}
	}{
		// when value is in freq vals, we should get precise estimate
		{NewEqualToOperator(NewProperty("doc.abv"), numberSixty), 6.0 / 50.0},
		{NewNotEqualToOperator(NewProperty("doc.abv"), numberSixty), 44.0 / 50.0},
		// if value is not in freq vals, and we dont have quantiles
		// it should fall back to defaults
		{NewEqualToOperator(NewProperty("doc.abv"), numberForty), 1.0 / 10.0},
		{NewNotEqualToOperator(NewProperty("doc.abv"), numberForty), 9.0 / 10.0},
	}

	for _, x := range tests {
		result := x.input.GetSelectivity(pathStatistics)
		if !reflect.DeepEqual(result, x.output) {
			t.Errorf("Expected %v for %v, got %v", x.output, x.input, result)
		}
	}
}

func TestCompareSelectivityWithStatsAndQuantiles(t *testing.T) {

	sixtyValue := 60.0
	numberSixty := NewLiteralNumber(sixtyValue)
	fortyValue := 40.0
	numberForty := NewLiteralNumber(fortyValue)

	abvStats := stats.DefaultPathStats(0, 100)
	abvStats.Rows = 50
	abvStats.DistinctValues = 10
	abvStats.MostFrequentValues.Consider(sixtyValue, 6)
	abvStats.Quantiles = []stats.QuantileRange{
		stats.QuantileRange{Start: 0.0, End: 20.0, Count: 10},
		stats.QuantileRange{Start: 21.0, End: 40.0, Count: 10},
		stats.QuantileRange{Start: 41.0, End: 60.0, Count: 10},
		stats.QuantileRange{Start: 61.0, End: 80.0, Count: 10},
		stats.QuantileRange{Start: 81.0, End: 100.0, Count: 10},
	}

	pathStatistics := map[string]stats.PathStatistics{
		"doc.abv": abvStats,
	}

	tests := []struct {
		input  BooleanExpression
		output interface{}
	}{
		// when value is in freq vals, we should get precise estimate
		{NewEqualToOperator(NewProperty("doc.abv"), numberSixty), 6.0 / 50.0},
		{NewNotEqualToOperator(NewProperty("doc.abv"), numberSixty), 44.0 / 50.0},
		// this time it should fall back to the quantiles
		{NewEqualToOperator(NewProperty("doc.abv"), numberForty), 10.0 / 50.0},
		{NewNotEqualToOperator(NewProperty("doc.abv"), numberForty), 40.0 / 50.0},
		// now test the ranges
		{NewLessThanOperator(NewProperty("doc.abv"), numberSixty), 30.0 / 50.0},
		{NewGreaterThanOperator(NewProperty("doc.abv"), numberSixty), 30.0 / 50.0},
		// note these two find exact values for the equals portion
		{NewLessThanOrEqualOperator(NewProperty("doc.abv"), numberSixty), 36.0 / 50.0},
		{NewGreaterThanOrEqualOperator(NewProperty("doc.abv"), numberSixty), 36.0 / 50.0},
		// these two should not
		{NewLessThanOrEqualOperator(NewProperty("doc.abv"), numberForty), 30.0 / 50.0},
		{NewGreaterThanOrEqualOperator(NewProperty("doc.abv"), numberForty), 50.0 / 50.0},
	}

	for _, x := range tests {
		result := x.input.GetSelectivity(pathStatistics)
		if !reflect.DeepEqual(result, x.output) {
			t.Errorf("Expected %v for %v, got %v", x.output, x.input, result)
		}
	}
}
