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
	"fmt"

	"github.com/couchbaselabs/tuqqedin/stats"
)

type BinaryOperator struct {
	left  Expression
	right Expression
}

func (this *BinaryOperator) compare(context Context) (int, error) {
	lv, err := this.left.Evaluate(context)
	if err != nil {
		return 0, err
	}
	rv, err := this.right.Evaluate(context)
	if err != nil {
		return 0, err
	}

	return CollateJSON(lv, rv), nil
}

func (this *BinaryOperator) ReferencedProperties() []Property {
	rv := make([]Property, 0, 0)

	rp := this.left.ReferencedProperties()
	for _, rpv := range rp {
		rv = append(rv, rpv)
	}
	rp = this.right.ReferencedProperties()
	for _, rpv := range rp {
		rv = append(rv, rpv)
	}

	return rv
}

func (this *BinaryOperator) IsSargable() bool {
	switch this.left.(type) {
	case *Property:
		switch this.right.(type) {
		case *LiteralBool, *LiteralNumber, *LiteralString:
			return true
		default:
			return false
		}
	case *LiteralBool, *LiteralNumber, *LiteralString:
		switch this.right.(type) {
		case *Property:
			return true
		default:
			return false
		}
	}
	return false
}

func (this *BinaryOperator) GetSargProperty() *Property {
	switch ltype := this.left.(type) {
	case *Property:
		return ltype
	case *LiteralBool, *LiteralNumber, *LiteralString:
		switch rtype := this.right.(type) {
		case *Property:
			return rtype
		}

	}
	return nil
}

func (this *BinaryOperator) GetSargValue() (interface{}, error) {
	switch this.left.(type) {
	case *Property:
		return this.right.Evaluate(nil)
	case *LiteralBool, *LiteralNumber, *LiteralString:
		return this.left.Evaluate(nil)
	}
	return nil, nil
}

type GreaterThanOperator struct {
	BinaryOperator
}

func NewGreaterThanOperator(left, right Expression) *GreaterThanOperator {
	return &GreaterThanOperator{
		BinaryOperator{
			left:  left,
			right: right,
		},
	}
}

func (this *GreaterThanOperator) EvaluateBoolean(context Context) (bool, error) {
	compare, err := this.BinaryOperator.compare(context)
	if err != nil {
		return false, err
	}
	return compare > 0, nil
}

func (this *GreaterThanOperator) Evaluate(context Context) (interface{}, error) {
	return this.EvaluateBoolean(context)
}

func (this *GreaterThanOperator) ConjunctiveNormalForm() BooleanExpression {
	return this
}

func (this *GreaterThanOperator) NegationNormalForm() BooleanExpression {
	return this
}

func (this *GreaterThanOperator) DistributeNot() BooleanExpression {
	return NewNotOperator(this)
}

func (this *GreaterThanOperator) ConvertToBooleanFactors() []BooleanExpression {
	return []BooleanExpression{this}
}

func (this *GreaterThanOperator) String() string {
	return fmt.Sprintf("%v > %v", this.left, this.right)
}

func (this *GreaterThanOperator) GetSelectivity(pathStats map[string]stats.PathStatistics) float64 {
	rv := 1.0 / 3.0

	// logic to improve this with pathStats
	if this.IsSargable() {
		path := this.GetSargProperty().Path
		pathStat, ok := pathStats[path]
		if ok {
			sargval, err := this.GetSargValue()
			if err == nil {
				rowsEstimate := RowsGreaterThanValue(pathStat, sargval)
				if rowsEstimate != -1 {
					return float64(rowsEstimate / float64(pathStat.Rows))
				}
			}
		}
	}

	return rv
}

type GreaterThanOrEqualOperator struct {
	BinaryOperator
}

func NewGreaterThanOrEqualOperator(left, right Expression) *GreaterThanOrEqualOperator {
	return &GreaterThanOrEqualOperator{
		BinaryOperator{
			left:  left,
			right: right,
		},
	}
}

