#include "UNQL2013Util.h"

UNQL2013Parser_input_return parse(char * inputString) {
    UNQL2013Parser_input_return ast; 
    pANTLR3_INPUT_STREAM           input;
    pUNQL2013Lexer               lex;
    pANTLR3_COMMON_TOKEN_STREAM    tokens;
    pUNQL2013Parser              parser;

    pANTLR3_UINT8 input_string = (pANTLR3_UINT8)inputString;
    input = antlr3StringStreamNew(input_string, ANTLR3_ENC_8BIT, strlen(input_string),(pANTLR3_UINT8)"ABCD"); 
    lex    = UNQL2013LexerNew(input);
    tokens = antlr3CommonTokenStreamSourceNew(ANTLR3_SIZE_HINT, TOKENSOURCE(lex));
    parser = UNQL2013ParserNew(tokens);

    ast = parser->input(parser);

    if (parser->pParser->rec->state->errorCount > 0) {
        fprintf(stderr, "The parser returned %d errors, tree walking aborted.\n", parser->pParser->rec->state->errorCount);
        exit(5);
    } else {
        printf("tree: \n%s\n\n", ast.tree->toStringTree(ast.tree)->chars);
    }

    GoWalkTree(ast.tree);

    // Must manually clean up
    //parser ->free(parser);
    //tokens ->free(tokens);
    //lex    ->free(lex);
    //input  ->close(input);

    return ast;
}

// void myToString(pANTLR3_BASE_TREE   tree)
// {
//     //return tree->toString(tree);
//     return;
// }

pANTLR3_BASE_TREE myVectorGet(pANTLR3_VECTOR vector, ANTLR3_UINT32 entry) {
    return (pANTLR3_BASE_TREE) vector->get(vector, entry);
}

ANTLR3_UINT32 myVectorSize(pANTLR3_VECTOR vector) {
    return vector->size(vector);
}

char * myPrintTree(pANTLR3_BASE_TREE   tree) {
    // ANTLR3_UINT32   i;
    // ANTLR3_UINT32   n;
    // pANTLR3_BASE_TREE   t;
    pANTLR3_STRING s;


    // s = tree->toString(tree);
    // printf("In node %s\n", (*s).chars);

    // if (tree->children == NULL || tree->children->size(tree->children) == 0) {
    //     return;
    // }

    // if (tree->children != NULL) {
    //     n = tree->children->size(tree->children);

    //     for (i = 0; i < n; i++) {
    //         t = (pANTLR3_BASE_TREE) tree->children->get(tree->children, i);

    //         walkTree(t);
    //     }
    // }
    s = tree->toString(tree);
    return (*s).chars;
}