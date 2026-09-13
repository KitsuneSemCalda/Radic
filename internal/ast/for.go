package ast

type For struct {
	Init Stmt
	Cond Expr
	Post Stmt
	Body *Block
}

func (*For) stmt() {}
