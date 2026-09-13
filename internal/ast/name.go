package ast

// Name is a reference to a type, such as the declared type of a Field,
// Param, or VarDecl. Builtin reports whether Lexeme names a primitive
// builtin type (e.g. "number") rather than a user-defined one.
type Name struct {
	Lexeme  string
	Builtin bool
}
