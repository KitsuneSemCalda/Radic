package ast

// Switch is a switch statement. Default is nil when there is no default
// clause.
type Switch struct {
	Disc    Expr
	Cases   []CaseClause
	Default []Stmt
}

func (*Switch) stmt() {}
