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

type BooleanExpression interface {
	Evaluate(context Context) (interface{}, error)
	EvaluateBoolean(context Context) (bool, error)
	ReferencedProperties() []Property
	ConjunctiveNormalForm() BooleanExpression
	NegationNormalForm() BooleanExpression
	DistributeNot() BooleanExpression
	ConvertToBooleanFactors() []BooleanExpression
	IsSargable() bool
	// these 2 methods only return meaningful results if IsSargable() == true
	GetSargProperty() *Property
	GetSargValue() (interface{}, error)
	GetSelectivity(pathStats map[string]stats.PathStatistics) float64
}

type BooleanOperator struct {
	operands []BooleanExpression
}

type AndOperator struct {
	BooleanOperator
}

func NewAndOperator(operands []BooleanExpression) *AndOperator {
	return &AndOperator{
		BooleanOperator{
			operands: operands,
		},
	}
}

func (this *AndOperator) EvaluateBoolean(context Context) (bool, error) {
	for _, operand := range this.operands {
		operandVal, err := operand.EvaluateBoolean(context)
		if err != nil {
			return false, err
		} else if !operandVal {
			return false, nil
		}
	}
	return true, nil
}

func (this *AndOperator) Evaluate(context Context) (interface{}, error) {
	return this.EvaluateBoolean(context)
}

func (this *AndOperator) ConjunctiveNormalForm() BooleanExpression {
	resultOperands := make([]BooleanExpression, 0, 0)
	for _, operand := range this.operands {
		// first get all the children into CNF
		operandCNF := operand.ConjunctiveNormalForm()
		// then if the child is also an AND collapse it into self
		switch operandCNF := operandCNF.(type) {
		case *AndOperator:
			for _, inner := range operandCNF.operands {
				resultOperands = append(resultOperands, inner)
			}
		default:
			resultOperands = append(resultOperands, operandCNF)
		}
	}
	return &AndOperator{
		BooleanOperator{
			operands: resultOperands,
		},
	}
}

func (this *AndOperator) NegationNormalForm() BooleanExpression {
	resultOperands := make([]BooleanExpression, 0, 0)
	for _, operand := range this.operands {
		// get all the children into NNF
		operandNNF := operand.NegationNormalForm()
		resultOperands = append(resultOperands, operandNNF)
	}
	return &AndOperator{
		BooleanOperator{
			operands: resultOperands,
		},
	}
}

func (this *AndOperator) DistributeNot() BooleanExpression {
	resultOperands := make([]BooleanExpression, 0, 0)
	for _, operand := range this.operands {
		operandAfter := operand.DistributeNot()
		resultOperands = append(resultOperands, operandAfter)
	}
	return &OrOperator{
		BooleanOperator{
			operands: resultOperands,
		},
	}
}

func (this *AndOperator) ConvertToBooleanFactors() []BooleanExpression {
	return this.operands
}

func (this *AndOperator) String() string {
	return fmt.Sprintf("AND %v", this.operands)
}

func (this *AndOperator) ReferencedProperties() []Property {
	rv := make([]Property, 0, 0)
	for _, v := range this.operands {
		rp := v.ReferencedProperties()
		for _, rpv := range rp {
			rv = append(rv, rpv)
		}
	}
	return rv
}

func (this *AndOperator) IsSargable() bool {
	return false
}

func (this *AndOperator) GetSargProperty() *Property {
	return nil
}

func (this *AndOperator) GetSargValue() (interface{}, error) {
	return nil, nil
}

func (this *AndOperator) GetSelectivity(pathStats map[string]stats.PathStatistics) float64 {
	// assume their selectivities are independent
	rv := 1.0
	for _, operand := range this.operands {
		operSelectivity := operand.GetSelectivity(pathStats)
		rv = rv * operSelectivity
	}
	return rv
}

type OrOperator struct {
	BooleanOperator
}

func NewOrOperator(operands []BooleanExpression) *OrOperator {
	return &OrOperator{
		BooleanOperator{
			operands: operands,
		},
	}
}

func (this *OrOperator) EvaluateBoolean(context Context) (bool, error) {
	for _, operand := range this.operands {
		operandVal, err := operand.EvaluateBoolean(context)
		if err != nil {
			return false, err
		} else if operandVal {
			return true, nil
		}
	}
	return false, nil
}

func (this *OrOperator) Evaluate(context Context) (interface{}, error) {
	return this.EvaluateBoolean(context)
}

