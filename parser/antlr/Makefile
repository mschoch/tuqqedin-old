ANTLR_PATH=~/antlr-3.4/lib/antlr-3.4-complete.jar

all: UNQL2013Lexer.c

UNQL2013Lexer.c: UNQL2013.g
	java -cp ${ANTLR_PATH} org.antlr.Tool -o . UNQL2013.g