func (this *GreaterThanOrEqualOperator) EvaluateBoolean(context Context) (bool, error) {
	compare, err := this.BinaryOperator.compare(context)
	if err != nil {
		return false, err
	}
	return compare >= 0, nil
}

func (this *GreaterThanOrEqualOperator) Evaluate(context Context) (interface{}, error) {
	return this.EvaluateBoolean(context)
}

func (this *GreaterThanOrEqualOperator) ConjunctiveNormalForm() BooleanExpression {
	return this
}

func (this *GreaterThanOrEqualOperator) NegationNormalForm() BooleanExpression {
	return this
}

func (this *GreaterThanOrEqualOperator) DistributeNot() BooleanExpression {
	return NewNotOperator(this)
}

func (this *GreaterThanOrEqualOperator) ConvertToBooleanFactors() []BooleanExpression {
	return []BooleanExpression{this}
}

func (this *GreaterThanOrEqualOperator) String() string {
	return fmt.Sprintf("%v >= %v", this.left, this.right)
}

func (this *GreaterThanOrEqualOperator) GetSelectivity(pathStats map[string]stats.PathStatistics) float64 {
	rv := 1.0 / 3.0

	// logic to improve this with pathStats
	if this.IsSargable() {
		path := this.GetSargProperty().Path
		pathStat, ok := pathStats[path]
		if ok {
			sargval, err := this.GetSargValue()
			if err == nil {
				lessEstimate := RowsGreaterThanValue(pathStat, sargval)
				equalEstimate := RowsWithValue(pathStat, sargval)
				if lessEstimate != -1 && equalEstimate != -1 {
					sum := lessEstimate + equalEstimate
					// since we're adding 2 estimates
					// ensure we dont go over the max
					if sum > float64(pathStat.Rows) {
						sum = float64(pathStat.Rows)
					}
					return float64(sum / float64(pathStat.Rows))
				}
			}
		}
	}

	return rv
}

type LessThanOperator struct {
	BinaryOperator
}

func NewLessThanOperator(left, right Expression) *LessThanOperator {
	return &LessThanOperator{
		BinaryOperator{
			left:  left,
			right: right,
		},
	}
}

func (this *LessThanOperator) EvaluateBoolean(context Context) (bool, error) {
	compare, err := this.BinaryOperator.compare(context)
	if err != nil {
		return false, err
	}
	return compare < 0, nil
}

func (this *LessThanOperator) Evaluate(context Context) (interface{}, error) {
	return this.EvaluateBoolean(context)
}

func (this *LessThanOperator) ConjunctiveNormalForm() BooleanExpression {
	return this
}

func (this *LessThanOperator) NegationNormalForm() BooleanExpression {
	return this
}

func (this *LessThanOperator) DistributeNot() BooleanExpression {
	return NewNotOperator(this)
}

func (this *LessThanOperator) ConvertToBooleanFactors() []BooleanExpression {
	return []BooleanExpression{this}
}

func (this *LessThanOperator) String() string {
	return fmt.Sprintf("%v < %v", this.left, this.right)
}

func (this *LessThanOperator) GetSelectivity(pathStats map[string]stats.PathStatistics) float64 {
	rv := 1.0 / 3.0

	// logic to improve this with pathStats
	if this.IsSargable() {
		path := this.GetSargProperty().Path
		pathStat, ok := pathStats[path]
		if ok {
			sargval, err := this.GetSargValue()
			if err == nil {
				rowsEstimate := RowsLessThanValue(pathStat, sargval)
				if rowsEstimate != -1 {
					return float64(rowsEstimate / float64(pathStat.Rows))
				}
			}
		}
	}

	return rv
}

type LessThanOrEqualOperator struct {
	BinaryOperator
}

