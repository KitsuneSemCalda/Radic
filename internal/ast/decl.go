package ast

type Decl interface {
	decl()
}

type StructDecl struct {
	Name   string
	Fields []Field
}

func (*StructDecl) decl() {}

type UnionDecl struct {
	Name   string
	Fields []Field
}

func (*UnionDecl) decl() {}

type EnumDecl struct {
	Name     string
	Variants []EnumVariant
}

func (*EnumDecl) decl() {}

type VarDecl struct {
	Type Name
	Name string
	Init Expr
}

func (*VarDecl) decl() {}

type FuncDecl struct {
	Name   string
	Params []Param
	Return Name
	Body   *Block
}

func (*FuncDecl) decl() {}
