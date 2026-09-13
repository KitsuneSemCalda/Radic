package ast

// Decl is implemented by every top-level declaration node in the AST.
type Decl interface {
	// decl is unexported so only types in this package can satisfy Decl.
	decl()
}

// StructDecl declares a struct type with a fixed set of named fields.
type StructDecl struct {
	Name   string
	Fields []Field
}

func (*StructDecl) decl() {}

// UnionDecl declares a union type whose fields share the same storage.
type UnionDecl struct {
	Name   string
	Fields []Field
}

func (*UnionDecl) decl() {}

// EnumDecl declares an enum type with a fixed set of named variants.
type EnumDecl struct {
	Name     string
	Variants []EnumVariant
}

func (*EnumDecl) decl() {}

// VarDecl declares a variable with an explicit type and an initializer
// expression.
type VarDecl struct {
	Type Name
	Name string
	Init Expr
}

func (*VarDecl) decl() {}

// FuncDecl declares a named function with its parameters, return type,
// and body.
type FuncDecl struct {
	Name   string
	Params []Param
	Return Name
	Body   *Block
}

func (*FuncDecl) decl() {}
