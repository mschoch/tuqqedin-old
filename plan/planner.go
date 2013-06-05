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
	"log"

	"github.com/couchbaselabs/tuqqedin/ast"
	"github.com/couchbaselabs/tuqqedin/datasource"
)

type Planner interface {
	Plan(statment ast.Statement) ([]Operator, error)
}

type CouchbasePlanner struct {
	dataSourceManager datasource.DataSourceManager
}

func NewCouchbasePlanner(dataSourceManager datasource.DataSourceManager) *CouchbasePlanner {
	return &CouchbasePlanner{
		dataSourceManager: dataSourceManager,
	}
}

func (this *CouchbasePlanner) Plan(statement ast.Statement) ([]Operator, error) {
	rv := make([]Operator, 0)

	switch statement.GetType() {
	case ast.STATEMENT_TYPE_SELECT:
		if len(statement.GetFrom()) != 1 {
			panic("Only 1 data source is currently supported")
		}

		namedDataSource := statement.GetFrom()[0]
		couchbaseDataSource, err := this.dataSourceManager.GetDataSource(namedDataSource.GetName())
		if err != nil {
			return nil, err
		}

		// put the WHERE clause in NNF
		nnf := statement.GetWhere().NegationNormalForm()
		cnf := nnf.ConjunctiveNormalForm()
		booleanFactors := cnf.ConvertToBooleanFactors()

		// look at each access path the datasource
		// and try to create a plan using it
		for _, accessPath := range couchbaseDataSource.AccessPaths() {
			if accessPath.ReturnsAll() || accessPath.Matches(booleanFactors) {
				var currentOperator Operator
				//sourceOperator, remainingFactors := accessPath.BuildOperator(booleanFactors)
				// NOTE here we don't currently use the "remaining factors"
				// because we if we do a fetch, we have to recheck anyway
				// FIXME if the query is covered by the index it can be avoided
				currentOperator, _ = buildOperatorForAccessPath(accessPath, booleanFactors)
				// FIXME need to check select clause to see if we need fetch
				currentOperator = NewFetch(currentOperator, couchbaseDataSource)
				currentOperator = NewFilter(currentOperator, booleanFactors)

				order := statement.GetOrder()
				if len(order) > 0 {
					currentOperator = NewOrder(currentOperator, order)
				}

				offset := statement.GetOffset()
				if offset > 0 {
					currentOperator = NewOffset(currentOperator, offset)
				}

				limit := statement.GetLimit()
				if limit >= 0 {
					currentOperator = NewLimit(currentOperator, limit)
				}

				projection := statement.GetSelect()
				currentOperator = NewProject(currentOperator, projection)

				// add operator as plan
				rv = append(rv, currentOperator)
			} else {
				log.Printf("cannot use ap: %v", accessPath)
			}
		}

	}

	return rv, nil
}

func buildOperatorForAccessPath(accessPath datasource.AccessPath, booleanFactors []ast.BooleanExpression) (Operator, []ast.BooleanExpression) {
	switch accessPath := accessPath.(type) {
	case *datasource.CouchbaseAllDocsAccessPath:
		return NewAllDocsScanner(accessPath), booleanFactors
	case *datasource.CouchbaseViewAccessPath:
		unsupportedFactors := make([]ast.BooleanExpression, 0, 0)
		viewScanner := NewViewScanner(accessPath)

		for _, booleanFactor := range booleanFactors {
			if !viewScanner.AddBooleanFactor(booleanFactor) {
				unsupportedFactors = append(unsupportedFactors, booleanFactor)
			}
		}

		return viewScanner, unsupportedFactors
	}

	panic("diedie")
}
