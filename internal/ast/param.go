package ast

// Param is a single named parameter in a FuncDecl.
type Param struct {
	Type Name
	Name string
}
