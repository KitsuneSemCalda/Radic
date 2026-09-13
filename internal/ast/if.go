package ast

type If struct {
	Cond Expr
	Then *Block
	Else Stmt
}

func (*If) stmt() {}
