package ast

// Stmt is implemented by every statement node in the AST.
type Stmt interface {
	// stmt is unexported so only types in this package can satisfy Stmt.
	stmt()
}

// ExprStmt is an expression used in statement position, such as a bare
// call expression followed by a semicolon.
type ExprStmt struct {
	X Expr
}

func (*ExprStmt) stmt() {}

// VarDeclStmt wraps a VarDecl so it can appear as a statement inside a
// function body, e.g. inside a Block.
type VarDeclStmt struct {
	VarDecl
}

func (*VarDeclStmt) stmt() {}