func NewLessThanOrEqualOperator(left, right Expression) *LessThanOrEqualOperator {
	return &LessThanOrEqualOperator{
		BinaryOperator{
			left:  left,
			right: right,
		},
	}
}

func (this *LessThanOrEqualOperator) EvaluateBoolean(context Context) (bool, error) {
	compare, err := this.BinaryOperator.compare(context)
	if err != nil {
		return false, err
	}
	return compare <= 0, nil
}

func (this *LessThanOrEqualOperator) Evaluate(context Context) (interface{}, error) {
	return this.EvaluateBoolean(context)
}

func (this *LessThanOrEqualOperator) ConjunctiveNormalForm() BooleanExpression {
	return this
}

func (this *LessThanOrEqualOperator) NegationNormalForm() BooleanExpression {
	return this
}

func (this *LessThanOrEqualOperator) DistributeNot() BooleanExpression {
	return NewNotOperator(this)
}

func (this *LessThanOrEqualOperator) ConvertToBooleanFactors() []BooleanExpression {
	return []BooleanExpression{this}
}

func (this *LessThanOrEqualOperator) String() string {
	return fmt.Sprintf("%v <= %v", this.left, this.right)
}

func (this *LessThanOrEqualOperator) GetSelectivity(pathStats map[string]stats.PathStatistics) float64 {
	rv := 1.0 / 3.0

	// logic to improve this with pathStats
	if this.IsSargable() {
		path := this.GetSargProperty().Path
		pathStat, ok := pathStats[path]
		if ok {
			sargval, err := this.GetSargValue()
			if err == nil {
				lessEstimate := RowsLessThanValue(pathStat, sargval)
				equalEstimate := RowsWithValue(pathStat, sargval)
				if lessEstimate != -1 && equalEstimate != -1 {
					sum := lessEstimate + equalEstimate
					// since we're adding 2 estimates
					// ensure we dont go over the max
					if sum > float64(pathStat.Rows) {
						sum = float64(pathStat.Rows)
					}
					return float64(sum / float64(pathStat.Rows))
				}
			}
		}
	}

	return rv
}

type EqualToOperator struct {
	BinaryOperator
}

func NewEqualToOperator(left, right Expression) *EqualToOperator {
	return &EqualToOperator{
		BinaryOperator{
			left:  left,
			right: right,
		},
	}
}

func (this *EqualToOperator) EvaluateBoolean(context Context) (bool, error) {
	compare, err := this.BinaryOperator.compare(context)
	if err != nil {
		return false, err
	}
	return compare == 0, nil
}

func (this *EqualToOperator) Evaluate(context Context) (interface{}, error) {
	return this.EvaluateBoolean(context)
}

func (this *EqualToOperator) ConjunctiveNormalForm() BooleanExpression {
	return this
}

func (this *EqualToOperator) NegationNormalForm() BooleanExpression {
	return this
}

func (this *EqualToOperator) DistributeNot() BooleanExpression {
	return NewNotOperator(this)
}

func (this *EqualToOperator) ConvertToBooleanFactors() []BooleanExpression {
	return []BooleanExpression{this}
}

func (this *EqualToOperator) String() string {
	return fmt.Sprintf("%v = %v", this.left, this.right)
}

func (this *EqualToOperator) GetSelectivity(pathStats map[string]stats.PathStatistics) float64 {
	rv := 1.0 / 10.0

	// logic to improve this with pathStats
	if this.IsSargable() {
		path := this.GetSargProperty().Path
		pathStat, ok := pathStats[path]
		if ok {
			sargval, err := this.GetSargValue()
			if err == nil {
				rowsEstimate := RowsWithValue(pathStat, sargval)
				if rowsEstimate != -1 {
					return float64(rowsEstimate / float64(pathStat.Rows))
				}
			}
		}
	}

	return rv
}

type NotEqualToOperator struct {
	BinaryOperator
}

func NewNotEqualToOperator(left, right Expression) *NotEqualToOperator {
	return &NotEqualToOperator{
		BinaryOperator{
			left:  left,
			right: right,
		},
	}
}

