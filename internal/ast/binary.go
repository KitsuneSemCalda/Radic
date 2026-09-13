package ast

import "radic/internal/token"

// Binary is a binary operator expression, e.g. "left op right".
type Binary struct {
	Op    token.TokenKind
	Left  Expr
	Right Expr
}

func (*Binary) expr() {}
