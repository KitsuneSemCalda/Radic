package ast

// Member is a member access expression, e.g. "object.Name" or, when
// Deref is true, "object->Name".
type Member struct {
	Object Expr
	Name   string
	Deref  bool
}

func (*Member) expr() {}
