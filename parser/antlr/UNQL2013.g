grammar UNQL2013;
options
{
	language=C;
	output=AST;
	backtrack=true;
	ASTLabelType=pANTLR3_BASE_TREE;
//	charVocabulary='\u0000'..'\uFFFE';
}
   
tokens
{
	OBJECT;
	ELEMENT;
	ARRAY;
	STRING;
	INTEGER;
	DOUBLE;
	QUERY;
	ELEMENT_MATCH_ANY;
	ELEMENT_MATCH_ALL;
	FUNCTION;
}

@members
{

}

// lexer rules

fragment A_ :	'a' | 'A';
fragment B_ :	'b' | 'B';
fragment C_ :	'c' | 'C';
fragment D_ :	'd' | 'D';
fragment E_ :	'e' | 'E';
fragment F_ :	'f' | 'F';
fragment G_ :	'g' | 'G';
fragment H_ :	'h' | 'H';
fragment I_ :	'i' | 'I';
fragment J_ :	'j' | 'J';
fragment K_ :	'k' | 'K';
fragment L_ :	'l' | 'L';
fragment M_ :	'm' | 'M';
fragment N_ :	'n' | 'N';
fragment O_ :	'o' | 'O';
fragment P_ :	'p' | 'P';
fragment Q_ :	'q' | 'Q';
fragment R_ :	'r' | 'R';
fragment S_ :	's' | 'S';
fragment T_ :	't' | 'T';
fragment U_ :	'u' | 'U';
fragment V_ :	'v' | 'V';
fragment W_ :	'w' | 'W';
fragment X_ :	'x' | 'X';
fragment Y_ :	'y' | 'Y';
fragment Z_ :	'z' | 'Z';

AND	:	A_ N_ D_;
ALL	:	A_ L_ L_;
ANY	:	A_ N_ Y_;
ASC	:	A_ S_ C_;
AVG	:	A_ V_ G_;
BY	:	B_ Y_;
COUNT	:	C_ O_ U_ N_ T_;
DESC	:	D_ E_ S_ C_;
DISTINCT:	D_ I_ S_ T_ I_ N_ C_ T_;
EXPLAIN	:	E_ X_ P_ L_ A_ I_ N_;
FALSE	:	F_ A_ L_ S_ E_;
GROUP	:	G_ R_ O_ U_ P_;
HAVING	:	H_ A_ V_ I_ N_ G_;
IS	: 	I_ S_;
LIKE	:	L_ I_ K_ E_;
LIMIT	:	L_ I_ M_ I_ T_;
MINN	:	M_ I_ N_;
MAXX	:	M_ A_ X_;
MISSING	:	M_ I_ S_ S_ I_ N_ G_;
NULLL	:	N_ U_ L_ L_;
NOT	:	N_ O_ T_;
OFFSET	:	O_ F_ F_ S_ E_ T_;
OR 	:	O_ R_;
ORDER	:	O_ R_ D_ E_ R_;
SELECT	:	S_ E_ L_ E_ C_ T_;
SUM	:	S_ U_ M_;
TRUE	:	T_ R_ U_ E_;
WHERE	:	W_ H_ E_ R_ E_;

COLON	:	':';
COMMA	:	',';
LBRACE	:	'{';
RBRACE	:	'}';
LBRACKET:	'[';
RBRACKET:	']';
LPAREN	:	'(';
RPAREN	:	')';
DOT	:	'.';
LTT	:	'<';
LTE	:	'<=';
GT	:	'>';
GTE	:	'>=';
EQ	:	'=' | '==';
NEQ 	:	'<>' | '!=';
PLUS	:	'+';
MINUS	:	'-';
ASTERISK:	'*';
DIVIDE	:	'/';
MOD	:	'%';


fragment Digit: '0' .. '9';
fragment HexDigit: ('0' .. '9' | 'A' .. 'F' | 'a' .. 'f');
fragment UnicodeChar: ~('"'| '\\');
fragment StringChar :  UnicodeChar | EscapeSequence;

fragment EscapeSequence
	: '\\' ('\"' | '\\' | '/' | 'b' | 'f' | 'n' | 'r' | 't' | 'u' HexDigit HexDigit HexDigit HexDigit)
	;

fragment Int: '-'? ('0' | '1'..'9' Digit*);
fragment Frac: DOT Digit+;
fragment Exp: ('e' | 'E') ('+' | '-')? Digit+;

WhiteSpace: (' ' | '\r' | '\t' | '\u000C' | '\n') { $channel=HIDDEN; };

Integer: Int;
Double:  Int (Frac Exp? | Exp);
String: '"' StringChar* '"';

ID:	
	( 'A'..'Z' | 'a'..'z' | '_' | '$') ( 'A'..'Z' | 'a'..'z' | '_' | '$' | '0'..'9' )*
;


// parser rules

input
	: unql_stmt EOF
	;

jsonObject
	: object
	;
	
jsonArray
	: array
	;	


