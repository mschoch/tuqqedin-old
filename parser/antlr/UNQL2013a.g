grammar UNQL2013a;
options
{
	language=Java;
	output=AST;
	backtrack=true;
//	ASTLabelType=pANTLR3_BASE_TREE;
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
AS	:	A_ S_;
ASC	:	A_ S_ C_;
AVG	:	A_ V_ G_;
BY	:	B_ Y_;
CASE	:	C_ A_ S_ E_;
COUNT	:	C_ O_ U_ N_ T_;
DESC	:	D_ E_ S_ C_;
DISTINCT:	D_ I_ S_ T_ I_ N_ C_ T_;
ELSE	:	E_ L_ S_ E_;
END	:	E_ N_ D_;
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
THEN	:	T_ H_ E_ N_;
TRUE	:	T_ R_ U_ E_;
UNIQUE	:	U_ N_ I_ Q_ U_ E_;
VALUED	:	V_ A_ L_ U_ E_ D_;
WHEN	:	W_ H_ E_ N_;
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
	: String COLON expr
	  -> ^(ELEMENT String expr)
	;	
	
array
	: LBRACKET expr (COMMA expr)* RBRACKET
	  -> ^(ARRAY expr+)
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
	: orExpr
	;

orExpr
	: andExpr ( OR^ andExpr )*
	;

andExpr
	: notExpr ( AND^ notExpr )*
	;

notExpr
	: (NOT^)? nullMissingValuedExpr
	;

nullMissingValuedExpr
	: compareExpr ( IS (NOT)? (NULLL | MISSING | VALUED) )?
	;
	
compareExpr
	:	expr1 (comparison_op expr1)?
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
	: value | property_path | functionCall | parenExpr | case_when
	;
	
case_when
	:	CASE (WHEN expr THEN expr)+ (ELSE expr)? END
	;
	
comparison_op
	: EQ | NEQ | GT | GTE | LTT | LTE | (NOT)? LIKE
	;
	
functionCall
	:	functionName LPAREN (a+=expr (COMMA a+=expr)*)? RPAREN
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
	:	expr (ASC | DESC)?
	;
	
select_core
	:	SELECT (DISTINCT|UNIQUE)? result_expr_list (WHERE expr)? (GROUP BY expr_list (HAVING expr)?)? 
	->	^(QUERY ^(SELECT result_expr_list) ^(WHERE expr))
	;
	
result_expr
	:	ASTERISK | expr AS ID
	;

result_expr_list
	:	result_expr (COMMA result_expr)*
	;

expr_list
	:	expr (COMMA expr)*
	;
