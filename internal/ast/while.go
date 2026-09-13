package ast

// While is a while loop statement.
type While struct {
	Cond Expr
	Body *Block
}

func (*While) stmt() {}
