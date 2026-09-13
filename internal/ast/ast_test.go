package ast

import (
	"testing"

	"radic/internal/token"
)

// TestExprImplementations lists every expression node type. Removing a
// type's unexported expr() method, or misspelling it, fails this test to
// compile instead of silently leaving the type unusable through Expr.
func TestExprImplementations(t *testing.T) {
	nodes := []Expr{
		&Ident{},
		&Unary{},
		&Binary{},
		&Call{},
		&Member{},
		&NumberLit{},
		&StringLit{},
		&FloatLit{},
		&CharLit{},
		&BoolLit{},
		&NilLit{},
	}
	if len(nodes) == 0 {
		t.Fatal("no Expr implementations listed")
	}
}

// TestStmtImplementations lists every statement node type. Removing a
// type's unexported stmt() method, or misspelling it, fails this test to
// compile instead of silently leaving the type unusable through Stmt.
func TestStmtImplementations(t *testing.T) {
	nodes := []Stmt{
		&Assign{},
		&Block{},
		&Break{},
		&Continue{},
		&ExprStmt{},
		&For{},
		&If{},
		&Return{},
		&Switch{},
		&VarDeclStmt{},
		&While{},
	}
	if len(nodes) == 0 {
		t.Fatal("no Stmt implementations listed")
	}
}

// TestDeclImplementations lists every top-level declaration node type.
// Removing a type's unexported decl() method, or misspelling it (e.g. an
// exported Decl() instead), fails this test to compile instead of
// silently leaving the type unusable through Decl.
func TestDeclImplementations(t *testing.T) {
	nodes := []Decl{
		&StructDecl{},
		&UnionDecl{},
		&EnumDecl{},
		&VarDecl{},
		&FuncDecl{},
	}
	if len(nodes) == 0 {
		t.Fatal("no Decl implementations listed")
	}
}

func TestProgramDecls(t *testing.T) {
	a := &StructDecl{Name: "A"}
	b := &FuncDecl{Name: "B"}
	p := Program{Decls: []Decl{a, b}}

	if len(p.Decls) != 2 {
		t.Fatalf("len(Decls) = %d, want 2", len(p.Decls))
	}
	if p.Decls[0] != Decl(a) || p.Decls[1] != Decl(b) {
		t.Errorf("Decls = %v, want [%v %v]", p.Decls, a, b)
	}
}

func TestBlockStmts(t *testing.T) {
	s1 := &Break{}
	s2 := &Continue{}
	b := Block{Stmts: []Stmt{s1, s2}}

	if len(b.Stmts) != 2 || b.Stmts[0] != Stmt(s1) || b.Stmts[1] != Stmt(s2) {
		t.Errorf("Stmts = %v, want [%v %v]", b.Stmts, s1, s2)
	}
}

func TestIfElseIfChaining(t *testing.T) {
	// An "else if" is represented as an *If nested in the outer If's
	// Else field, since Else is typed as Stmt rather than *Block.
	inner := &If{Cond: &BoolLit{Value: "false"}, Then: &Block{}}
	outer := &If{Cond: &BoolLit{Value: "true"}, Then: &Block{}, Else: inner}

	nested, ok := outer.Else.(*If)
	if !ok {
		t.Fatalf("outer.Else = %T, want *If", outer.Else)
	}
	if nested != inner {
		t.Errorf("nested If = %v, want %v", nested, inner)
	}
}

func TestIfWithoutElse(t *testing.T) {
	i := &If{Cond: &BoolLit{Value: "true"}, Then: &Block{}}
	if i.Else != nil {
		t.Errorf("Else = %v, want nil", i.Else)
	}
}

func TestForOptionalClauses(t *testing.T) {
	f := &For{Body: &Block{}}
	if f.Init != nil || f.Cond != nil || f.Post != nil {
		t.Errorf("For{} = %+v, want Init, Cond, and Post all nil", f)
	}
}

func TestSwitchOptionalDefault(t *testing.T) {
	s := &Switch{
		Disc:  &Ident{Name: "x"},
		Cases: []CaseClause{{Value: &NumberLit{Value: "1"}}},
	}
	if s.Default != nil {
		t.Errorf("Default = %v, want nil", s.Default)
	}
	if len(s.Cases) != 1 {
		t.Fatalf("len(Cases) = %d, want 1", len(s.Cases))
	}
}

func TestReturnWithoutValue(t *testing.T) {
	r := &Return{}
	if r.Value != nil {
		t.Errorf("Value = %v, want nil", r.Value)
	}
}

func TestEnumVariantImpliedValue(t *testing.T) {
	v := EnumVariant{Name: "A"}
	if v.Value != nil {
		t.Errorf("Value = %v, want nil for an implied variant", v.Value)
	}
}

func TestNameBuiltin(t *testing.T) {
	tests := []struct {
		name    Name
		builtin bool
	}{
		{Name{Lexeme: "number", Builtin: true}, true},
		{Name{Lexeme: "Point", Builtin: false}, false},
	}

	for _, tt := range tests {
		if tt.name.Builtin != tt.builtin {
			t.Errorf("Name{%q}.Builtin = %v, want %v", tt.name.Lexeme, tt.name.Builtin, tt.builtin)
		}
	}
}

func TestVarDeclStmtEmbedsVarDecl(t *testing.T) {
	s := &VarDeclStmt{VarDecl: VarDecl{
		Type: Name{Lexeme: "number", Builtin: true},
		Name: "x",
		Init: &NumberLit{Value: "1"},
	}}

	if s.Name != "x" {
		t.Errorf("Name = %q, want %q", s.Name, "x")
	}

	// VarDecl.decl() has a pointer receiver, so embedding VarDecl by
	// value promotes it to *VarDeclStmt too: a VarDeclStmt unintentionally
	// also satisfies Decl, alongside its intended Stmt. This test pins
	// down that current behavior; update it if that is fixed on purpose.
	if _, ok := interface{}(s).(Decl); !ok {
		t.Error("VarDeclStmt no longer satisfies Decl via its embedded VarDecl")
	}
}

func TestExprNodeFields(t *testing.T) {
	ident := &Ident{Name: "x"}
	bin := &Binary{Op: token.TokenPlus, Left: ident, Right: &NumberLit{Value: "1"}}
	un := &Unary{Op: token.TokenMinus, X: ident}
	call := &Call{Fn: ident, Args: []Expr{&NumberLit{Value: "1"}, &NumberLit{Value: "2"}}}
	mem := &Member{Object: ident, Name: "field", Deref: true}

	if bin.Left != Expr(ident) {
		t.Errorf("Binary.Left = %v, want %v", bin.Left, ident)
	}
	if un.X != Expr(ident) {
		t.Errorf("Unary.X = %v, want %v", un.X, ident)
	}
	if len(call.Args) != 2 {
		t.Errorf("len(Call.Args) = %d, want 2", len(call.Args))
	}
	if !mem.Deref {
		t.Error("Member.Deref = false, want true")
	}
}
