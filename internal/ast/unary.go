package ast

import "radic/internal/token"

type Unary struct {
	Op   token.TokenKind
	X    Expr
	Post bool
}

func (*Unary) expr() {}
