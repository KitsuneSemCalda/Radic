package ast

// Expr is implemented by every expression node in the AST.
type Expr interface {
	// expr is unexported so only types in this package can satisfy Expr.
	expr()
}