func (this *OrOperator) ConjunctiveNormalForm() BooleanExpression {
	andOperands := make([]*AndOperator, 0, 0)
	otherOperands := make([]BooleanExpression, 0, 0)
	for _, operand := range this.operands {
		// first get all the children into CNF
		operandCNF := operand.ConjunctiveNormalForm()
		// then if the child is also an OR collapse it into self
		switch operandCNF := operandCNF.(type) {
		case *AndOperator:
			andOperands = append(andOperands, operandCNF)
		case *OrOperator:
			for _, inner := range operandCNF.operands {
				otherOperands = append(otherOperands, inner)
			}
		default:
			otherOperands = append(otherOperands, operandCNF)
		}
	}

	resultPieces := make([]*OrOperator, 0, 0)
	// start with a single OrOperator with the other operands
	resultPieces = append(resultPieces, NewOrOperator(otherOperands))

	// now we should have two lists, the ANDs we need to distribute over and the other OR items
	if len(andOperands) > 0 {

		for _, andOperand := range andOperands {
			newResultPieces := make([]*OrOperator, 0, 0)
			for _, innerAnd := range andOperand.operands {
				for _, resultPiece := range resultPieces {
					newResultOperands := make([]BooleanExpression, len(resultPiece.operands))
					copy(newResultOperands, resultPiece.operands)
					newResultOperands = append(newResultOperands, innerAnd)
					newResultPieces = append(newResultPieces, NewOrOperator(newResultOperands))
				}
			}
			resultPieces = newResultPieces
		}
	} else {
		return NewOrOperator(otherOperands)
	}

	// ugly conversion code that go won't do impliclitly
	booleanExpressions := make([]BooleanExpression, len(resultPieces))
	for i, v := range resultPieces {
		booleanExpressions[i] = BooleanExpression(v)
	}
	return NewAndOperator(booleanExpressions)
}

func (this *OrOperator) NegationNormalForm() BooleanExpression {
	resultOperands := make([]BooleanExpression, 0, 0)
	for _, operand := range this.operands {
		// get all the children into NNF
		operandNNF := operand.NegationNormalForm()
		resultOperands = append(resultOperands, operandNNF)
	}
	return &OrOperator{
		BooleanOperator{
			operands: resultOperands,
		},
	}
}

func (this *OrOperator) DistributeNot() BooleanExpression {
	resultOperands := make([]BooleanExpression, 0, 0)
	for _, operand := range this.operands {
		operandAfter := operand.DistributeNot()
		resultOperands = append(resultOperands, operandAfter)
	}
	return &AndOperator{
		BooleanOperator{
			operands: resultOperands,
		},
	}
}

func (this *OrOperator) ConvertToBooleanFactors() []BooleanExpression {
	return []BooleanExpression{this}
}

func (this *OrOperator) String() string {
	return fmt.Sprintf("OR %v", this.operands)
}

func (this *OrOperator) ReferencedProperties() []Property {
	rv := make([]Property, 0, 0)
	for _, v := range this.operands {
		rp := v.ReferencedProperties()
		for _, rpv := range rp {
			rv = append(rv, rpv)
		}
	}
	return rv
}

func (this *OrOperator) IsSargable() bool {
	return false
}

func (this *OrOperator) GetSargProperty() *Property {
	return nil
}

func (this *OrOperator) GetSargValue() (interface{}, error) {
	return nil, nil
}

func (this *OrOperator) GetSelectivity(pathStats map[string]stats.PathStatistics) float64 {
	rv := 0.0
	for _, operand := range this.operands {
		operSelectivity := operand.GetSelectivity(pathStats)
		left := rv + operSelectivity
		right := rv * operSelectivity
		rv = left - right
	}
	return rv
}

type NotOperator struct {
	operand BooleanExpression
}

func NewNotOperator(operand BooleanExpression) *NotOperator {
	return &NotOperator{
		operand: operand,
	}
}

func (this *NotOperator) EvaluateBoolean(context Context) (bool, error) {
	ov, err := this.operand.EvaluateBoolean(context)
	if err != nil {
		return false, err
	}

	return !ov, nil
}

func (this *NotOperator) Evaluate(context Context) (interface{}, error) {
	return this.EvaluateBoolean(context)
}

func (this *NotOperator) ConjunctiveNormalForm() BooleanExpression {
	// get the operand into CNF
	resultOperand := this.operand.ConjunctiveNormalForm()
	return &NotOperator{
		operand: resultOperand,
	}
}

func (this *NotOperator) NegationNormalForm() BooleanExpression {
	// put the operand into negation normal form
	resultOperand := this.operand.NegationNormalForm()
	// now distribute the not through
	return resultOperand.DistributeNot()
}

func (this *NotOperator) DistributeNot() BooleanExpression {
	// eliminate double negative
	return this.operand
}

func (this *NotOperator) ConvertToBooleanFactors() []BooleanExpression {
	return []BooleanExpression{this}
}

func (this *NotOperator) ReferencedProperties() []Property {
	return this.operand.ReferencedProperties()
}

func (this *NotOperator) IsSargable() bool {
	return false
}

func (this *NotOperator) GetSargProperty() *Property {
	return nil
}

func (this *NotOperator) GetSargValue() (interface{}, error) {
	return nil, nil
}

func (this *NotOperator) GetSelectivity(pathStats map[string]stats.PathStatistics) float64 {
	rv := 1.0 - this.operand.GetSelectivity(pathStats)
	return rv
}
