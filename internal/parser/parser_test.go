package parser

import (
	"strings"
	"testing"

	"radic/internal/ast"
	"radic/internal/lexer"
	"radic/internal/token"
)

// parseSrc lexes and parses src, failing the test if either step reports
// an error.
func parseSrc(t *testing.T, src string) *ast.Program {
	t.Helper()
	prog, err := Parse(lexer.New(src).Tokenize())
	if err != nil {
		t.Fatalf("Parse(%q) returned error: %v", src, err)
	}
	return prog
}

// parseSrcErr lexes and parses src, failing the test if it does NOT
// report an error, and returns that error.
func parseSrcErr(t *testing.T, src string) error {
	t.Helper()
	_, err := Parse(lexer.New(src).Tokenize())
	if err == nil {
		t.Fatalf("Parse(%q) = nil error, want an error", src)
	}
	return err
}

func TestParseEmptyProgram(t *testing.T) {
	prog := parseSrc(t, "")
	if len(prog.Decls) != 0 {
		t.Errorf("Decls = %v, want empty", prog.Decls)
	}
}

func TestParseVarDeclBuiltinType(t *testing.T) {
	prog := parseSrc(t, "number x = 1;")
	if len(prog.Decls) != 1 {
		t.Fatalf("len(Decls) = %d, want 1", len(prog.Decls))
	}
	d, ok := prog.Decls[0].(*ast.VarDecl)
	if !ok {
		t.Fatalf("Decls[0] = %T, want *ast.VarDecl", prog.Decls[0])
	}
	if d.Type != (ast.Name{Lexeme: "number", Builtin: true}) {
		t.Errorf("Type = %+v, want {number true}", d.Type)
	}
	if d.Name != "x" {
		t.Errorf("Name = %q, want %q", d.Name, "x")
	}
	lit, ok := d.Init.(*ast.NumberLit)
	if !ok || lit.Value != "1" {
		t.Errorf("Init = %+v, want NumberLit{1}", d.Init)
	}
}

func TestParseVarDeclPinnedWidthTypes(t *testing.T) {
	tests := []string{"i8", "i16", "i32", "i64", "u8", "u16", "u32", "u64"}
	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			prog := parseSrc(t, name+" x = 1;")
			d := prog.Decls[0].(*ast.VarDecl)
			if d.Type != (ast.Name{Lexeme: name, Builtin: true}) {
				t.Errorf("Type = %+v, want {%s true}", d.Type, name)
			}
		})
	}
}

func TestParseVarDeclWithoutInit(t *testing.T) {
	prog := parseSrc(t, "number x;")
	d := prog.Decls[0].(*ast.VarDecl)
	if d.Init != nil {
		t.Errorf("Init = %v, want nil", d.Init)
	}
}

func TestParseVarDeclUserDefinedType(t *testing.T) {
	prog := parseSrc(t, "Point p;")
	d := prog.Decls[0].(*ast.VarDecl)
	if d.Type != (ast.Name{Lexeme: "Point", Builtin: false}) {
		t.Errorf("Type = %+v, want {Point false}", d.Type)
	}
	if d.Name != "p" {
		t.Errorf("Name = %q, want %q", d.Name, "p")
	}
}

func TestParseStructDecl(t *testing.T) {
	prog := parseSrc(t, "struct Point { number x; number y; }")
	d, ok := prog.Decls[0].(*ast.StructDecl)
	if !ok {
		t.Fatalf("Decls[0] = %T, want *ast.StructDecl", prog.Decls[0])
	}
	if d.Name != "Point" {
		t.Errorf("Name = %q, want %q", d.Name, "Point")
	}
	want := []ast.Field{
		{Type: ast.Name{Lexeme: "number", Builtin: true}, Name: "x"},
		{Type: ast.Name{Lexeme: "number", Builtin: true}, Name: "y"},
	}
	if len(d.Fields) != len(want) || d.Fields[0] != want[0] || d.Fields[1] != want[1] {
		t.Errorf("Fields = %+v, want %+v", d.Fields, want)
	}
}

