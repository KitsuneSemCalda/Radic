package ast

type Stmt interface {
	stmt()
}

type ExprStmt struct {
	X Expr
}

func (*ExprStmt) stmt() {}

type VarDeclStmt struct {
	VarDecl
}

func (*VarDeclStmt) stmt() {}
