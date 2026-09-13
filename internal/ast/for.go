package ast

// For is a C-style for loop statement: "for (Init; Cond; Post) Body".
// Init and Post may be nil when omitted.
type For struct {
	Init Stmt
	Cond Expr
	Post Stmt
	Body *Block
}

func (*For) stmt() {}
