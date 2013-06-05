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
)

type OrderedExpression interface {
	Expression() Expression
	Order() bool
}

type SortExpression struct {
	Expr      Expression
	Ascending bool
}

func NewSortExpression(expr Expression, asc bool) *SortExpression {
	return &SortExpression{
		Expr:      expr,
		Ascending: asc,
	}
}

func (this *SortExpression) Expression() Expression {
	return this.Expr
}

func (this *SortExpression) Order() bool {
	return this.Ascending
}

func (this *SortExpression) String() string {
	if this.Ascending {
		return fmt.Sprintf("%v ASC", this.Expr)
	}
	return fmt.Sprintf("%v DESC", this.Expr)
}
