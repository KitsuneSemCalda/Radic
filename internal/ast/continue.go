package ast

// Continue is a continue statement.
type Continue struct{}

func (*Continue) stmt() {}
