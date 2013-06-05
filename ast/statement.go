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
	"log"
)

const STATEMENT_TYPE_SELECT = "SELECT"

type Statement interface {
	GetType() string
	GetFrom() []DataSource
	GetWhere() BooleanExpression
	SetFrom([]DataSource)
	GetSelect() Expression
	GetOrder() []OrderedExpression
	GetLimit() int
	GetOffset() int
}

type SelectStatement struct {
	From   []DataSource
	Where  BooleanExpression
	Select Expression
	Order  []OrderedExpression
	Limit  int
	Offset int
}

func (this *SelectStatement) GetType() string {
	return STATEMENT_TYPE_SELECT
}

func (this *SelectStatement) GetWhere() BooleanExpression {
	return this.Where
}

func (this *SelectStatement) GetFrom() []DataSource {
	return this.From
}

func (this *SelectStatement) SetFrom(from []DataSource) {
	this.From = from
}

func (this *SelectStatement) GetSelect() Expression {
	return this.Select
}

func (this *SelectStatement) GetOrder() []OrderedExpression {
	return this.Order
}

func (this *SelectStatement) GetLimit() int {
	return this.Limit
}

func (this *SelectStatement) GetOffset() int {
	return this.Offset
}

func (this *SelectStatement) String() string {
	return fmt.Sprintf("SELECT %v WHERE %v ORDER BY %v LIMIT %v OFFSET %v", this.Select, this.Where, this.Order, this.Limit, this.Offset)
}

func NewSelectStatement() *SelectStatement {
	return &SelectStatement{
		From:   make([]DataSource, 0),
		Where:  NewLiteralBool(true),
		Order:  make([]OrderedExpression, 0),
		Limit:  -1,
		Offset: 0,
	}
}

func NewStatementFromJSONRequestToBucket(bucket string, body map[string]interface{}) (Statement, error) {
	statementJSON, err := validateTopLevel(body)
	if err != nil {
		return nil, err
	}

	statement, err := parseStatementJSON(statementJSON)
	if err != nil {
		return nil, err
	}

	statement.SetFrom([]DataSource{NewNamedDataSource(bucket)})

	return statement, nil
}

func parseStatementJSON(statementJSON map[string]interface{}) (Statement, error) {
	statementType, ok := statementJSON["type"]
	if !ok {
		return nil, fmt.Errorf("Statement must specify type")
	}
	switch statementType {
	case "select":
		return parseSelectStatementJSON(statementJSON)
	}
	return nil, fmt.Errorf("Statement type %v is not supported", statementType)
}

func parseSelectStatementJSON(statementJSON map[string]interface{}) (Statement, error) {
	selectStatement := NewSelectStatement()
	whereJSON, ok := statementJSON["where"]
	if ok {
		switch whereJSON := whereJSON.(type) {
		case map[string]interface{}:
			// there is a where clause
			whereClause, err := parseBooleanExpression(whereJSON)
			if err != nil {
				return nil, err
			}
			selectStatement.Where = whereClause
		default:
			return nil, fmt.Errorf("where element must be an object")
		}
	}

	selectJSON, ok := statementJSON["select"]
	if ok {
		switch selectJSON := selectJSON.(type) {
		case map[string]interface{}:
			// there is a where clause
			selectClause, err := parseExpression(selectJSON)
			if err != nil {
				return nil, err
			}
			selectStatement.Select = selectClause
		default:
			return nil, fmt.Errorf("select element must be an object")
		}
	}

	orderJSON, ok := statementJSON["order"]
	if ok {
		switch orderJSON := orderJSON.(type) {
		case []interface{}:
			// there is an oder by clause
			orderClause, err := parseOrderBy(orderJSON)
			if err != nil {
				return nil, err
			}
			selectStatement.Order = orderClause
		default:
			return nil, fmt.Errorf("order element must be an array")
		}
	}
	log.Printf("looking for limit")
	limitJSON, ok := statementJSON["limit"]
	if ok {
		log.Printf("see limit")
		switch limitJSON := limitJSON.(type) {
		case float64:
			// FIXME probably shouldnt allow LIMIT 3.2 to actually work
			selectStatement.Limit = int(limitJSON)
		default:
			return nil, fmt.Errorf("limit must be an integer")
		}

		// offset is only valid if there is a limit (per SQL)
		offsetJSON, ok := statementJSON["offset"]
		if ok {
			switch offsetJSON := offsetJSON.(type) {
			case float64:
				// FIXME probably shouldnt allow LIMIT 3.2 to actually work
				selectStatement.Offset = int(offsetJSON)
			default:
				return nil, fmt.Errorf("offset must be an integer")
			}
		}
	}

	return selectStatement, nil
}

