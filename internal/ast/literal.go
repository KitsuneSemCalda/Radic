package ast

type NumberLit struct {
	Value string
}

func (*NumberLit) expr() {}

type StringLit struct {
	Value string
}

func (*StringLit) expr() {}

type FloatLit struct {
	Value string
}

func (*FloatLit) expr() {}

type CharLit struct {
	Value string
}

func (*CharLit) expr() {}

type BoolLit struct {
	Value string
}

func (*BoolLit) expr() {}

type NilLit struct {
}

func (*NilLit) expr() {}
