package ast

type EnumVariant struct {
	Name  string
	Value Expr // Can be nil if implied
}