func parseOrderBy(orderJSON []interface{}) ([]OrderedExpression, error) {
	rv := make([]OrderedExpression, 0, len(orderJSON))
	for _, v := range orderJSON {
		switch v := v.(type) {
		case map[string]interface{}:
			se, err := parseSortExpression(v)
			if err != nil {
				return nil, err
			}
			rv = append(rv, se)
		default:
			return nil, fmt.Errorf("members of sort array must be sort objects")
		}
	}
	return rv, nil
}

func parseSortExpression(sortJSON map[string]interface{}) (OrderedExpression, error) {
	ascending := true
	ascendingJSON, ok := sortJSON["ascending"]
	if ok {
		switch ascendingJSON := ascendingJSON.(type) {
		case bool:
			ascending = ascendingJSON
		default:
			return nil, fmt.Errorf("ascending must be type boolean")
		}
	}
	exprJSON, ok := sortJSON["expr"]
	if !ok {
		return nil, fmt.Errorf("sort object is imissing expression")
	}
	switch exprJSON := exprJSON.(type) {
	case map[string]interface{}:
		expr, err := parseExpression(exprJSON)
		if err != nil {
			return nil, err
		}
		return NewSortExpression(expr, ascending), nil
	}

	return nil, fmt.Errorf("sort expression must be type object")
}

func parseExpression(expressionJSON map[string]interface{}) (Expression, error) {
	expressionType, ok := expressionJSON["type"]
	if !ok {
		return nil, fmt.Errorf("expression is missing type")
	}

	switch expressionType {
	case "compare", "and", "or", "not":
		return parseBooleanExpression(expressionJSON)
	case "literal":
		return parseLiteral(expressionJSON)
	case "property":
		return parseProperty(expressionJSON)
	case "arithmetic":
		return parseArithmetic(expressionJSON)
	}

	return nil, fmt.Errorf("Unrecognized expression type %v", expressionType)
}

func parseArithmetic(expressionJSON map[string]interface{}) (Expression, error) {
	operator, ok := expressionJSON["operator"]
	if !ok {
		return nil, fmt.Errorf("arithmetic must specify operator")
	}
	left, right, err := parseBinaryOperatorArguments(expressionJSON)
	if err != nil {
		return nil, err
	}
	switch operator {
	case "plus":
		return NewPlusOperator(left, right), nil
	case "minus":
		return NewSubtractOperator(left, right), nil
	case "mult":
		return NewMultiplyOperator(left, right), nil
	case "div":
		return NewDivideOperator(left, right), nil
	}
	return nil, fmt.Errorf("Unsupported arithmetic operator %v", operator)
}

func parseProperty(expressionJSON map[string]interface{}) (Expression, error) {
	path, ok := expressionJSON["path"]
	if !ok {
		return nil, fmt.Errorf("property must contain path")
	}
	switch path := path.(type) {
	case string:
		return NewProperty(path), nil
	}
	return nil, fmt.Errorf("property path must be a string")
}

func parseLiteral(expressionJSON map[string]interface{}) (Expression, error) {
	value, ok := expressionJSON["value"]
	if !ok {
		return nil, fmt.Errorf("literal requires value")
	}
	switch value := value.(type) {
	case bool:
		return NewLiteralBool(value), nil
	case float64:
		return NewLiteralNumber(value), nil
	case string:
		return NewLiteralString(value), nil
	case []interface{}:
		exprArr := make([]Expression, 0, len(value))
		for _, childExpressionJSON := range value {
			switch childExpressionJSON := childExpressionJSON.(type) {
			case map[string]interface{}:
				childexpr, err := parseExpression(childExpressionJSON)
				if err != nil {
					return nil, err
				}
				exprArr = append(exprArr, childexpr)
			default:
				return nil, fmt.Errorf("members of literal array must be expression objects")
			}

		}
		return NewLiteralArray(exprArr), nil
	case map[string]interface{}:
		exprObj := make(map[string]Expression)
		for k, v := range value {
			switch v := v.(type) {
			case map[string]interface{}:
				childexpr, err := parseExpression(v)
				if err != nil {
					return nil, err
				}
				exprObj[k] = childexpr
			default:
				return nil, fmt.Errorf("values of literal object must be expression objects")
			}
		}
		return NewLiteralObject(exprObj), nil
	}
	return nil, fmt.Errorf("Unexpected type %T", value)
}

