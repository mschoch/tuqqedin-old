
//line unql.y:2
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

//line unql.y:16
type yySymType struct {
	yys int 
s string 
n int
f float64}

const INT = 57346
const REAL = 57347
const STRING = 57348
const TRUE = 57349
const FALSE = 57350
const NULL = 57351
const IDENTIFIER = 57352
const DOT = 57353
const LBRACKET = 57354
const RBRACKET = 57355
const COMMA = 57356
const LBRACE = 57357
const RBRACE = 57358
const COLON = 57359
const PLUS = 57360
const MINUS = 57361
const MULT = 57362
const DIV = 57363
const SELECT = 57364
const WHERE = 57365
const ORDER = 57366
const BY = 57367
const ASC = 57368
const DESC = 57369
const OFFSET = 57370
const LIMIT = 57371
const LPAREN = 57372
const RPAREN = 57373
const AND = 57374
const OR = 57375
const NOT = 57376
const LT = 57377
const LTE = 57378
const GT = 57379
const GTE = 57380
const EQ = 57381
const NE = 57382
const MOD = 57383
const QUESTION = 57384

var yyToknames = []string{
	"INT",
	"REAL",
	"STRING",
	"TRUE",
	"FALSE",
	"NULL",
	"IDENTIFIER",
	"DOT",
	"LBRACKET",
	"RBRACKET",
	"COMMA",
	"LBRACE",
	"RBRACE",
	"COLON",
	"PLUS",
	"MINUS",
	"MULT",
	"DIV",
	"SELECT",
	"WHERE",
	"ORDER",
	"BY",
	"ASC",
	"DESC",
	"OFFSET",
	"LIMIT",
	"LPAREN",
	"RPAREN",
	"AND",
	"OR",
	"NOT",
	"LT",
	"LTE",
	"GT",
	"GTE",
	"EQ",
	"NE",
	"MOD",
	"QUESTION",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 60
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 192

var yyAct = []int{

	56, 63, 55, 52, 21, 15, 84, 14, 37, 38,
	39, 40, 36, 83, 34, 88, 89, 35, 9, 61,
	11, 7, 41, 42, 16, 44, 45, 46, 47, 43,
	48, 57, 2, 80, 78, 62, 65, 37, 38, 39,
	40, 87, 49, 66, 67, 68, 69, 70, 71, 72,
	73, 74, 75, 76, 77, 37, 38, 39, 40, 82,
	79, 81, 86, 58, 85, 59, 31, 54, 53, 41,
	50, 51, 44, 45, 46, 47, 43, 48, 19, 18,
	60, 91, 33, 90, 64, 92, 12, 6, 65, 93,
	37, 38, 39, 40, 10, 5, 4, 32, 8, 3,
	1, 0, 0, 0, 0, 0, 0, 44, 45, 46,
	47, 43, 48, 22, 24, 25, 26, 27, 20, 31,
	0, 29, 0, 0, 28, 0, 0, 0, 23, 0,
	0, 7, 0, 0, 0, 0, 0, 0, 0, 30,
	0, 0, 0, 17, 22, 24, 25, 26, 27, 20,
	31, 0, 29, 0, 0, 28, 0, 0, 0, 23,
	13, 22, 24, 25, 26, 27, 20, 31, 0, 29,
	30, 0, 28, 0, 17, 0, 23, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 30, 0, 0,
	0, 17,
}
var yyPact = []int{

	-1, -1000, -1000, -6, -1000, -3, 140, -1000, -15, -8,
	-1000, 157, -1000, -1000, -1000, -10, -1000, 157, -1000, -1000,
	-1000, -1000, -1000, 66, -1000, -1000, -1000, -1000, 61, 157,
	109, 54, -1000, -9, 157, 157, -1000, 157, 157, 157,
	157, 157, 157, 157, 157, 157, 157, 157, 157, -1000,
	-1000, -1000, 18, 46, 16, 48, 45, -18, -25, 56,
	-1000, 157, -1000, -1000, 27, -11, -1000, -1000, -1000, -1000,
	72, 37, 19, 19, 19, 19, 19, 19, -1000, 61,
	157, -1000, 157, -1000, -1000, -1000, -1000, 157, -1000, -1000,
	-1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 100, 32, 99, 98, 97, 96, 95, 94, 87,
	86, 0, 1, 84, 82, 80, 5, 24, 79, 78,
	4, 3, 2, 68,
}
var yyR1 = []int{

	0, 1, 2, 3, 6, 7, 9, 10, 10, 8,
	8, 4, 4, 12, 12, 13, 13, 13, 5, 5,
	5, 14, 15, 11, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 16, 16, 16, 16, 17, 17, 18,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	19, 19, 19, 22, 22, 21, 21, 23, 20, 20,
}
var yyR2 = []int{

	0, 1, 3, 1, 2, 2, 1, 1, 1, 0,
	2, 0, 3, 1, 3, 1, 2, 2, 0, 1,
	2, 2, 2, 1, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 1, 2, 1, 1,
	1, 1, 1, 2, 1, 2, 1, 1, 1, 3,
	3, 3, 3, 1, 3, 1, 3, 3, 1, 3,
}
var yyChk = []int{

	-1000, -1, -2, -3, -6, -7, -9, 22, -4, 24,
	-8, 23, -10, 20, -11, -16, -17, 34, -18, -19,
	9, -20, 4, 19, 5, 6, 7, 8, 15, 12,
	30, 10, -5, -14, 29, 25, -11, 18, 19, 20,
	21, 32, 33, 39, 35, 36, 37, 38, 40, -17,
	4, 5, -21, -23, 6, -22, -11, -11, -2, 11,
	-15, 28, -11, -12, -13, -11, -16, -16, -16, -16,
	-16, -16, -16, -16, -16, -16, -16, -16, 16, 14,
	17, 13, 14, 31, 31, -20, -11, 14, 26, 27,
	-21, -11, -22, -12,
}
var yyDef = []int{

	0, -2, 1, 11, 3, 9, 0, 6, 18, 0,
	4, 0, 5, 7, 8, 23, 36, 0, 38, 39,
	40, 41, 42, 0, 44, 46, 47, 48, 0, 0,
	0, 58, 2, 19, 0, 0, 10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 37,
	43, 45, 0, 55, 0, 0, 53, 0, 0, 0,
	20, 0, 21, 12, 13, 15, 24, 25, 26, 27,
	28, 29, 30, 31, 32, 33, 34, 35, 49, 0,
	0, 50, 0, 51, 52, 59, 22, 0, 16, 17,
	56, 57, 54, 14,
}
var yyTok1 = []int{

	1,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c > 0 && c <= len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return fmt.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return fmt.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		fmt.Printf("lex %U %s\n", uint(char), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		fmt.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				fmt.Printf("%s", yyStatname(yystate))
				fmt.Printf("saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					fmt.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				fmt.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		fmt.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line unql.y:38
		{ 
		logDebugGrammar("INPUT") 
	}
	case 2:
		//line unql.y:43
		{
		logDebugGrammar("SELECT_STMT")
	}
	case 3:
		//line unql.y:48
		{ 
		logDebugGrammar("SELECT_COMPOUND") 
	}
	case 4:
		//line unql.y:53
		{ 
		logDebugGrammar("SELECT_CORE")
	}
	case 5:
		//line unql.y:58
		{
		logDebugGrammar("SELECT_SELECT")
	}
	case 6:
		//line unql.y:63
		{ 
		logDebugGrammar("SELECT_SELECT_HEAD")
		if parsingStatement == nil {
			parsingStatement = ast.NewSelectStatement()
		}
	}
	case 7:
		//line unql.y:71
		{ 
		logDebugGrammar("SELECT SELECT TAIL - STAR")
	}
	case 8:
		//line unql.y:74
		{ 
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
	case 9:
		//line unql.y:88
		{ 
		logDebugGrammar("SELECT WHERE - EMPTY")
	}
	case 10:
		//line unql.y:92
		{
		logDebugGrammar("SELECT WHERE - EXPR")
		where_part := parsingStack.Pop().(ast.BooleanExpression)
		switch parsingStatement := parsingStatement.(type) {
		case *ast.SelectStatement:
			parsingStatement.Where = where_part
		default:
			logDebugGrammar("This statement does not support WHERE")
		}
	}
	case 12:
		//line unql.y:106
		{
		
	}
	case 13:
		//line unql.y:112
		{
		
	}
	case 14:
		//line unql.y:116
		{
		
	}
	case 15:
		//line unql.y:121
		{ 
		thisExpression := ast.NewSortExpression(parsingStack.Pop().(ast.Expression), true)
		switch parsingStatement := parsingStatement.(type) {
		case *ast.SelectStatement:
			parsingStatement.Order = append(parsingStatement.Order, thisExpression)
		default:
			logDebugGrammar("This statement does not support ORDER BY")
		}
	}
	case 16:
		//line unql.y:131
		{ 
		thisExpression := ast.NewSortExpression(parsingStack.Pop().(ast.Expression), true)
		switch parsingStatement := parsingStatement.(type) {
		case *ast.SelectStatement:
			parsingStatement.Order = append(parsingStatement.Order, thisExpression)
		default:
			logDebugGrammar("This statement does not support ORDER BY")
		}
	}
	case 17:
		//line unql.y:141
		{ 
		thisExpression := ast.NewSortExpression(parsingStack.Pop().(ast.Expression), false)
		switch parsingStatement := parsingStatement.(type) {
		case *ast.SelectStatement:
			parsingStatement.Order = append(parsingStatement.Order, thisExpression)
		default:
			logDebugGrammar("This statement does not support ORDER BY")
		}
	}
	case 18:
		//line unql.y:152
		{
		
	}
	case 19:
		//line unql.y:156
		{
		
	}
	case 20:
		//line unql.y:160
		{
		
	}
	case 21:
		//line unql.y:166
		{
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
	}
	case 22:
		//line unql.y:182
		{ 
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
	}
	case 23:
		//line unql.y:198
		{
		logDebugGrammar("EXPRESSION")
	}
	case 24:
		//line unql.y:203
		{
		logDebugGrammar("EXPR - PLUS")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewPlusOperator(left.(ast.Expression), right.(ast.Expression)) 
		parsingStack.Push(thisExpression)
	}
	case 25:
		//line unql.y:211
		{
		logDebugGrammar("EXPR - MINUS")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewSubtractOperator(left.(ast.Expression), right.(ast.Expression)) 
		parsingStack.Push(thisExpression)
	}
	case 26:
		//line unql.y:219
		{
		logDebugGrammar("EXPR - MULT")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewMultiplyOperator(left.(ast.Expression), right.(ast.Expression)) 
		parsingStack.Push(thisExpression)
	}
	case 27:
		//line unql.y:227
		{
		logDebugGrammar("EXPR - DIV")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewDivideOperator(left.(ast.Expression), right.(ast.Expression)) 
		parsingStack.Push(thisExpression)
	}
	case 28:
		//line unql.y:235
		{
		logDebugGrammar("EXPR - AND")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewAndOperator([]ast.BooleanExpression{left.(ast.BooleanExpression), right.(ast.BooleanExpression)}) 
		parsingStack.Push(thisExpression)
	}
	case 29:
		//line unql.y:243
		{
		logDebugGrammar("EXPR - OR")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewOrOperator([]ast.BooleanExpression{left.(ast.BooleanExpression), right.(ast.BooleanExpression)}) 
		parsingStack.Push(thisExpression)
	}
	case 30:
		//line unql.y:251
		{
		logDebugGrammar("EXPR - EQ")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewEqualToOperator(left.(ast.Expression), right.(ast.Expression)) 
		parsingStack.Push(thisExpression)
	}
	case 31:
		//line unql.y:259
		{
		logDebugGrammar("EXPR - LT")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewLessThanOperator(left.(ast.Expression), right.(ast.Expression)) 
		parsingStack.Push(thisExpression)
	}
	case 32:
		//line unql.y:267
		{
		logDebugGrammar("EXPR - LTE")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewLessThanOrEqualOperator(left.(ast.Expression), right.(ast.Expression)) 
		parsingStack.Push(thisExpression)
	}
	case 33:
		//line unql.y:275
		{
		logDebugGrammar("EXPR - GT")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewGreaterThanOperator(left.(ast.Expression), right.(ast.Expression)) 
		parsingStack.Push(thisExpression)
	}
	case 34:
		//line unql.y:283
		{
		logDebugGrammar("EXPR - GTE")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewGreaterThanOrEqualOperator(left.(ast.Expression), right.(ast.Expression)) 
		parsingStack.Push(thisExpression)
	}
	case 35:
		//line unql.y:291
		{
		logDebugGrammar("EXPR - NE")
		right := parsingStack.Pop()
		left := parsingStack.Pop()
		thisExpression := ast.NewNotEqualToOperator(left.(ast.Expression), right.(ast.Expression)) 
		parsingStack.Push(thisExpression)
	}
	case 36:
		//line unql.y:299
		{
		
	}
	case 37:
		//line unql.y:305
		{
		logDebugGrammar("EXPR - NOT")
	}
	case 38:
		//line unql.y:309
		{
		
	}
	case 39:
		//line unql.y:314
		{
		logDebugGrammar("SUFFIX_EXPR")
	}
	case 40:
		//line unql.y:319
		{
		logDebugGrammar("NULL")
		thisExpression := ast.NewLiteralNull()
		parsingStack.Push(thisExpression)
	}
	case 41:
		//line unql.y:325
		{
	
	}
	case 42:
		//line unql.y:338
		{ 
		thisExpression := ast.NewLiteralNumber(float64(yyS[yypt-0].n))
		parsingStack.Push(thisExpression)
	}
	case 43:
		//line unql.y:343
		{
		thisExpression := ast.NewLiteralNumber(float64(-yyS[yypt-1].n))
		parsingStack.Push(thisExpression)
	}
	case 44:
		//line unql.y:348
		{
		thisExpression := ast.NewLiteralNumber(yyS[yypt-0].f)
		parsingStack.Push(thisExpression)
	}
	case 45:
		//line unql.y:353
		{
		thisExpression := ast.NewLiteralNumber(-yyS[yypt-1].f)
		parsingStack.Push(thisExpression)
	}
	case 46:
		//line unql.y:358
		{
		thisExpression := ast.NewLiteralString(yyS[yypt-0].s) 
		parsingStack.Push(thisExpression)
	}
	case 47:
		//line unql.y:363
		{
		thisExpression := ast.NewLiteralBool(true) 
		parsingStack.Push(thisExpression)
	}
	case 48:
		//line unql.y:368
		{
		thisExpression := ast.NewLiteralBool(false) 
		parsingStack.Push(thisExpression)
	}
	case 49:
		//line unql.y:373
		{
		logDebugGrammar("ATOM - {}")
	}
	case 50:
		//line unql.y:377
		{
	    logDebugGrammar("ATOM - []")
		exp_list := parsingStack.Pop().([]ast.Expression)
		thisExpression := ast.NewLiteralArray(exp_list)
		parsingStack.Push(thisExpression)
	}
	case 51:
		//line unql.y:384
		{
		
	}
	case 52:
		//line unql.y:388
		{
		
	}
	case 53:
		//line unql.y:393
		{
		logDebugGrammar("EXPRESSION_LIST - EXPRESSION")
		exp_list := make([]ast.Expression, 0)
		exp_list = append(exp_list, parsingStack.Pop().(ast.Expression))
		parsingStack.Push(exp_list)
	}
	case 54:
		//line unql.y:400
		{ 
		logDebugGrammar("EXPRESSION_LIST - EXPRESSION COMMA EXPRESSION_LIST")
		rest := parsingStack.Pop().([]ast.Expression)
		last := parsingStack.Pop()
		new_list := make([]ast.Expression, 0, len(rest) + 1)
		new_list = append(new_list, last.(ast.Expression))
		for _, v := range rest {
			new_list = append(new_list, v)
		}
		parsingStack.Push(new_list)
	}
	case 55:
		//line unql.y:413
		{
		
	}
	case 56:
		//line unql.y:417
		{
		last := parsingStack.Pop().(*ast.LiteralObject)
		rest := parsingStack.Pop().(*ast.LiteralObject)
		for k,v := range last.Value {
			rest.Value[k] = v
		}
		parsingStack.Push(rest)
	}
	case 57:
		//line unql.y:427
		{  
		thisKey := yyS[yypt-2].s
		thisValue := parsingStack.Pop().(ast.Expression)
		thisExpression := ast.NewLiteralObject(map[string]ast.Expression{thisKey: thisValue})
		parsingStack.Push(thisExpression) 
	}
	case 58:
		//line unql.y:435
		{
		thisExpression := ast.NewProperty(yyS[yypt-0].s) 
		parsingStack.Push(thisExpression) 
	}
	case 59:
		//line unql.y:440
		{
		thisValue := parsingStack.Pop().(*ast.Property)
		thisExpression := ast.NewProperty(yyS[yypt-2].s + "." + thisValue.Path)
		parsingStack.Push(thisExpression)
	}
	}
	goto yystack /* stack new state and value */
}
