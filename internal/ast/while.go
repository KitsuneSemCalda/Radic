package ast

type While struct {
	Cond Expr
	Body *Block
}

func (*While) stmt() {}
