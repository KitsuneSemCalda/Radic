package ast

// If is an if statement. Else is nil when there is no else branch; it
// holds a *Block for a plain else branch or another *If for an
// "else if" chain.
type If struct {
	Cond Expr
	Then *Block
	Else Stmt
}

func (*If) stmt() {}
