package main

import (
	"fmt"
	"strings"
	"unicode"
)

/******************************************************************************
 * 3) LEXER
 ******************************************************************************/

var codasylKeywords = map[string]TokenType{
	"RECORD":     RECORD,
	"SET":        SETKW,
	"OWNER":      OWNER,
	"ORDER":      ORDER,
	"SORTED":     SORTED,
	"BY":         BY,
	"KEY":        KEY,
	"DESCENDING": DESCENDING,
	"MEMBER":     MEMBER,
	"INSERTION":  INSERTION,
	"AUTOMATIC":  AUTOMATIC,
	"RETENTION":  RETENTION,
	"MANDATORY":  MANDATORY,
	"LOCATION":   LOCATION,
	"MODE":       MODE,
	"IS":         IS,
	"CALC":       CALC,
	"USING":      USING,
	"DUPLICATES": DUPLICATES,
	"ARE":        ARE,
	"NOT":        NOT,
	"ALLOWED":    ALLOWED,
	"SYSTEM":     SYSTEM,
	"OCCURS":     OCCURS,
	"TIMES":      TIMES,
	"TYPE":       TYPEKW,
	"DECIMAL":    DECIMAL,
	"FIXED":      FIXED,
	"CHARACTER":  CHARACTER,

	// DML
	"FIND":      FIND,
	"DUPLICATE": DUPLICATE,
	"GET":       GET,
	"NEXT":      NEXT,
	"FIRST":     FIRST,
	"PRIOR":     PRIOR,
	"LAST":      LAST,
	"OF":        OWNEROF,
	"STORE":     STORE,
	"INSERT":    INSERT,
	"REMOVE":    REMOVE,
	"MODIFY":    MODIFY,
	"DELETE":    DELETE,
	"ALL":       ALLKW,
	"CALC-KEY":  CALCKEY,
}

type Lexer struct {
	input  string
	pos    int
	line   int
	column int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: EOF, Literal: "", Line: l.line, Column: l.column}
	}

	ch := l.input[l.pos]

	// Check newline
	if ch == '\n' {
		tok := Token{Type: EOL, Literal: "", Line: l.line, Column: l.column}
		l.advance()
		l.line++
		l.column = 0
		return tok
	}

	switch ch {
	case '=':
		tok := l.newToken(EQUALS, string(ch))
		l.advance()
		return tok
	case ',':
		tok := l.newToken(COMMA, string(ch))
		l.advance()
		return tok
	case '(':
		tok := l.newToken(LPAREN, string(ch))
		l.advance()
		return tok
	case ')':
		tok := l.newToken(RPAREN, string(ch))
		l.advance()
		return tok
	}

	// If we see "CALC-KEY":
	if strings.HasPrefix(strings.ToUpper(l.input[l.pos:]), "CALC-KEY") {
		start := l.column
		for i := 0; i < len("CALC-KEY"); i++ {
			l.advance()
		}
		return Token{Type: CALCKEY, Literal: "CALC-KEY", Line: l.line, Column: start}
	}

	// Ident/keyword
	if isLetter(ch) || isDigit(ch) {
		startCol := l.column
		ident := l.readIdent()
		upperIdent := strings.ToUpper(ident)
		if tt, ok := codasylKeywords[upperIdent]; ok {
			return Token{Type: tt, Literal: ident, Line: l.line, Column: startCol}
		}
		return Token{Type: IDENT, Literal: ident, Line: l.line, Column: startCol}
	}

	// Otherwise, unrecognized
	tok := l.newToken(EOF, fmt.Sprintf("UNKNOWN(%c)", ch))
	l.advance()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
		} else {
			break
		}
	}
}

func (l *Lexer) readIdent() string {
	start := l.pos
	for l.pos < len(l.input) &&
		(isLetter(l.input[l.pos]) || isDigit(l.input[l.pos]) ||
			l.input[l.pos] == '-' || l.input[l.pos] == '_') {
		l.advance()
	}
	return l.input[start:l.pos]
}

func (l *Lexer) advance() {
	l.pos++
	l.column++
}

func (l *Lexer) newToken(t TokenType, lit string) Token {
	return Token{Type: t, Literal: lit, Line: l.line, Column: l.column}
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}
func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}