func parseBooleanExpression(expressionJSON map[string]interface{}) (BooleanExpression, error) {
	expressionType, ok := expressionJSON["type"]
	if !ok {
		return nil, fmt.Errorf("expression is missing type")
	}

	switch expressionType {
	case "compare":
		left, right, err := parseBinaryOperatorArguments(expressionJSON)
		if err != nil {
			return nil, err
		}
		return parseCompareExpression(left, right, expressionJSON)
	case "and", "or":
		left, right, err := parseBinaryBooleanOperatorArguments(expressionJSON)
		if err != nil {
			return nil, err
		}
		if expressionType == "and" {
			return NewAndOperator([]BooleanExpression{left, right}), nil
		} else {
			return NewOrOperator([]BooleanExpression{left, right}), nil
		}
	case "not":
		operand, err := parseUnaryBooleanOperatorArguments(expressionJSON)
		if err != nil {
			return nil, err
		}
		return NewNotOperator(operand), nil
	}
	return nil, fmt.Errorf("Unrecognized expression type %v", expressionType)
}

func parseCompareExpression(left Expression, right Expression, expressionJSON map[string]interface{}) (BooleanExpression, error) {
	operator, ok := expressionJSON["operator"]
	if ok {
		switch operator {
		case "gt":
			return NewGreaterThanOperator(left, right), nil
		case "gte":
			return NewGreaterThanOrEqualOperator(left, right), nil
		case "lt":
			return NewLessThanOperator(left, right), nil
		case "lte":
			return NewLessThanOrEqualOperator(left, right), nil
		case "eq":
			return NewEqualToOperator(left, right), nil
		case "neq":
			return NewNotEqualToOperator(left, right), nil
		}
		return nil, fmt.Errorf("unsupported comparison operator %v", operator)
	}
	return nil, fmt.Errorf("comparison element must specify operator")
}

func parseUnaryBooleanOperatorArguments(expressionJSON map[string]interface{}) (BooleanExpression, error) {
	operandJSON, ok := expressionJSON["operand"]
	if !ok {
		return nil, fmt.Errorf("unary operator is missing element operand")
	}
	switch operandJSON := operandJSON.(type) {
	case map[string]interface{}:

		operand, err := parseBooleanExpression(operandJSON)
		if err != nil {
			return nil, err
		}
		return operand, nil
	}
	return nil, fmt.Errorf("unary operator element operand must be an object")
}

func parseBinaryBooleanOperatorArguments(expressionJSON map[string]interface{}) (BooleanExpression, BooleanExpression, error) {
	leftJSON, ok := expressionJSON["left"]
	if !ok {
		return nil, nil, fmt.Errorf("binary operator is missing element left")
	}
	rightJSON, ok := expressionJSON["right"]
	if !ok {
		return nil, nil, fmt.Errorf("binary operator is missing element right")
	}
	switch leftJSON := leftJSON.(type) {
	case map[string]interface{}:
		switch rightJSON := rightJSON.(type) {
		case map[string]interface{}:
			left, err := parseBooleanExpression(leftJSON)
			if err != nil {
				return nil, nil, err
			}
			right, err := parseBooleanExpression(rightJSON)
			if err != nil {
				return nil, nil, err
			}
			return left, right, nil
		default:
			return nil, nil, fmt.Errorf("binary operator element right must be object")
		}
	}
	return nil, nil, fmt.Errorf("binary operator element left must be an object")
}

func parseBinaryOperatorArguments(expressionJSON map[string]interface{}) (Expression, Expression, error) {
	leftJSON, ok := expressionJSON["left"]
	if !ok {
		return nil, nil, fmt.Errorf("binary operator is missing element left")
	}
	rightJSON, ok := expressionJSON["right"]
	if !ok {
		return nil, nil, fmt.Errorf("binary operator is missing element right")
	}
	switch leftJSON := leftJSON.(type) {
	case map[string]interface{}:
		switch rightJSON := rightJSON.(type) {
		case map[string]interface{}:
			left, err := parseExpression(leftJSON)
			if err != nil {
				return nil, nil, err
			}
			right, err := parseExpression(rightJSON)
			if err != nil {
				return nil, nil, err
			}
			return left, right, nil
		default:
			return nil, nil, fmt.Errorf("binary operator element right must be object")
		}
	}
	return nil, nil, fmt.Errorf("binary operator element left must be an object")
}

func validateTopLevel(body map[string]interface{}) (map[string]interface{}, error) {
	requestType, ok := body["type"]
	if !ok || requestType != "cbqast" {
		return nil, fmt.Errorf("Unrecognized request type %v", requestType)
	}
	requestVersion, ok := body["version"]
	if !ok || requestVersion != "1" {
		return nil, fmt.Errorf("Unrecognized version %v", requestVersion)
	}
	statementJSON, ok := body["statement"]
	if !ok {
		return nil, fmt.Errorf("Missing required element statement at top-level")
	}
	switch statementJSON := statementJSON.(type) {
	case map[string]interface{}:
		return statementJSON, nil
	}
	return nil, fmt.Errorf("statement element must be an object")
}
