package ast

// Field is a single named field within a StructDecl or UnionDecl.
type Field struct {
	Type Name
	Name string
}
