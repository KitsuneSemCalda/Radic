package ast

type Return struct {
	Value Expr
}

func (*Return) stmt() {}