func (this *NotEqualToOperator) EvaluateBoolean(context Context) (bool, error) {
	compare, err := this.BinaryOperator.compare(context)
	if err != nil {
		return false, err
	}
	return compare != 0, nil
}

func (this *NotEqualToOperator) Evaluate(context Context) (interface{}, error) {
	return this.EvaluateBoolean(context)
}

func (this *NotEqualToOperator) ConjunctiveNormalForm() BooleanExpression {
	return this
}

func (this *NotEqualToOperator) NegationNormalForm() BooleanExpression {
	return this
}

func (this *NotEqualToOperator) DistributeNot() BooleanExpression {
	return NewNotOperator(this)
}

func (this *NotEqualToOperator) ConvertToBooleanFactors() []BooleanExpression {
	return []BooleanExpression{this}
}

func (this *NotEqualToOperator) String() string {
	return fmt.Sprintf("%v != %v", this.left, this.right)
}

func (this *NotEqualToOperator) GetSelectivity(pathStats map[string]stats.PathStatistics) float64 {
	rv := 9.0 / 10.0

	// logic to improve this with pathStats
	if this.IsSargable() {
		path := this.GetSargProperty().Path
		pathStat, ok := pathStats[path]
		if ok {
			sargval, err := this.GetSargValue()
			if err == nil {
				rowsEstimate := RowsWithoutValue(pathStat, sargval)
				if rowsEstimate != -1 {
					return float64(rowsEstimate / float64(pathStat.Rows))
				}
			}
		}
	}

	return rv
}

// some utility functions for interpretting path statistics

// Most of these estimation functions are very conservative
// and overestimate the number of rows that will be returned
// if they just dont know they return -1

func RowsWithValue(ps stats.PathStatistics, value interface{}) float64 {
	// first see if this is one of the most frequent values
	itemCount := ps.MostFrequentValues.NumItemsWithKey(value)
	if itemCount > 0 {
		return itemCount
	}

	// if not, try the quantiles
	for _, quantile := range ps.Quantiles {
		if CollateJSON(quantile.Start, value) <= 0 && CollateJSON(quantile.End, value) >= 0 {
			return float64(quantile.Count)
		}
	}

	return -1
}

func RowsWithoutValue(ps stats.PathStatistics, value interface{}) float64 {
	// see if we can estimate how many have the value
	rv := RowsWithValue(ps, value)
	if rv != -1 {
		// then subtract from the total
		rv = float64(ps.Rows) - rv
	}
	return rv
}

func RowsBetweenValues(ps stats.PathStatistics, left interface{}, right interface{}) float64 {
	if len(ps.Quantiles) <= 0 {
		return -1
	}
	rv := float64(ps.Rows)
	rv = rv - RowsLessThanValue(ps, left)
	rv = rv - RowsGreaterThanValue(ps, right)
	return rv
}

func RowsLessThanValue(ps stats.PathStatistics, value interface{}) float64 {
	if len(ps.Quantiles) <= 0 {
		return -1
	}

	rv := 0.0
	// now walk through the quantiles
	// as long as the start value of the quantile
	// is less than this value
	// add those rows
	for _, quantile := range ps.Quantiles {
		if CollateJSON(quantile.Start, value) < 0 {
			rv = rv + float64(quantile.Count)
		}
	}
	return rv
}

func RowsGreaterThanValue(ps stats.PathStatistics, value interface{}) float64 {
	if len(ps.Quantiles) <= 0 {
		return -1
	}

	rv := float64(ps.Rows)
	// now walk through the quantiles
	// as long as the end value of the quantile
	// is less than this value
	// subtract those rows
	for _, quantile := range ps.Quantiles {
		if CollateJSON(quantile.End, value) < 0 {
			rv = rv - float64(quantile.Count)
		}
	}
	return rv
}
