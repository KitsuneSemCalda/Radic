package ast

// NumberLit is an integer literal expression, holding the literal text
// as written in the source (decimal, hex, binary, or octal).
type NumberLit struct {
	Value string
}

func (*NumberLit) expr() {}

// StringLit is a double-quoted string literal expression, holding the
// literal text including its surrounding quotes.
type StringLit struct {
	Value string
}

func (*StringLit) expr() {}

// FloatLit is a floating-point literal expression, holding the literal
// text as written in the source.
type FloatLit struct {
	Value string
}

func (*FloatLit) expr() {}

// CharLit is a single-quoted character literal expression, holding the
// literal text including its surrounding quotes.
type CharLit struct {
	Value string
}

func (*CharLit) expr() {}

// BoolLit is a boolean literal expression, holding the literal text
// ("true" or "false") as written in the source.
type BoolLit struct {
	Value string
}

func (*BoolLit) expr() {}

// NilLit is the nil literal expression.
type NilLit struct {
}

func (*NilLit) expr() {}
