//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package ast

type PlusOperator struct {
	left  Expression
	right Expression
}

func NewPlusOperator(left, right Expression) *PlusOperator {
	return &PlusOperator{
		left:  left,
		right: right,
	}
}

func (this *PlusOperator) Evaluate(context Context) (interface{}, error) {
	lv, err := this.left.Evaluate(context)
	if err != nil {
		return nil, err
	}
	rv, err := this.right.Evaluate(context)
	if err != nil {
		return nil, err
	}

	switch lv := lv.(type) {
	case string:
		switch rv := rv.(type) {
		case string:
			// if both values are strings append
			return lv + rv, nil
		default:
			return nil, nil
		}
	case float64:
		switch rv := rv.(type) {
		case float64:
			// if both values are numeric add
			return lv + rv, nil
		default:
			return nil, nil
		}
	default:
		return nil, nil
	}

	return nil, nil
}

func (this *PlusOperator) ReferencedProperties() []Property {
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

type SubtractOperator struct {
	left  Expression
	right Expression
}

func NewSubtractOperator(left, right Expression) *SubtractOperator {
	return &SubtractOperator{
		left:  left,
		right: right,
	}
}

func (this *SubtractOperator) Evaluate(context Context) (interface{}, error) {
	lv, err := this.left.Evaluate(context)
	if err != nil {
		return nil, err
	}
	rv, err := this.right.Evaluate(context)
	if err != nil {
		return nil, err
	}

	switch lv := lv.(type) {
	case float64:
		switch rv := rv.(type) {
		case float64:
			// if both values are numeric subtract
			return lv - rv, nil
		default:
			return nil, nil
		}
	default:
		return nil, nil
	}

	return nil, nil
}

func (this *SubtractOperator) ReferencedProperties() []Property {
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

type MultiplyOperator struct {
	left  Expression
	right Expression
}

func NewMultiplyOperator(left, right Expression) *MultiplyOperator {
	return &MultiplyOperator{
		left:  left,
		right: right,
	}
}

func (this *MultiplyOperator) Evaluate(context Context) (interface{}, error) {
	lv, err := this.left.Evaluate(context)
	if err != nil {
		return nil, err
	}
	rv, err := this.right.Evaluate(context)
	if err != nil {
		return nil, err
	}

	switch lv := lv.(type) {
	case float64:
		switch rv := rv.(type) {
		case float64:
			// if both values are numeric multiply
			return lv * rv, nil
		default:
			return nil, nil
		}
	default:
		return nil, nil
	}

	return nil, nil
}

func (this *MultiplyOperator) ReferencedProperties() []Property {
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

type DivideOperator struct {
	left  Expression
	right Expression
}

func NewDivideOperator(left, right Expression) *DivideOperator {
	return &DivideOperator{
		left:  left,
		right: right,
	}
}

func (this *DivideOperator) Evaluate(context Context) (interface{}, error) {
	lv, err := this.left.Evaluate(context)
	if err != nil {
		return nil, err
	}
	rv, err := this.right.Evaluate(context)
	if err != nil {
		return nil, err
	}

	switch lv := lv.(type) {
	case float64:
		switch rv := rv.(type) {
		case float64:
			// if both values are numeric divide
			return lv / rv, nil
		default:
			return nil, nil
		}
	default:
		return nil, nil
	}

	return nil, nil
}

func (this *DivideOperator) ReferencedProperties() []Property {
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
