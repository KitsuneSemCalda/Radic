package ast

type CaseClause struct {
	Value Expr
	Stmts []Stmt
}
