package ast

type Block struct {
	Stmts []Stmt
}

func (*Block) stmt() {}
