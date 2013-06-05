#include "UNQL2013Lexer.h"
#include "UNQL2013Parser.h"
 
UNQL2013Parser_input_return parse(char * inputString);
char * myPrintTree(pANTLR3_BASE_TREE   tree);
//void myToString(pANTLR3_BASE_TREE   tree)
ANTLR3_UINT32 myVectorSize(pANTLR3_VECTOR vector);
pANTLR3_BASE_TREE myVectorGet(pANTLR3_VECTOR vector, ANTLR3_UINT32 entry);