package ast

// Call is a function call expression, e.g. "fn(args...)".
type Call struct {
	Fn   Expr
	Args []Expr
}

func (*Call) expr() {}
