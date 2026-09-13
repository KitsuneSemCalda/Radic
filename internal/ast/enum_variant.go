package ast

// EnumVariant is a single variant of an EnumDecl.
type EnumVariant struct {
	Name  string
	Value Expr // Can be nil if implied
}
