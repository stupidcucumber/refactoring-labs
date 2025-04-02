package main

/******************************************************************************
 * 1) TOKEN DEFINITIONS
 ******************************************************************************/

type TokenType string

const (
	// Special
	EOF      TokenType = "EOF"
	EOL      TokenType = "EOL"
	IDENT    TokenType = "IDENT"
	INTEGER  TokenType = "INTEGER"
	LABEL    TokenType = "LABEL"
	SKIPLINE TokenType = "SKIPLINE"

	// Operators/punctuation
	EQUALS TokenType = "="
	COMMA  TokenType = ","
	LPAREN TokenType = "("
	RPAREN TokenType = ")"

	// CODASYL keywords
	RECORD     TokenType = "RECORD"
	SETKW      TokenType = "SET"
	OWNER      TokenType = "OWNER"
	ORDER      TokenType = "ORDER"
	SORTED     TokenType = "SORTED"
	BY         TokenType = "BY"
	KEY        TokenType = "KEY"
	DESCENDING TokenType = "DESCENDING"
	MEMBER     TokenType = "MEMBER"
	INSERTION  TokenType = "INSERTION"
	AUTOMATIC  TokenType = "AUTOMATIC"
	RETENTION  TokenType = "RETENTION"
	MANDATORY  TokenType = "MANDATORY"
	LOCATION   TokenType = "LOCATION"
	MODE       TokenType = "MODE"
	IS         TokenType = "IS"
	CALC       TokenType = "CALC"
	USING      TokenType = "USING"
	DUPLICATES TokenType = "DUPLICATES"
	ARE        TokenType = "ARE"
	NOT        TokenType = "NOT"
	ALLOWED    TokenType = "ALLOWED"
	SYSTEM     TokenType = "SYSTEM"
	OCCURS     TokenType = "OCCURS"
	TIMES      TokenType = "TIMES"
	TYPEKW     TokenType = "TYPE"
	DECIMAL    TokenType = "DECIMAL"
	FIXED      TokenType = "FIXED"
	CHARACTER  TokenType = "CHARACTER"

	// DML / statements
	FIND      TokenType = "FIND"
	DUPLICATE TokenType = "DUPLICATE"
	GET       TokenType = "GET"
	NEXT      TokenType = "NEXT"
	FIRST     TokenType = "FIRST"
	PRIOR     TokenType = "PRIOR"
	LAST      TokenType = "LAST"
	OWNEROF   TokenType = "OWNEROF" // e.g. "FIND OWNER OF ..."
	STORE     TokenType = "STORE"
	INSERT    TokenType = "INSERT"
	REMOVE    TokenType = "REMOVE"
	MODIFY    TokenType = "MODIFY"
	DELETE    TokenType = "DELETE"
	ALLKW     TokenType = "ALL"

	CALCKEY TokenType = "CALC-KEY" // you can parse as one token or handle differently
)

/******************************************************************************
 * 2) TOKEN STRUCT
 ******************************************************************************/

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}
