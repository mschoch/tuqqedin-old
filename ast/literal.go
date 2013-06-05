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

type LiteralNull struct {
}

func NewLiteralNull() *LiteralNull {
	return &LiteralNull{}
}

func (this *LiteralNull) Evaluate(context Context) (interface{}, error) {
	return nil, nil
}

func (this *LiteralNull) String() string {
	return fmt.Sprintf("null")
}

func (this *LiteralNull) ReferencedProperties() []Property {
	return []Property{}
}

type LiteralBool struct {
	Value bool
}

func NewLiteralBool(Value bool) *LiteralBool {
	return &LiteralBool{
		Value: Value,
	}
}

func (this *LiteralBool) EvaluateBoolean(context Context) (bool, error) {
	return this.Value, nil
}

func (this *LiteralBool) Evaluate(context Context) (interface{}, error) {
	return this.EvaluateBoolean(context)
}

func (this *LiteralBool) ConjunctiveNormalForm() BooleanExpression {
	return this
}

func (this *LiteralBool) NegationNormalForm() BooleanExpression {
	return this
}

func (this *LiteralBool) DistributeNot() BooleanExpression {
	return NewNotOperator(this)
}

func (this *LiteralBool) ConvertToBooleanFactors() []BooleanExpression {
	return []BooleanExpression{this}
}

func (this *LiteralBool) String() string {
	return fmt.Sprintf("%v", this.Value)
}

func (this *LiteralBool) ReferencedProperties() []Property {
	return []Property{}
}

func (this *LiteralBool) IsSargable() bool {
	return false
}

func (this *LiteralBool) GetSargProperty() *Property {
	return nil
}

func (this *LiteralBool) GetSargValue() (interface{}, error) {
	return nil, nil
}

func (this *LiteralBool) GetSelectivity(pathStats map[string]stats.PathStatistics) float64 {
	if this.Value {
		return 1.0
	}
	return 0.0
}

type LiteralNumber struct {
	Value float64
}

func NewLiteralNumber(Value float64) *LiteralNumber {
	return &LiteralNumber{
		Value: Value,
	}
}

func (this *LiteralNumber) Evaluate(context Context) (interface{}, error) {
	return this.Value, nil
}

func (this *LiteralNumber) String() string {
	return fmt.Sprintf("%v", this.Value)
}

func (this *LiteralNumber) ReferencedProperties() []Property {
	return []Property{}
}

type LiteralString struct {
	Value string
}

func NewLiteralString(Value string) *LiteralString {
	return &LiteralString{
		Value: Value,
	}
}

func (this *LiteralString) Evaluate(context Context) (interface{}, error) {
	return this.Value, nil
}

func (this *LiteralString) String() string {
	return fmt.Sprintf("%v", this.Value)
}

func (this *LiteralString) ReferencedProperties() []Property {
	return []Property{}
}

type LiteralArray struct {
	Value []Expression
}

func NewLiteralArray(Value []Expression) *LiteralArray {
	return &LiteralArray{
		Value: Value,
	}
}

func (this *LiteralArray) Evaluate(context Context) (interface{}, error) {
	rv := make([]interface{}, 0, len(this.Value))
	for _, v := range this.Value {
		ev, err := v.Evaluate(context)
		if err != nil {
			return nil, err
		}
		rv = append(rv, ev)
	}
	return rv, nil
}

func (this *LiteralArray) String() string {
	return fmt.Sprintf("%v", this.Value)
}

func (this *LiteralArray) ReferencedProperties() []Property {
	rv := make([]Property, 0, 0)
	for _, v := range this.Value {
		rp := v.ReferencedProperties()
		for _, rpv := range rp {
			rv = append(rv, rpv)
		}
	}
	return rv
}

type LiteralObject struct {
	Value map[string]Expression
}

func NewLiteralObject(Value map[string]Expression) *LiteralObject {
	return &LiteralObject{
		Value: Value,
	}
}

func (this *LiteralObject) Evaluate(context Context) (interface{}, error) {
	rv := make(map[string]interface{}, len(this.Value))
	for k, v := range this.Value {
		ev, err := v.Evaluate(context)
		if err != nil {
			return nil, err
		}
		rv[k] = ev
	}
	return rv, nil
}

func (this *LiteralObject) String() string {
	return fmt.Sprintf("%v", this.Value)
}

func (this *LiteralObject) ReferencedProperties() []Property {
	rv := make([]Property, 0, 0)
	for _, v := range this.Value {
		rp := v.ReferencedProperties()
		for _, rpv := range rp {
			rv = append(rv, rpv)
		}
	}
	return rv
}