func TestParseStructDeclEmptyBody(t *testing.T) {
	prog := parseSrc(t, "struct Empty {}")
	d := prog.Decls[0].(*ast.StructDecl)
	if len(d.Fields) != 0 {
		t.Errorf("Fields = %v, want empty", d.Fields)
	}
}

func TestParseUnionDecl(t *testing.T) {
	prog := parseSrc(t, "union Value { number n; float f; }")
	d, ok := prog.Decls[0].(*ast.UnionDecl)
	if !ok {
		t.Fatalf("Decls[0] = %T, want *ast.UnionDecl", prog.Decls[0])
	}
	if d.Name != "Value" || len(d.Fields) != 2 {
		t.Errorf("UnionDecl = %+v, want Name=Value with 2 fields", d)
	}
}

func TestParseEnumDecl(t *testing.T) {
	prog := parseSrc(t, "enum Color { Red, Green = 5, Blue }")
	d, ok := prog.Decls[0].(*ast.EnumDecl)
	if !ok {
		t.Fatalf("Decls[0] = %T, want *ast.EnumDecl", prog.Decls[0])
	}
	if d.Name != "Color" {
		t.Errorf("Name = %q, want %q", d.Name, "Color")
	}
	if len(d.Variants) != 3 {
		t.Fatalf("len(Variants) = %d, want 3", len(d.Variants))
	}
	if d.Variants[0].Name != "Red" || d.Variants[0].Value != nil {
		t.Errorf("Variants[0] = %+v, want Red with no value", d.Variants[0])
	}
	greenVal, ok := d.Variants[1].Value.(*ast.NumberLit)
	if d.Variants[1].Name != "Green" || !ok || greenVal.Value != "5" {
		t.Errorf("Variants[1] = %+v, want Green = 5", d.Variants[1])
	}
	if d.Variants[2].Name != "Blue" || d.Variants[2].Value != nil {
		t.Errorf("Variants[2] = %+v, want Blue with no value", d.Variants[2])
	}
}

func TestParseFuncDeclFull(t *testing.T) {
	prog := parseSrc(t, "func add(number a, number b) -> number { return a + b; }")
	d, ok := prog.Decls[0].(*ast.FuncDecl)
	if !ok {
		t.Fatalf("Decls[0] = %T, want *ast.FuncDecl", prog.Decls[0])
	}
	if d.Name != "add" {
		t.Errorf("Name = %q, want %q", d.Name, "add")
	}
	want := []ast.Param{
		{Type: ast.Name{Lexeme: "number", Builtin: true}, Name: "a"},
		{Type: ast.Name{Lexeme: "number", Builtin: true}, Name: "b"},
	}
	if len(d.Params) != 2 || d.Params[0] != want[0] || d.Params[1] != want[1] {
		t.Errorf("Params = %+v, want %+v", d.Params, want)
	}
	if d.Return != (ast.Name{Lexeme: "number", Builtin: true}) {
		t.Errorf("Return = %+v, want {number true}", d.Return)
	}
	if len(d.Body.Stmts) != 1 {
		t.Fatalf("len(Body.Stmts) = %d, want 1", len(d.Body.Stmts))
	}
	if _, ok := d.Body.Stmts[0].(*ast.Return); !ok {
		t.Errorf("Body.Stmts[0] = %T, want *ast.Return", d.Body.Stmts[0])
	}
}

func TestParseFuncDeclNoParamsNoReturn(t *testing.T) {
	prog := parseSrc(t, "func run() { }")
	d := prog.Decls[0].(*ast.FuncDecl)
	if len(d.Params) != 0 {
		t.Errorf("Params = %v, want empty", d.Params)
	}
	if d.Return != (ast.Name{}) {
		t.Errorf("Return = %+v, want zero value", d.Return)
	}
}

func TestParseInclude(t *testing.T) {
	err := parseSrcErr(t, `include "foo";`)
	if !strings.Contains(err.Error(), "include directives are not supported yet") {
		t.Errorf("err = %v, want mention of unsupported include", err)
	}
}

