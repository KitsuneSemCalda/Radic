package ast

import "radic/internal/token"

type Binary struct {
	Op    token.TokenKind
	Left  Expr
	Right Expr
}

func (*Binary) expr() {}
