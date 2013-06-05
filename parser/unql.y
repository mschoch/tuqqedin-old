%{
package parser
import "fmt"
import "log"
import "github.com/couchbaselabs/tuqqedin/ast"

func logDebugGrammar(format string, v ...interface{}) {
    if DebugGrammar && len(v) > 0 {
        log.Printf("DEBUG GRAMMAR " + format, v)
    } else if DebugGrammar {
        log.Printf("DEBUG GRAMMAR " + format)
    }
}
%}

%union { 
s string 
n int
f float64}

%token INT REAL STRING TRUE FALSE NULL
%token IDENTIFIER DOT
%token LBRACKET RBRACKET COMMA LBRACE RBRACE COLON
%token PLUS MINUS MULT DIV
%token SELECT WHERE ORDER BY ASC DESC
%token OFFSET LIMIT
%token LPAREN RPAREN
%token AND OR NOT
%token LT LTE GT GTE EQ NE 
%left OR
%left AND 
%left EQ LT LTE GT GTE NE
%left PLUS MINUS MULT DIV MOD
%right NOT
%right QUESTION
%%

input: select_stmt { 
	logDebugGrammar("INPUT") 
}
;

select_stmt:	select_compound select_order select_limit_offset {
	logDebugGrammar("SELECT_STMT")
}
;

select_compound:    select_core { 
	logDebugGrammar("SELECT_COMPOUND") 
}
;

select_core:    select_select select_where { 
	logDebugGrammar("SELECT_CORE")
}
;

select_select:  select_select_head select_select_tail {
	logDebugGrammar("SELECT_SELECT")
}
;

select_select_head:  SELECT { 
	logDebugGrammar("SELECT_SELECT_HEAD")
	if parsingStatement == nil {
		parsingStatement = ast.NewSelectStatement()
	}
}
;

select_select_tail:		MULT { 
	logDebugGrammar("SELECT SELECT TAIL - STAR")
} 
|	expression { 
	logDebugGrammar("SELECT SELECT TAIL - EXPR")
	select_part := parsingStack.Pop().(ast.Expression)
	switch parsingStatement := parsingStatement.(type) {
	case *ast.SelectStatement:
		parsingStatement.Select = select_part
		logDebugGrammar("set a select")
	default:
		logDebugGrammar("This statement does not support SELECT")
	}
}
;

select_where:   
/* empty */ { 
	logDebugGrammar("SELECT WHERE - EMPTY")
}
|
WHERE expression {
	logDebugGrammar("SELECT WHERE - EXPR")
	where_part := parsingStack.Pop().(ast.BooleanExpression)
	switch parsingStatement := parsingStatement.(type) {
	case *ast.SelectStatement:
		parsingStatement.Where = where_part
	default:
		logDebugGrammar("This statement does not support WHERE")
	}
};

select_order:   
/* empty */
|
ORDER BY sorting_list {
	
}
;

sorting_list:
sorting_single {
	
}
|
sorting_single COMMA sorting_list {
	
};

sorting_single:
expression { 
	thisExpression := ast.NewSortExpression(parsingStack.Pop().(ast.Expression), true)
	switch parsingStatement := parsingStatement.(type) {
	case *ast.SelectStatement:
		parsingStatement.Order = append(parsingStatement.Order, thisExpression)
	default:
		logDebugGrammar("This statement does not support ORDER BY")
	}
}
|
expression ASC { 
	thisExpression := ast.NewSortExpression(parsingStack.Pop().(ast.Expression), true)
	switch parsingStatement := parsingStatement.(type) {
	case *ast.SelectStatement:
		parsingStatement.Order = append(parsingStatement.Order, thisExpression)
	default:
		logDebugGrammar("This statement does not support ORDER BY")
	}
}
|
expression DESC { 
	thisExpression := ast.NewSortExpression(parsingStack.Pop().(ast.Expression), false)
	switch parsingStatement := parsingStatement.(type) {
	case *ast.SelectStatement:
		parsingStatement.Order = append(parsingStatement.Order, thisExpression)
	default:
		logDebugGrammar("This statement does not support ORDER BY")
	}
};

