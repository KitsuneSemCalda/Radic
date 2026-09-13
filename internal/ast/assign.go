package ast

import "radic/internal/token"

// Assign is an assignment statement, e.g. "target op= value".
type Assign struct {
	Target Expr
	Op     token.TokenKind
	Value  Expr
}

func (*Assign) stmt() {}
