package main

import (
	_ "embed"
	"fmt"
)

//go:embed CODASYL.txt
var dbdInput string

func main() {
	lexer := NewLexer(dbdInput)
	var tokens []Token
	for {
		tok := lexer.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}

	parser := NewParser(tokens)
	records, sets, finds, dml, err := parser.ParseAll()
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}

	fmt.Println("=== RECORDS ===")
	for i, r := range records {
		fmt.Printf("Record #%d: Name=%q\n", i+1, r.Name)
		for _, ln := range r.Lines {
			fmt.Printf("   -> %s\n", ln)
		}
	}

	fmt.Println("\n=== SETS ===")
	for i, s := range sets {
		fmt.Printf("Set #%d: Name=%q\n", i+1, s.Name)
		for _, ln := range s.Description {
			fmt.Printf("   -> %s\n", ln)
		}
	}

	fmt.Println("\n=== FIND Statements (2-variable only) ===")
	for i, f := range finds {
		fmt.Printf("Find #%d:\n", i+1)
		fmt.Printf("   FullText: %q\n", f.FullText)
		fmt.Printf("   Alias:    %q\n", f.Alias)
		fmt.Printf("   SetName:  %q\n", f.SetName)
	}

	fmt.Println("\n=== DML Statements ===")
	for i, dd := range dml {
		fmt.Printf("DML #%d: %q\n", i+1, dd.Text)
	}
}