select_limit_offset:
/* empty */ {
	
}
|
select_limit {
	
}
|
select_limit select_offset {
	
}
;

select_limit:
LIMIT expression {
	thisExpression := parsingStack.Pop()
	switch thisExpression := thisExpression.(type) {
	case *ast.LiteralNumber:
		switch parsingStatement := parsingStatement.(type) {
		case *ast.SelectStatement:
			parsingStatement.Limit = int(thisExpression.Value)
		default:
			logDebugGrammar("This statement does not support LIMIT")
		}
	default:
		logDebugGrammar("limit must be literal integer")
	}
};

select_offset:
OFFSET expression { 
	thisExpression := parsingStack.Pop()
	switch thisExpression := thisExpression.(type) {
	case *ast.LiteralNumber:
		switch parsingStatement := parsingStatement.(type) {
		case *ast.SelectStatement:
			parsingStatement.Offset = int(thisExpression.Value)
		default:
			logDebugGrammar("This statement does not support OFFSET")
		}
	default:
		logDebugGrammar("offset must be literal integer")
	}
};

expression:
expr {
	logDebugGrammar("EXPRESSION")
};

expr:
expr PLUS expr {
	logDebugGrammar("EXPR - PLUS")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewPlusOperator(left.(ast.Expression), right.(ast.Expression)) 
	parsingStack.Push(thisExpression)
}
|
expr MINUS expr {
	logDebugGrammar("EXPR - MINUS")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewSubtractOperator(left.(ast.Expression), right.(ast.Expression)) 
	parsingStack.Push(thisExpression)
}
|
expr MULT expr {
	logDebugGrammar("EXPR - MULT")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewMultiplyOperator(left.(ast.Expression), right.(ast.Expression)) 
	parsingStack.Push(thisExpression)
}
|
expr DIV expr {
	logDebugGrammar("EXPR - DIV")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewDivideOperator(left.(ast.Expression), right.(ast.Expression)) 
	parsingStack.Push(thisExpression)
}
|
expr AND expr {
	logDebugGrammar("EXPR - AND")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewAndOperator([]ast.BooleanExpression{left.(ast.BooleanExpression), right.(ast.BooleanExpression)}) 
	parsingStack.Push(thisExpression)
}
|
expr OR expr {
	logDebugGrammar("EXPR - OR")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewOrOperator([]ast.BooleanExpression{left.(ast.BooleanExpression), right.(ast.BooleanExpression)}) 
	parsingStack.Push(thisExpression)
}
|
expr EQ expr {
	logDebugGrammar("EXPR - EQ")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewEqualToOperator(left.(ast.Expression), right.(ast.Expression)) 
	parsingStack.Push(thisExpression)
}
|
expr LT expr {
	logDebugGrammar("EXPR - LT")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewLessThanOperator(left.(ast.Expression), right.(ast.Expression)) 
	parsingStack.Push(thisExpression)
}
|
expr LTE expr {
	logDebugGrammar("EXPR - LTE")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewLessThanOrEqualOperator(left.(ast.Expression), right.(ast.Expression)) 
	parsingStack.Push(thisExpression)
}
|
expr GT expr {
	logDebugGrammar("EXPR - GT")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewGreaterThanOperator(left.(ast.Expression), right.(ast.Expression)) 
	parsingStack.Push(thisExpression)
}
|
expr GTE expr {
	logDebugGrammar("EXPR - GTE")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewGreaterThanOrEqualOperator(left.(ast.Expression), right.(ast.Expression)) 
	parsingStack.Push(thisExpression)
}
|
expr NE expr {
	logDebugGrammar("EXPR - NE")
	right := parsingStack.Pop()
	left := parsingStack.Pop()
	thisExpression := ast.NewNotEqualToOperator(left.(ast.Expression), right.(ast.Expression)) 
	parsingStack.Push(thisExpression)
}
|
prefix_expr {
	
}
;

