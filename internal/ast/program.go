package ast

// Program is the root of the AST for a single Radic source file, holding
// its top-level declarations in source order.
type Program struct {
	Decls []Decl
}
