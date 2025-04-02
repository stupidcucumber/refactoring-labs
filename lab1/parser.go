package main

import (
	"fmt"
	"strings"
)

/******************************************************************************
 * 5) PARSER
 ******************************************************************************/

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) curToken() Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return Token{Type: EOF, Literal: ""}
}
func (p *Parser) nextToken() {
	p.pos++
}
func (p *Parser) skipEOL() {
	for p.curToken().Type == EOL {
		p.nextToken()
	}
}

// We will store only the FIND statements that have 3 variables (Alias + 2 fields).
func (p *Parser) ParseAll() ([]*Record, []*SetDef, []*FindStmt, []*DMLStmt, error) {
	var records []*Record
	var sets []*SetDef
	var finds []*FindStmt
	var dmls []*DMLStmt

	for {
		p.skipEOL()
		if p.curToken().Type == EOF {
			break
		}

		switch p.curToken().Type {
		case RECORD:
			rec, err := p.parseRecord()
			if err != nil {
				return nil, nil, nil, nil, err
			}
			records = append(records, rec)

		case SETKW:
			s, err := p.parseSetDef()
			if err != nil {
				return nil, nil, nil, nil, err
			}
			sets = append(sets, s)

		case FIND:
			// We'll parse the FIND statement in detail. If it
			// doesn't match the pattern with 3 variables,
			// we just skip it (don't append).
			fstmt, err := p.parseFindStmt()
			if err != nil {
				// we can skip or treat as error. Let's skip it to avoid halting.
				// return nil, nil, nil, nil, err
				fmt.Printf("Skipping non-2-variable FIND statement: %v\n", err)
			} else {
				finds = append(finds, fstmt)
			}

		case STORE, INSERT, REMOVE, MODIFY, DELETE:
			d, err := p.parseDML()
			if err != nil {
				return nil, nil, nil, nil, err
			}
			dmls = append(dmls, d)

		default:
			tok := p.curToken()
			return nil, nil, nil, nil, fmt.Errorf("unexpected token %q at line %d col %d", tok.Literal, tok.Line, tok.Column)
		}
	}
	return records, sets, finds, dmls, nil
}

/******************************************************************************
 * 5a) parseRecord
 *****************************************************************************/
func (p *Parser) parseRecord() (*Record, error) {
	if p.curToken().Type != RECORD {
		return nil, fmt.Errorf("expected RECORD, got %q", p.curToken().Literal)
	}
	p.nextToken() // consume RECORD

	nameTok := p.curToken()
	if nameTok.Type != IDENT {
		return nil, fmt.Errorf("parseRecord: expected record name, got %q", nameTok.Literal)
	}
	recName := nameTok.Literal
	p.nextToken()

	rec := &Record{Name: recName}

	// read until next top-level statement or EOF
	for {
		p.skipEOL()
		switch p.curToken().Type {
		case EOF, RECORD, SETKW, FIND, STORE, INSERT, REMOVE, MODIFY, DELETE:
			return rec, nil
		default:
			line := p.collectLine()
			rec.Lines = append(rec.Lines, line)
		}
	}
}

// gather tokens until EOL/EOF or next top-level statement.
func (p *Parser) collectLine() string {
	var parts []string
	for p.curToken().Type != EOL && p.curToken().Type != EOF {
		parts = append(parts, p.curToken().Literal)
		p.nextToken()
	}
	if p.curToken().Type == EOL {
		p.nextToken()
	}
	return strings.Join(parts, " ")
}

/******************************************************************************
 * 5b) parseSetDef
 *****************************************************************************/