prefix_expr: 
NOT prefix_expr {
	logDebugGrammar("EXPR - NOT")
}
|
suffix_expr {
	
};

suffix_expr: 
atom {
	logDebugGrammar("SUFFIX_EXPR")
};

atom: 
NULL {
	logDebugGrammar("NULL")
	thisExpression := ast.NewLiteralNull()
	parsingStack.Push(thisExpression)
}
|
property {

}
// FIXME enable element match and bracket member
//|
//property LBRACKET expression RBRACKET {
//	logDebugGrammar("ATOM - prop[]")
//	rightExpr := parsingStack.Pop().(ast.Expression)
//	leftProp := parsingStack.Pop().(*ast.Property)
//	thisExpression := parser.NewBracketMemberExpression(leftProp, rightExpr)
//	parsingStack.Push(thisExpression)
//}
|
INT { 
	thisExpression := ast.NewLiteralNumber(float64($1.n))
	parsingStack.Push(thisExpression)
}
|
MINUS INT {
	thisExpression := ast.NewLiteralNumber(float64(-$1.n))
	parsingStack.Push(thisExpression)
}
|
REAL {
	thisExpression := ast.NewLiteralNumber($1.f)
	parsingStack.Push(thisExpression)
}
|
MINUS REAL {
	thisExpression := ast.NewLiteralNumber(-$1.f)
	parsingStack.Push(thisExpression)
}
|
STRING {
	thisExpression := ast.NewLiteralString($1.s) 
	parsingStack.Push(thisExpression)
}
|
TRUE {
	thisExpression := ast.NewLiteralBool(true) 
	parsingStack.Push(thisExpression)
}
|
FALSE {
	thisExpression := ast.NewLiteralBool(false) 
	parsingStack.Push(thisExpression)
}
|
LBRACE named_expression_list RBRACE {
	logDebugGrammar("ATOM - {}")
}
|
LBRACKET expression_list RBRACKET {
    logDebugGrammar("ATOM - []")
	exp_list := parsingStack.Pop().([]ast.Expression)
	thisExpression := ast.NewLiteralArray(exp_list)
	parsingStack.Push(thisExpression)
}
|
LPAREN expression RPAREN {
	
}
|
LPAREN select_stmt RPAREN {
	
};

expression_list:
expression {
	logDebugGrammar("EXPRESSION_LIST - EXPRESSION")
	exp_list := make([]ast.Expression, 0)
	exp_list = append(exp_list, parsingStack.Pop().(ast.Expression))
	parsingStack.Push(exp_list)
}
|
expression COMMA expression_list { 
	logDebugGrammar("EXPRESSION_LIST - EXPRESSION COMMA EXPRESSION_LIST")
	rest := parsingStack.Pop().([]ast.Expression)
	last := parsingStack.Pop()
	new_list := make([]ast.Expression, 0, len(rest) + 1)
	new_list = append(new_list, last.(ast.Expression))
	for _, v := range rest {
		new_list = append(new_list, v)
	}
	parsingStack.Push(new_list)
};

named_expression_list:
named_expression_single {
	
}
|
named_expression_single COMMA named_expression_list {
	last := parsingStack.Pop().(*ast.LiteralObject)
	rest := parsingStack.Pop().(*ast.LiteralObject)
	for k,v := range last.Value {
		rest.Value[k] = v
	}
	parsingStack.Push(rest)
};

named_expression_single:   
STRING COLON expression {  
	thisKey := $1.s
	thisValue := parsingStack.Pop().(ast.Expression)
	thisExpression := ast.NewLiteralObject(map[string]ast.Expression{thisKey: thisValue})
	parsingStack.Push(thisExpression) 
};

property:
IDENTIFIER {
	thisExpression := ast.NewProperty($1.s) 
	parsingStack.Push(thisExpression) 
}
|
IDENTIFIER DOT property {
	thisValue := parsingStack.Pop().(*ast.Property)
	thisExpression := ast.NewProperty($1.s + "." + thisValue.Path)
	parsingStack.Push(thisExpression)
};