package ast

import "radic/internal/token"

type Assign struct {
	Target Expr
	Op     token.TokenKind
	Value  Expr
}

func (*Assign) stmt() {}