func (p *Parser) parseSetDef() (*SetDef, error) {
	if p.curToken().Type != SETKW {
		return nil, fmt.Errorf("expected SET, got %q", p.curToken().Literal)
	}
	p.nextToken() // consume SET

	nameTok := p.curToken()
	if nameTok.Type != IDENT {
		return nil, fmt.Errorf("parseSetDef: expected set name, got %q", nameTok.Literal)
	}
	setName := nameTok.Literal
	p.nextToken()

	s := &SetDef{Name: setName}

	for {
		p.skipEOL()
		switch p.curToken().Type {
		case EOF, RECORD, SETKW, FIND, STORE, INSERT, REMOVE, MODIFY, DELETE:
			return s, nil
		default:
			line := p.collectLine()
			s.Description = append(s.Description, line)
		}
	}
}

/******************************************************************************
 * 5c) parseFindStmt: We want only if the pattern is:
 *    FIND <Alias> RECORD IN <SetName> (SET)?
 * Example:
 *    FIND R1 RECORD IN CN SET
 *
 * If it doesn't match that, we skip.
 *****************************************************************************/
func (p *Parser) parseFindStmt() (*FindStmt, error) {
	// We know curToken is FIND
	fullText := p.collectLine() // the entire line
	// Re-lex or re-parse that line to get details, or parse in place below:
	// For simplicity, let's parse "in place" using a mini grammar.

	// We'll do a *backup* of the current parser state, parse carefully, and if it fails, restore.
	// But simpler: let's do a short sub-parser with a fresh token slice from that line.

	// 1) Re-inject the line into a separate lexer, parse step by step:
	miniLexer := NewLexer(fullText)
	var miniTokens []Token
	for {
		tk := miniLexer.NextToken()
		miniTokens = append(miniTokens, tk)
		if tk.Type == EOF {
			break
		}
	}
	mParser := &Parser{tokens: miniTokens, pos: 0}

	// 2) Expect: FIND
	if mParser.curToken().Type != FIND {
		return nil, fmt.Errorf("not a FIND statement at all")
	}
	mParser.nextToken() // consume FIND

	// 3) Expect alias as IDENT, e.g. R1
	aliasTok := mParser.curToken()
	if aliasTok.Type != IDENT {
		return nil, fmt.Errorf("expected alias after FIND, got %q", aliasTok.Literal)
	}
	alias := aliasTok.Literal
	mParser.nextToken() // consume alias

	// 4) Expect RECORD
	if mParser.curToken().Type != RECORD {
		return nil, fmt.Errorf("expected RECORD after alias, got %q", mParser.curToken().Literal)
	}
	mParser.nextToken() // consume RECORD

	// 5) Expect IN
	inTok := mParser.curToken()
	if inTok.Type != IDENT || strings.ToUpper(inTok.Literal) != "IN" {
		return nil, fmt.Errorf("expected 'IN' after RECORD, got %q", inTok.Literal)
	}
	mParser.nextToken() // consume 'IN'

	// 6) Expect set name (IDENT), e.g. CN
	setTok := mParser.curToken()
	if setTok.Type != IDENT {
		return nil, fmt.Errorf("expected set name, got %q", setTok.Literal)
	}
	setName := setTok.Literal
	mParser.nextToken()

	// 7) Optionally "SET"
	if mParser.curToken().Type == SETKW || (mParser.curToken().Type == IDENT && strings.ToUpper(mParser.curToken().Literal) == "SET") {
		mParser.nextToken() // consume SET
	}

	// if there's more stuff after the second field (besides EOL/EOF), it fails our 2-variable pattern
	// let's see if we reach EOF or EOL
	for {
		if mParser.curToken().Type == EOF || mParser.curToken().Type == EOL {
			break
		}
		// If there's an unexpected token, we bail
		return nil, fmt.Errorf("found extra tokens after second field: %q", mParser.curToken().Literal)
	}

	// We matched the pattern => return a FindStmt
	return &FindStmt{
		Alias:    alias,
		SetName:  setName,
		FullText: fullText,
	}, nil
}

/******************************************************************************
 * 5d) parseDML
 *****************************************************************************/
func (p *Parser) parseDML() (*DMLStmt, error) {
	line := p.collectLine()
	return &DMLStmt{Text: line}, nil
}
