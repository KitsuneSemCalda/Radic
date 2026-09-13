package ast

import "radic/internal/token"

// Unary is a unary operator expression, such as "-x" (Post false) or
// "x++" (Post true).
type Unary struct {
	Op   token.TokenKind
	X    Expr
	Post bool
}

func (*Unary) expr() {}
