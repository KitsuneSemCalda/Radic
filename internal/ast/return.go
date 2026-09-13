package ast

// Return is a return statement. Value is nil for a bare "return" with
// no value.
type Return struct {
	Value Expr
}

func (*Return) stmt() {}
