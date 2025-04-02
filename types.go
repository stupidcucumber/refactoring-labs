package main

/******************************************************************************
 * 4) DATA STRUCTS FOR PARSED OBJECTS
 ******************************************************************************/

// Record: e.g. RECORD Customer ...
type Record struct {
	Name  string
	Lines []string
}

// SetDef: e.g. SET CT ...
type SetDef struct {
	Name        string
	Description []string
}

// We refine FindStmt to hold the parsed form:
type FindStmt struct {
	Alias    string // e.g. "R1"
	SetName  string // e.g. "CN"
	FullText string // raw text (optional)
}

// DMLStmt for any other statements like STORE, INSERT, etc.
type DMLStmt struct {
	Text string
}
