package ast

// Block is a brace-delimited sequence of statements.
type Block struct {
	Stmts []Stmt
}

func (*Block) stmt() {}