func TestParseUnexpectedTopLevelToken(t *testing.T) {
	err := parseSrcErr(t, "42;")
	if !strings.Contains(err.Error(), "expected a declaration") {
		t.Errorf("err = %v, want \"expected a declaration\"", err)
	}
}

func TestParseMultipleDecls(t *testing.T) {
	prog := parseSrc(t, "number a; number b; func f() {}")
	if len(prog.Decls) != 3 {
		t.Fatalf("len(Decls) = %d, want 3", len(prog.Decls))
	}
}

func TestParseIfWithoutElse(t *testing.T) {
	prog := parseSrc(t, "func f() { if (x) { return; } }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	ifStmt, ok := fn.Body.Stmts[0].(*ast.If)
	if !ok {
		t.Fatalf("Stmts[0] = %T, want *ast.If", fn.Body.Stmts[0])
	}
	if _, ok := ifStmt.Cond.(*ast.Ident); !ok {
		t.Errorf("Cond = %T, want *ast.Ident", ifStmt.Cond)
	}
	if ifStmt.Else != nil {
		t.Errorf("Else = %v, want nil", ifStmt.Else)
	}
}

func TestParseIfElse(t *testing.T) {
	prog := parseSrc(t, "func f() { if (x) { return; } else { return; } }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	ifStmt := fn.Body.Stmts[0].(*ast.If)
	if _, ok := ifStmt.Else.(*ast.Block); !ok {
		t.Errorf("Else = %T, want *ast.Block", ifStmt.Else)
	}
}

func TestParseElseIfChain(t *testing.T) {
	prog := parseSrc(t, "func f() { if (a) { return; } else if (b) { return; } else { return; } }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	outer := fn.Body.Stmts[0].(*ast.If)
	inner, ok := outer.Else.(*ast.If)
	if !ok {
		t.Fatalf("outer.Else = %T, want *ast.If", outer.Else)
	}
	if _, ok := inner.Else.(*ast.Block); !ok {
		t.Errorf("inner.Else = %T, want *ast.Block", inner.Else)
	}
}

func TestParseWhile(t *testing.T) {
	prog := parseSrc(t, "func f() { while (x) { break; } }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	w, ok := fn.Body.Stmts[0].(*ast.While)
	if !ok {
		t.Fatalf("Stmts[0] = %T, want *ast.While", fn.Body.Stmts[0])
	}
	if len(w.Body.Stmts) != 1 {
		t.Fatalf("len(Body.Stmts) = %d, want 1", len(w.Body.Stmts))
	}
	if _, ok := w.Body.Stmts[0].(*ast.Break); !ok {
		t.Errorf("Body.Stmts[0] = %T, want *ast.Break", w.Body.Stmts[0])
	}
}

func TestParseForFullClauses(t *testing.T) {
	prog := parseSrc(t, "func f() { for (number i = 0; i < 10; i = i + 1) { continue; } }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	f, ok := fn.Body.Stmts[0].(*ast.For)
	if !ok {
		t.Fatalf("Stmts[0] = %T, want *ast.For", fn.Body.Stmts[0])
	}
	if _, ok := f.Init.(*ast.VarDeclStmt); !ok {
		t.Errorf("Init = %T, want *ast.VarDeclStmt", f.Init)
	}
	if f.Cond == nil {
		t.Error("Cond = nil, want a condition")
	}
	if _, ok := f.Post.(*ast.Assign); !ok {
		t.Errorf("Post = %T, want *ast.Assign", f.Post)
	}
	if len(f.Body.Stmts) != 1 {
		t.Errorf("len(Body.Stmts) = %d, want 1", len(f.Body.Stmts))
	}
}

func TestParseForOmittedClauses(t *testing.T) {
	prog := parseSrc(t, "func f() { for (;;) { break; } }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	f := fn.Body.Stmts[0].(*ast.For)
	if f.Init != nil || f.Cond != nil || f.Post != nil {
		t.Errorf("For = %+v, want Init, Cond, and Post all nil", f)
	}
}

func TestParseForAssignInit(t *testing.T) {
	prog := parseSrc(t, "func f() { for (i = 0; i < 10; i++) {} }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	f := fn.Body.Stmts[0].(*ast.For)
	if _, ok := f.Init.(*ast.Assign); !ok {
		t.Errorf("Init = %T, want *ast.Assign", f.Init)
	}
	post, ok := f.Post.(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Post = %T, want *ast.ExprStmt", f.Post)
	}
	if _, ok := post.X.(*ast.Unary); !ok {
		t.Errorf("Post.X = %T, want *ast.Unary", post.X)
	}
}

func TestParseForExprInit(t *testing.T) {
	prog := parseSrc(t, "func f() { for (log(); i < 10; i++) {} }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	f := fn.Body.Stmts[0].(*ast.For)
	if _, ok := f.Init.(*ast.ExprStmt); !ok {
		t.Errorf("Init = %T, want *ast.ExprStmt", f.Init)
	}
}

func TestParseSwitchCaseAndDefault(t *testing.T) {
	prog := parseSrc(t, `func f() {
		switch (x) {
		case 1 { break; }
		default { continue; }
		}
	}`)
	fn := prog.Decls[0].(*ast.FuncDecl)
	sw, ok := fn.Body.Stmts[0].(*ast.Switch)
	if !ok {
		t.Fatalf("Stmts[0] = %T, want *ast.Switch", fn.Body.Stmts[0])
	}
	if len(sw.Cases) != 1 {
		t.Fatalf("len(Cases) = %d, want 1", len(sw.Cases))
	}
	val, ok := sw.Cases[0].Value.(*ast.NumberLit)
	if !ok || val.Value != "1" {
		t.Errorf("Cases[0].Value = %+v, want NumberLit{1}", sw.Cases[0].Value)
	}
	if len(sw.Default) != 1 {
		t.Errorf("len(Default) = %d, want 1", len(sw.Default))
	}
}

func TestParseSwitchWithoutDefault(t *testing.T) {
	prog := parseSrc(t, "func f() { switch (x) { case 1 {} } }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	sw := fn.Body.Stmts[0].(*ast.Switch)
	if sw.Default != nil {
		t.Errorf("Default = %v, want nil", sw.Default)
	}
}

func TestParseSwitchExpectsCaseOrDefault(t *testing.T) {
	err := parseSrcErr(t, "func f() { switch (x) { 1 {} } }")
	if !strings.Contains(err.Error(), "expected case or default") {
		t.Errorf("err = %v, want \"expected case or default\"", err)
	}
}

func TestParseReturnWithValue(t *testing.T) {
	prog := parseSrc(t, "func f() { return 1; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	r := fn.Body.Stmts[0].(*ast.Return)
	if lit, ok := r.Value.(*ast.NumberLit); !ok || lit.Value != "1" {
		t.Errorf("Value = %+v, want NumberLit{1}", r.Value)
	}
}

func TestParseReturnWithoutValue(t *testing.T) {
	prog := parseSrc(t, "func f() { return; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	r := fn.Body.Stmts[0].(*ast.Return)
	if r.Value != nil {
		t.Errorf("Value = %v, want nil", r.Value)
	}
}

func TestParseNestedBlock(t *testing.T) {
	prog := parseSrc(t, "func f() { { break; } }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	if _, ok := fn.Body.Stmts[0].(*ast.Block); !ok {
		t.Errorf("Stmts[0] = %T, want *ast.Block", fn.Body.Stmts[0])
	}
}

func TestParseVarDeclStmtInsideFunc(t *testing.T) {
	prog := parseSrc(t, "func f() { number x = 1; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	if _, ok := fn.Body.Stmts[0].(*ast.VarDeclStmt); !ok {
		t.Errorf("Stmts[0] = %T, want *ast.VarDeclStmt", fn.Body.Stmts[0])
	}
}

// TestParseBareLiteralStatement is a regression test: the lexer reuses
// TokenNumber, TokenString, TokenChar, and TokenFloat both for the
// matching builtin type keyword and for literals of that kind, so
// isDeclStart must check for a trailing identifier even for builtin type
// tokens or these bare literal expression statements are misread as the
// start of a (malformed) variable declaration.
func TestParseBareLiteralStatement(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{"number", "func f() { 42; }"},
		{"string", `func f() { "hi"; }`},
		{"char", "func f() { 'a'; }"},
		{"float", "func f() { 1.5; }"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog := parseSrc(t, tt.src)
			fn := prog.Decls[0].(*ast.FuncDecl)
			if _, ok := fn.Body.Stmts[0].(*ast.ExprStmt); !ok {
				t.Errorf("Stmts[0] = %T, want *ast.ExprStmt", fn.Body.Stmts[0])
			}
		})
	}
}

func TestParseExprStmtCall(t *testing.T) {
	prog := parseSrc(t, "func f() { foo(); }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	es, ok := fn.Body.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Stmts[0] = %T, want *ast.ExprStmt", fn.Body.Stmts[0])
	}
	if _, ok := es.X.(*ast.Call); !ok {
		t.Errorf("X = %T, want *ast.Call", es.X)
	}
}

func TestParseAssignment(t *testing.T) {
	prog := parseSrc(t, "func f() { x = 1; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	a, ok := fn.Body.Stmts[0].(*ast.Assign)
	if !ok {
		t.Fatalf("Stmts[0] = %T, want *ast.Assign", fn.Body.Stmts[0])
	}
	if _, ok := a.Target.(*ast.Ident); !ok {
		t.Errorf("Target = %T, want *ast.Ident", a.Target)
	}
	if a.Op != token.TokenAssign {
		t.Errorf("Op = %v, want %v", a.Op, token.TokenAssign)
	}
}

func TestParseCompoundAssignment(t *testing.T) {
	prog := parseSrc(t, "func f() { x += 1; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	a := fn.Body.Stmts[0].(*ast.Assign)
	if a.Op != token.TokenPlusEquals {
		t.Errorf("Op = %v, want %v", a.Op, token.TokenPlusEquals)
	}
}

func TestParseAssignmentToMember(t *testing.T) {
	prog := parseSrc(t, "func f() { obj.field = 1; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	a := fn.Body.Stmts[0].(*ast.Assign)
	if _, ok := a.Target.(*ast.Member); !ok {
		t.Errorf("Target = %T, want *ast.Member", a.Target)
	}
}

func TestParseInvalidAssignmentTarget(t *testing.T) {
	err := parseSrcErr(t, "func f() { 1 = 2; }")
	if !strings.Contains(err.Error(), "invalid assignment target") {
		t.Errorf("err = %v, want \"invalid assignment target\"", err)
	}
}

func TestParseBinaryPrecedence(t *testing.T) {
	// "1 + 2 * 3" should bind as "1 + (2 * 3)".
	prog := parseSrc(t, "func f() { return 1 + 2 * 3; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	r := fn.Body.Stmts[0].(*ast.Return)
	top, ok := r.Value.(*ast.Binary)
	if !ok || top.Op != token.TokenPlus {
		t.Fatalf("Value = %+v, want top-level Binary(+)", r.Value)
	}
	if _, ok := top.Left.(*ast.NumberLit); !ok {
		t.Errorf("Left = %T, want *ast.NumberLit", top.Left)
	}
	right, ok := top.Right.(*ast.Binary)
	if !ok || right.Op != token.TokenStar {
		t.Errorf("Right = %+v, want Binary(*)", top.Right)
	}
}

func TestParseLeftAssociativity(t *testing.T) {
	// "1 - 2 - 3" should bind as "(1 - 2) - 3".
	prog := parseSrc(t, "func f() { return 1 - 2 - 3; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	r := fn.Body.Stmts[0].(*ast.Return)
	top := r.Value.(*ast.Binary)
	if top.Op != token.TokenMinus {
		t.Fatalf("top.Op = %v, want %v", top.Op, token.TokenMinus)
	}
	left, ok := top.Left.(*ast.Binary)
	if !ok || left.Op != token.TokenMinus {
		t.Fatalf("Left = %+v, want nested Binary(-)", top.Left)
	}
	if _, ok := top.Right.(*ast.NumberLit); !ok {
		t.Errorf("Right = %T, want *ast.NumberLit", top.Right)
	}
}

func TestParseLogicalPrecedence(t *testing.T) {
	// "a or b and c" should bind as "a or (b and c)" since and > or.
	prog := parseSrc(t, "func f() { return a or b and c; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	r := fn.Body.Stmts[0].(*ast.Return)
	top, ok := r.Value.(*ast.Binary)
	if !ok || top.Op != token.TokenOr {
		t.Fatalf("Value = %+v, want top-level Binary(or)", r.Value)
	}
	right, ok := top.Right.(*ast.Binary)
	if !ok || right.Op != token.TokenAnd {
		t.Errorf("Right = %+v, want Binary(and)", top.Right)
	}
}

func TestParseParenthesizedExpression(t *testing.T) {
	// "(1 + 2) * 3" should override precedence via parentheses.
	prog := parseSrc(t, "func f() { return (1 + 2) * 3; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	r := fn.Body.Stmts[0].(*ast.Return)
	top := r.Value.(*ast.Binary)
	if top.Op != token.TokenStar {
		t.Fatalf("top.Op = %v, want %v", top.Op, token.TokenStar)
	}
	if _, ok := top.Left.(*ast.Binary); !ok {
		t.Errorf("Left = %T, want *ast.Binary", top.Left)
	}
}

func TestParseUnaryOperators(t *testing.T) {
	tests := []struct {
		name string
		src  string
		op   token.TokenKind
	}{
		{"minus", "-x", token.TokenMinus},
		{"plus", "+x", token.TokenPlus},
		{"logical not", "not x", token.TokenNot},
		{"bitwise not", "~x", token.TokenBitwiseNot},
		{"deref", "*x", token.TokenStar},
		{"address of", "&x", token.TokenBitwiseAnd},
		{"pre increment", "++x", token.TokenPlusPlus},
		{"pre decrement", "--x", token.TokenMinusMinus},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog := parseSrc(t, "func f() { return "+tt.src+"; }")
			fn := prog.Decls[0].(*ast.FuncDecl)
			r := fn.Body.Stmts[0].(*ast.Return)
			u, ok := r.Value.(*ast.Unary)
			if !ok {
				t.Fatalf("Value = %T, want *ast.Unary", r.Value)
			}
			if u.Op != tt.op {
				t.Errorf("Op = %v, want %v", u.Op, tt.op)
			}
			if u.Post {
				t.Errorf("Post = true, want false")
			}
		})
	}
}

func TestParsePostfixIncrementDecrement(t *testing.T) {
	tests := []struct {
		src string
		op  token.TokenKind
	}{
		{"x++", token.TokenPlusPlus},
		{"x--", token.TokenMinusMinus},
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			prog := parseSrc(t, "func f() { "+tt.src+"; }")
			fn := prog.Decls[0].(*ast.FuncDecl)
			es := fn.Body.Stmts[0].(*ast.ExprStmt)
			u, ok := es.X.(*ast.Unary)
			if !ok || !u.Post || u.Op != tt.op {
				t.Errorf("X = %+v, want post Unary(%v)", es.X, tt.op)
			}
		})
	}
}

func TestParseMemberAccess(t *testing.T) {
	prog := parseSrc(t, "func f() { return obj.field; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	r := fn.Body.Stmts[0].(*ast.Return)
	m, ok := r.Value.(*ast.Member)
	if !ok {
		t.Fatalf("Value = %T, want *ast.Member", r.Value)
	}
	if m.Name != "field" || m.Deref {
		t.Errorf("Member = %+v, want Name=field, Deref=false", m)
	}
}

func TestParseMemberAccessDeref(t *testing.T) {
	prog := parseSrc(t, "func f() { return ptr->field; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	r := fn.Body.Stmts[0].(*ast.Return)
	m := r.Value.(*ast.Member)
	if !m.Deref {
		t.Errorf("Deref = false, want true")
	}
}

func TestParseCallWithArgs(t *testing.T) {
	prog := parseSrc(t, "func f() { return sum(1, 2, 3); }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	r := fn.Body.Stmts[0].(*ast.Return)
	call, ok := r.Value.(*ast.Call)
	if !ok {
		t.Fatalf("Value = %T, want *ast.Call", r.Value)
	}
	if _, ok := call.Fn.(*ast.Ident); !ok {
		t.Errorf("Fn = %T, want *ast.Ident", call.Fn)
	}
	if len(call.Args) != 3 {
		t.Errorf("len(Args) = %d, want 3", len(call.Args))
	}
}

func TestParseCallNoArgs(t *testing.T) {
	prog := parseSrc(t, "func f() { return sum(); }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	r := fn.Body.Stmts[0].(*ast.Return)
	call := r.Value.(*ast.Call)
	if len(call.Args) != 0 {
		t.Errorf("Args = %v, want empty", call.Args)
	}
}

func TestParseChainedPostfix(t *testing.T) {
	prog := parseSrc(t, "func f() { return obj.get().field; }")
	fn := prog.Decls[0].(*ast.FuncDecl)
	r := fn.Body.Stmts[0].(*ast.Return)
	outer, ok := r.Value.(*ast.Member)
	if !ok || outer.Name != "field" {
		t.Fatalf("Value = %+v, want outer Member(field)", r.Value)
	}
	call, ok := outer.Object.(*ast.Call)
	if !ok {
		t.Fatalf("Object = %T, want *ast.Call", outer.Object)
	}
	inner, ok := call.Fn.(*ast.Member)
	if !ok || inner.Name != "get" {
		t.Errorf("Fn = %+v, want Member(get)", call.Fn)
	}
}

func TestParseLiterals(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want ast.Expr
	}{
		{"number", "42", &ast.NumberLit{Value: "42"}},
		{"float", "1.5", &ast.FloatLit{Value: "1.5"}},
		{"string", `"hi"`, &ast.StringLit{Value: `"hi"`}},
		{"char", "'a'", &ast.CharLit{Value: "'a'"}},
		{"true", "true", &ast.BoolLit{Value: "true"}},
		{"false", "false", &ast.BoolLit{Value: "false"}},
		{"nil", "nil", &ast.NilLit{}},
		{"identifier", "x", &ast.Ident{Name: "x"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog := parseSrc(t, "func f() { return "+tt.src+"; }")
			fn := prog.Decls[0].(*ast.FuncDecl)
			r := fn.Body.Stmts[0].(*ast.Return)
			if r.Value == nil {
				t.Fatal("Value = nil")
			}
			gotT, wantT := typeName(r.Value), typeName(tt.want)
			if gotT != wantT {
				t.Fatalf("Value = %T, want %T", r.Value, tt.want)
			}
		})
	}
}

// typeName is a tiny helper for comparing dynamic types in TestParseLiterals.
func typeName(x ast.Expr) string {
	switch x.(type) {
	case *ast.NumberLit:
		return "NumberLit"
	case *ast.FloatLit:
		return "FloatLit"
	case *ast.StringLit:
		return "StringLit"
	case *ast.CharLit:
		return "CharLit"
	case *ast.BoolLit:
		return "BoolLit"
	case *ast.NilLit:
		return "NilLit"
	case *ast.Ident:
		return "Ident"
	default:
		return "unknown"
	}
}

func TestParseBoolLiteralValue(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		{"true", "true"},
		{"false", "false"},
	}
	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			prog := parseSrc(t, "func f() { return "+tt.src+"; }")
			fn := prog.Decls[0].(*ast.FuncDecl)
			r := fn.Body.Stmts[0].(*ast.Return)
			b, ok := r.Value.(*ast.BoolLit)
			if !ok {
				t.Fatalf("Value = %T, want *ast.BoolLit", r.Value)
			}
			if b.Value != tt.want {
				t.Errorf("Value = %q, want %q", b.Value, tt.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr string
	}{
		{"missing semicolon", "number x = 1", "expected SEMICOLON"},
		{"unclosed struct body", "struct S { number x;", "expected RIGHT_BRACE"},
		{"missing func params paren", "func f) {}", "expected LEFT_PAREN"},
		{"missing identifier after type", "func f(number) {}", "expected IDENTIFIER"},
		{"missing if condition paren", "func f() { if x { } }", "expected LEFT_PAREN"},
		{"unterminated expression", "func f() { return; ", "expected RIGHT_BRACE"},
		{"unexpected token in expression", "func f() { return }; }", "expected expression"},
		{"struct missing name", "struct { number x; }", "expected IDENTIFIER"},
		{"struct missing brace after name", "struct S number x; }", "expected LEFT_BRACE"},
		{"struct field missing semicolon", "struct S { number x }", "expected SEMICOLON"},
		{"struct field bad name", "struct S { number 1; }", "expected IDENTIFIER"},
		{"enum missing name", "enum { A }", "expected IDENTIFIER"},
		{"enum missing brace", "enum E A }", "expected LEFT_BRACE"},
		{"enum missing right brace", "enum E { A", "expected RIGHT_BRACE"},
		{"enum bad variant value", "enum E { A = }", "expected expression"},
		{"vardecl missing init expr", "number x = ;", "expected expression"},
		{"while missing paren", "func f() { while x {} }", "expected LEFT_PAREN"},
		{"while missing close paren", "func f() { while (x {} }", "expected RIGHT_PAREN"},
		{"while bad body", "func f() { while (x) y; }", "expected LEFT_BRACE"},
		{"if missing close paren", "func f() { if (x {} }", "expected RIGHT_PAREN"},
		{"if else missing stmt", "func f() { if (x) {} else } }", "expected expression"},
		{"switch missing paren", "func f() { switch x {} }", "expected LEFT_PAREN"},
		{"switch missing close paren", "func f() { switch (x {} }", "expected RIGHT_PAREN"},
		{"switch missing left brace", "func f() { switch (x) case 1 {} }", "expected LEFT_BRACE"},
		{"switch case missing body", "func f() { switch (x) { case 1 } }", "expected LEFT_BRACE"},
		{"switch missing right brace", "func f() { switch (x) { case 1 {} }", "expected RIGHT_BRACE"},
		{"for missing semicolon after cond", "func f() { for (; true) {} }", "expected SEMICOLON"},
		{"for missing paren after post", "func f() { for (;;i++ {} }", "expected RIGHT_PAREN"},
		{"for bad body", "func f() { for (;;) x; }", "expected LEFT_BRACE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseSrcErr(t, tt.src)
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("err = %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestIsTypeToken(t *testing.T) {
	tests := []struct {
		kind token.TokenKind
		want bool
	}{
		{token.TokenNumber, true},
		{token.TokenString, true},
		{token.TokenChar, true},
		{token.TokenFloat, true},
		{token.TokenBool, true},
		{token.TokenError, true},
		{token.TokenVoid, true},
		{token.TokenI8, true},
		{token.TokenI16, true},
		{token.TokenI32, true},
		{token.TokenI64, true},
		{token.TokenU8, true},
		{token.TokenU16, true},
		{token.TokenU32, true},
		{token.TokenU64, true},
		{token.TokenIdentifier, false},
		{token.TokenPlus, false},
	}
	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := isTypeToken(tt.kind); got != tt.want {
				t.Errorf("isTypeToken(%v) = %v, want %v", tt.kind, got, tt.want)
			}
		})
	}
}

func TestIsAssignOp(t *testing.T) {
	tests := []struct {
		kind token.TokenKind
		want bool
	}{
		{token.TokenAssign, true},
		{token.TokenPlusEquals, true},
		{token.TokenMinusEquals, true},
		{token.TokenStarEquals, true},
		{token.TokenSlashEquals, true},
		{token.TokenPercentEquals, true},
		{token.TokenEquals, false},
		{token.TokenPlus, false},
	}
	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := isAssignOp(tt.kind); got != tt.want {
				t.Errorf("isAssignOp(%v) = %v, want %v", tt.kind, got, tt.want)
			}
		})
	}
}
