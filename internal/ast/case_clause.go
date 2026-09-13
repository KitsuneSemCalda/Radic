package ast

// CaseClause is a single "case Value:" clause within a Switch statement.
type CaseClause struct {
	Value Expr
	Stmts []Stmt
}