object
	: LBRACE (objectElement (COMMA objectElement)*)? RBRACE
	  -> ^(OBJECT objectElement*)
	;
	
objectElement
	: String COLON boolExprOrExpr
	  -> ^(ELEMENT String boolExprOrExpr)
	;	
	
array
	: LBRACKET boolExprOrExpr (COMMA boolExprOrExpr)* RBRACKET
	  -> ^(ARRAY boolExprOrExpr+)
	;

	
value
	: String -> ^(STRING String)
	| Integer -> ^(INTEGER Integer)
	| Double -> ^(DOUBLE Double)
	| object  
	| array  
	| TRUE
	| FALSE
	| NULLL
	;
		
expr
	: expr1
	;
	
	
expr1	
	: expr2 ((PLUS | MINUS) expr2)?
	;
	
expr2
	: expr3 ((ASTERISK | DIVIDE | MOD) expr3)?
	;
	
expr3
	: MINUS? expr4
	;

expr4
	: expr5 (DOT property_path | LBRACKET Integer RBRACKET)*
	;
	
expr5
	: value | property_path | functionCall | parenExpr
	;

boolExpr
	: boolExpr1
	;
	
boolExpr1
	: boolExpr2 ( OR^ boolExpr2 )*
	;
	
boolExpr2
	: boolExpr3 ( AND^ boolExpr3 )*
	;
	
boolExpr3
	: (NOT^)? boolExpr4
	;
	
boolExpr4
	: boolExpr5 ( IS (NOT)? NULLL )?
	;
	
boolExpr5
	:	boolExpr6 | parenBoolExpr | boolFunctionCall
	;
	
boolExpr6
@init{
	int comp = 0;
	int all = 0;
}
	: expr1 (comparison_op {comp=1;} expr1 | LBRACKET (ANY | ALL {all=1;})? boolExpr RBRACKET)
	-> {comp == 1}? ^(comparison_op expr1 expr1)
	-> {all == 1}? ^(ELEMENT_MATCH_ALL boolExpr)
	-> ^(ELEMENT_MATCH_ANY boolExpr)
	;
	
boolExprOrExpr
	: boolExprOrExpr1
	;
	
boolExprOrExpr1
	: boolExprOrExpr2 ( OR^ boolExprOrExpr2 )*
	;
	
boolExprOrExpr2
	: boolExprOrExpr3 ( AND^ boolExprOrExpr3 )*
	;
	
boolExprOrExpr3
	: (NOT^)? boolExprOrExpr4
	;
	
boolExprOrExpr4
	: boolExprOrExpr5 ( IS (NOT)? NULLL )?
	;
	
boolExprOrExpr5
	:	boolExprOrExpr6 | parenBoolExpr | boolFunctionCall
	;
	
boolExprOrExpr6
@init{
	int comp = 0;
	int all = 0;
}
	: expr1 (comparison_op {comp=true;} expr1 | LBRACKET (ANY | ALL {all=true;})? boolExpr RBRACKET)?
/*	-> {comp}? ^(comparison_op expr1 expr1)
	-> {all}? ^(ELEMENT_MATCH_ALL boolExpr)
	-> ^(ELEMENT_MATCH_ANY boolExpr)*/
	;
	
boolFunctionCall
	:	boolFunctionName LPAREN (a+=boolExprOrExpr (COMMA a+=boolExprOrExpr)*)? RPAREN
	-> ^(FUNCTION boolFunctionName $a)
	;
	
boolFunctionName
	:	MISSING
	;	
	
parenBoolExpr 
	:	LPAREN! boolExprOrExpr RPAREN!
	;
	
comparison_op
	: EQ | NEQ | GT | GTE | LTT | LTE | (NOT)? LIKE
	;
	
functionCall
	:	functionName LPAREN (a+=boolExprOrExpr (COMMA a+=boolExprOrExpr)*)? RPAREN
	-> ^(FUNCTION functionName $a)
	;
	
parenExpr 
	:	LPAREN! expr RPAREN!
	;
	
functionName
	:	SUM | AVG | MINN | MAXX | COUNT
	;
	
property_path
	:	ID
	;

unql_stmt
	: (EXPLAIN)? select_stmt
	;

select_stmt
	: select_core (ORDER BY ordering_term_list)? ( LIMIT Integer (OFFSET Integer)? )?
	;
	
ordering_term_list
	: ordering_term (COMMA ordering_term)*
	;
	
ordering_term 
	:	boolExprOrExpr (ASC | DESC)?
	;
	
select_core
	:	SELECT (ALL|DISTINCT)? result_expr (WHERE boolExpr)? (GROUP BY expr_list (HAVING boolExpr)?)? 
	->	^(QUERY ^(SELECT result_expr) ^(WHERE boolExpr))
	;
	
result_expr
	:	ASTERISK | boolExprOrExpr
	;
	
expr_list 
	:	boolExprOrExpr (COMMA boolExprOrExpr)*
	;
