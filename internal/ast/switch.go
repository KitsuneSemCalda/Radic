package ast

type Switch struct {
	Disc    Expr
	Cases   []CaseClause
	Default []Stmt
}

func (*Switch) stmt() {}
