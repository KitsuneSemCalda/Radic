// Package parser implements a recursive-descent, precedence-climbing
// parser that turns a stream of tokens produced by the lexer into a
// Radic abstract syntax tree (see radic/internal/ast).
package parser

import (
	"fmt"

	"radic/internal/ast"
	"radic/internal/token"
)

// Parser walks a fixed slice of tokens left-to-right, with one token of
// lookahead, building an AST as it goes.
type Parser struct {
	src []token.Token
	pos int
}

// Parse parses src as a complete Radic program and returns its AST, or
// an error describing the first syntax error encountered.
func Parse(src []token.Token) (*ast.Program, error) {
	p := &Parser{src: src}

	return p.parseProgram()
}

// isTypeToken reports whether k is a builtin type keyword: number,
// string, char, float, bool, error, void, or a pinned width (i8..u64).
func isTypeToken(k token.TokenKind) bool {
	switch k {
	case token.TokenNumber, token.TokenString, token.TokenChar,
		token.TokenFloat, token.TokenBool, token.TokenError, token.TokenVoid,
		token.TokenI8, token.TokenI16, token.TokenI32, token.TokenI64,
		token.TokenU8, token.TokenU16, token.TokenU32, token.TokenU64:
		return true
	}
	return false
}

// isAssignOp reports whether k is a plain or compound assignment operator.
func isAssignOp(k token.TokenKind) bool {
	return k == token.TokenAssign || k.IsCompoundOperator()
}

// isAtEnd reports whether the parser has consumed every token up to EOF.
func (p *Parser) isAtEnd() bool { return p.peek().TokenType == token.TokenEOF }

// peek returns the next token to be consumed without advancing the
// parser. It returns a synthetic EOF token once the input is exhausted.
func (p *Parser) peek() token.Token {
	if p.pos >= len(p.src) {
		return token.Token{TokenType: token.TokenEOF}
	}
	return p.src[p.pos]
}

// advance consumes and returns the next token, or a synthetic EOF token
// once the input is exhausted.
func (p *Parser) advance() token.Token {
	t := p.peek()
	if p.pos < len(p.src) {
		p.pos++
	}
	return t
}

// match consumes and returns true if the next token's kind is one of
// kinds; otherwise it leaves the parser position unchanged and returns
// false.
func (p *Parser) match(kinds ...token.TokenKind) bool {
	for _, k := range kinds {
		if p.peek().TokenType == k {
			p.advance()
			return true
		}
	}
	return false
}

// expect consumes the next token if it has kind k, returning an error
// describing the mismatch otherwise.
func (p *Parser) expect(k token.TokenKind) (token.Token, error) {
	t := p.peek()
	if t.TokenType != k {
		return t, fmt.Errorf("expected %s, got %q (%s)", k, t.Lexeme, t.TokenType)
	}
	return p.advance(), nil
}

// errorf builds a syntax error for token t, formatting msg the same way
// as fmt.Sprintf.
func (p *Parser) errorf(t token.Token, format string, args ...any) error {
	return fmt.Errorf("%s at %q (%s)", fmt.Sprintf(format, args...), t.Lexeme, t.TokenType)
}

// nextIsIdentifier reports whether the token after the current one is an
// identifier, used to disambiguate a user-defined type name from an
// expression that happens to start with the same identifier.
func (p *Parser) nextIsIdentifier() bool {
	n := p.pos + 1
	return n < len(p.src) && p.src[n].TokenType == token.TokenIdentifier
}

// isDeclStart reports whether the parser is positioned at the start of a
// variable declaration: a type name, builtin or user-defined, followed
// by another identifier. The trailing identifier check is required even
// for builtin types because the lexer reuses their token kind for same-
// named literals (e.g. TokenNumber is both the "number" type keyword and
// a numeric literal like "42"), which would otherwise be indistinguishable
// from a declaration by kind alone.
func (p *Parser) isDeclStart() bool {
	t := p.peek()
	return (isTypeToken(t.TokenType) || t.TokenType == token.TokenIdentifier) &&
		p.nextIsIdentifier()
}

// parseProgram parses a full source file as a sequence of top-level
// declarations.
func (p *Parser) parseProgram() (*ast.Program, error) {
	prog := &ast.Program{}
	for !p.isAtEnd() {
		d, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}
		prog.Decls = append(prog.Decls, d)
	}
	return prog, nil
}

// parseDeclaration parses a single top-level declaration: a function,
// struct, union, enum, or variable declaration.
func (p *Parser) parseDeclaration() (ast.Decl, error) {
	switch p.peek().TokenType {
	case token.TokenFunc:
		return p.parseFuncDecl()
	case token.TokenStruct:
		return p.parseAggregate(token.TokenStruct)
	case token.TokenUnion:
		return p.parseAggregate(token.TokenUnion)
	case token.TokenEnum:
		return p.parseEnumDecl()
	case token.TokenInclude:
		return nil, p.errorf(p.peek(), "include directives are not supported yet")
	default:
		if !p.isDeclStart() {
			return nil, p.errorf(p.peek(), "expected a declaration")
		}
		return p.parseVarDecl(true)
	}
}

// parseAggregate parses a struct or union declaration; kind selects
// which keyword introduces it.
func (p *Parser) parseAggregate(kind token.TokenKind) (ast.Decl, error) {
	p.advance()
	nameTok, err := p.expect(token.TokenIdentifier)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(token.TokenLeftBrace); err != nil {
		return nil, err
	}
	fields, err := p.parseFields()
	if err != nil {
		return nil, err
	}
	if kind == token.TokenStruct {
		return &ast.StructDecl{Name: nameTok.Lexeme, Fields: fields}, nil
	}
	return &ast.UnionDecl{Name: nameTok.Lexeme, Fields: fields}, nil
}

// parseFields parses the semicolon-terminated field list of a struct or
// union body, up to and including the closing brace.
func (p *Parser) parseFields() ([]ast.Field, error) {
	var fields []ast.Field
	for !p.isAtEnd() && p.peek().TokenType != token.TokenRightBrace {
		typ := p.parseType()
		nameTok, err := p.expect(token.TokenIdentifier)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.TokenSemicolon); err != nil {
			return nil, err
		}
		fields = append(fields, ast.Field{Type: typ, Name: nameTok.Lexeme})
	}
	if _, err := p.expect(token.TokenRightBrace); err != nil {
		return nil, err
	}
	return fields, nil
}

// parseEnumDecl parses an enum declaration and its comma-separated
// variants, each with an optional explicit value.
func (p *Parser) parseEnumDecl() (*ast.EnumDecl, error) {
	p.advance()
	nameTok, err := p.expect(token.TokenIdentifier)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(token.TokenLeftBrace); err != nil {
		return nil, err
	}
	decl := &ast.EnumDecl{Name: nameTok.Lexeme}
	for !p.isAtEnd() && p.peek().TokenType != token.TokenRightBrace {
		vName, err := p.expect(token.TokenIdentifier)
		if err != nil {
			return nil, err
		}
		v := ast.EnumVariant{Name: vName.Lexeme}
		if p.match(token.TokenAssign) {
			val, err := p.parseExpression(1)
			if err != nil {
				return nil, err
			}
			v.Value = val
		}
		decl.Variants = append(decl.Variants, v)
		if !p.match(token.TokenComma) {
			break
		}
	}
	if _, err := p.expect(token.TokenRightBrace); err != nil {
		return nil, err
	}
	return decl, nil
}

// parseFuncDecl parses a function declaration: its name, parameter list,
// optional return type, and body.
func (p *Parser) parseFuncDecl() (*ast.FuncDecl, error) {
	p.advance()
	nameTok, err := p.expect(token.TokenIdentifier)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(token.TokenLeftParen); err != nil {
		return nil, err
	}
	fn := &ast.FuncDecl{Name: nameTok.Lexeme, Return: ast.Name{}}
	for !p.isAtEnd() && p.peek().TokenType != token.TokenRightParen {
		typ := p.parseType()
		param, err := p.expect(token.TokenIdentifier)
		if err != nil {
			return nil, err
		}
		fn.Params = append(fn.Params, ast.Param{Type: typ, Name: param.Lexeme})
		if !p.match(token.TokenComma) {
			break
		}
	}
	if _, err := p.expect(token.TokenRightParen); err != nil {
		return nil, err
	}
	if p.match(token.TokenArrow) {
		fn.Return = p.parseType()
	}
	fn.Body, err = p.parseBlock()
	return fn, err
}

// parseVarDecl parses a variable declaration with its type, name, and
// optional initializer. If semi is true, it also consumes the
// terminating semicolon; callers parsing a for-loop initializer pass
// false since the semicolon there is consumed by the caller.
func (p *Parser) parseVarDecl(semi bool) (*ast.VarDecl, error) {
	typ := p.parseType()
	nameTok, err := p.expect(token.TokenIdentifier)
	if err != nil {
		return nil, err
	}
	d := &ast.VarDecl{Type: typ, Name: nameTok.Lexeme}
	if p.match(token.TokenAssign) {
		init, err := p.parseExpression(1)
		if err != nil {
			return nil, err
		}
		d.Init = init
	}
	if semi {
		if _, err := p.expect(token.TokenSemicolon); err != nil {
			return nil, err
		}
	}
	return d, nil
}

// parseType parses a type name, marking it as Builtin when it names one
// of the primitive types.
func (p *Parser) parseType() ast.Name {
	t := p.advance()
	if isTypeToken(t.TokenType) {
		return ast.Name{Lexeme: t.Lexeme, Builtin: true}
	}
	return ast.Name{Lexeme: t.Lexeme}
}

// parseBlock parses a brace-delimited sequence of statements.
func (p *Parser) parseBlock() (*ast.Block, error) {
	if _, err := p.expect(token.TokenLeftBrace); err != nil {
		return nil, err
	}
	b := &ast.Block{}
	for !p.isAtEnd() && p.peek().TokenType != token.TokenRightBrace {
		s, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		b.Stmts = append(b.Stmts, s)
	}
	if _, err := p.expect(token.TokenRightBrace); err != nil {
		return nil, err
	}
	return b, nil
}

// parseStatement parses a single statement, dispatching on the next
// token's kind.
func (p *Parser) parseStatement() (ast.Stmt, error) {
	switch p.peek().TokenType {
	case token.TokenLeftBrace:
		return p.parseBlock()
	case token.TokenIf:
		return p.parseIf()
	case token.TokenWhile:
		return p.parseWhile()
	case token.TokenFor:
		return p.parseFor()
	case token.TokenSwitch:
		return p.parseSwitch()
	case token.TokenBreak:
		p.advance()
		if _, err := p.expect(token.TokenSemicolon); err != nil {
			return nil, err
		}
		return &ast.Break{}, nil
	case token.TokenContinue:
		p.advance()
		if _, err := p.expect(token.TokenSemicolon); err != nil {
			return nil, err
		}
		return &ast.Continue{}, nil
	case token.TokenReturn:
		return p.parseReturn()
	default:
		if p.isDeclStart() {
			d, err := p.parseVarDecl(true)
			if err != nil {
				return nil, err
			}
			return &ast.VarDeclStmt{VarDecl: *d}, nil
		}
		return p.parseExprOrAssign()
	}
}

// parseIf parses an if statement, including an optional else branch or
// else-if chain.
func (p *Parser) parseIf() (ast.Stmt, error) {
	p.advance()
	if _, err := p.expect(token.TokenLeftParen); err != nil {
		return nil, err
	}
	cond, err := p.parseExpression(1)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(token.TokenRightParen); err != nil {
		return nil, err
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	st := &ast.If{Cond: cond, Then: body}
	if p.match(token.TokenElse) {
		els, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		st.Else = els
	}
	return st, nil
}

// parseWhile parses a while loop.
func (p *Parser) parseWhile() (ast.Stmt, error) {
	p.advance()
	if _, err := p.expect(token.TokenLeftParen); err != nil {
		return nil, err
	}
	cond, err := p.parseExpression(1)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(token.TokenRightParen); err != nil {
		return nil, err
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	return &ast.While{Cond: cond, Body: body}, nil
}

// parseFor parses a C-style for loop, whose init, condition, and post
// clauses are each optional.
func (p *Parser) parseFor() (ast.Stmt, error) {
	p.advance()
	if _, err := p.expect(token.TokenLeftParen); err != nil {
		return nil, err
	}
	f := &ast.For{}
	if !p.match(token.TokenSemicolon) {
		if p.isDeclStart() {
			d, err := p.parseVarDecl(false)
			if err != nil {
				return nil, err
			}
			f.Init = &ast.VarDeclStmt{VarDecl: *d}
		} else {
			init, err := p.parseSimpleStmt()
			if err != nil {
				return nil, err
			}
			f.Init = init
		}
		if _, err := p.expect(token.TokenSemicolon); err != nil {
			return nil, err
		}
	}
	if !p.match(token.TokenSemicolon) {
		cond, err := p.parseExpression(1)
		if err != nil {
			return nil, err
		}
		f.Cond = cond
		if _, err := p.expect(token.TokenSemicolon); err != nil {
			return nil, err
		}
	}
	if !p.match(token.TokenRightParen) {
		post, err := p.parseSimpleStmt()
		if err != nil {
			return nil, err
		}
		f.Post = post
		if _, err := p.expect(token.TokenRightParen); err != nil {
			return nil, err
		}
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	f.Body = body
	return f, nil
}

// parseSwitch parses a switch statement and its case and default clauses.
func (p *Parser) parseSwitch() (ast.Stmt, error) {
	p.advance()
	if _, err := p.expect(token.TokenLeftParen); err != nil {
		return nil, err
	}
	disc, err := p.parseExpression(1)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(token.TokenRightParen); err != nil {
		return nil, err
	}
	if _, err := p.expect(token.TokenLeftBrace); err != nil {
		return nil, err
	}
	sw := &ast.Switch{Disc: disc}
	for !p.isAtEnd() && p.peek().TokenType != token.TokenRightBrace {
		switch p.peek().TokenType {
		case token.TokenCase:
			p.advance()
			val, err := p.parseExpression(1)
			if err != nil {
				return nil, err
			}
			body, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			sw.Cases = append(sw.Cases, ast.CaseClause{Value: val, Stmts: body.Stmts})
		case token.TokenDefault:
			p.advance()
			body, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			sw.Default = body.Stmts
		default:
			return nil, p.errorf(p.peek(), "expected case or default")
		}
	}
	if _, err := p.expect(token.TokenRightBrace); err != nil {
		return nil, err
	}
	return sw, nil
}

// parseReturn parses a return statement, whose value is optional.
func (p *Parser) parseReturn() (ast.Stmt, error) {
	p.advance()
	r := &ast.Return{}
	if !p.match(token.TokenSemicolon) {
		x, err := p.parseExpression(1)
		if err != nil {
			return nil, err
		}
		r.Value = x
		if _, err := p.expect(token.TokenSemicolon); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// parseExprOrAssign parses a statement that starts with an expression,
// which is either a bare expression statement or an assignment to an
// identifier or member expression, consuming the terminating semicolon.
func (p *Parser) parseExprOrAssign() (ast.Stmt, error) {
	s, err := p.parseSimpleStmt()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(token.TokenSemicolon); err != nil {
		return nil, err
	}
	return s, nil
}

// parseSimpleStmt parses a bare expression or an assignment to an
// identifier or member expression, without consuming any terminator.
// It is shared by parseExprOrAssign, which consumes a trailing
// semicolon, and parseFor's post clause, which is instead terminated by
// the loop's closing parenthesis.
func (p *Parser) parseSimpleStmt() (ast.Stmt, error) {
	x, err := p.parseExpression(1)
	if err != nil {
		return nil, err
	}
	if isAssignOp(p.peek().TokenType) {
		op := p.advance()
		var ok bool
		switch x.(type) {
		case *ast.Ident, *ast.Member:
			ok = true
		}
		if !ok {
			return nil, fmt.Errorf("invalid assignment target: %T", x)
		}
		val, err := p.parseExpression(1)
		if err != nil {
			return nil, err
		}
		return &ast.Assign{Target: x, Op: op.TokenType, Value: val}, nil
	}
	return &ast.ExprStmt{X: x}, nil
}

// parseExpression parses a binary expression using precedence climbing,
// only consuming operators whose precedence is at least minPrec.
func (p *Parser) parseExpression(minPrec int) (ast.Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for {
		t := p.peek()
		prec, isBin := binaryPrecedence[t.TokenType]
		if !isBin || prec < minPrec {
			return left, nil
		}
		p.advance()
		right, err := p.parseExpression(prec + 1)
		if err != nil {
			return nil, err
		}
		left = &ast.Binary{Op: t.TokenType, Left: left, Right: right}
	}
}

// parseUnary parses a prefix unary expression, or falls through to a
// postfix expression when the next token is not a unary operator.
func (p *Parser) parseUnary() (ast.Expr, error) {
	t := p.peek()
	switch t.TokenType {
	case token.TokenNot, token.TokenMinus, token.TokenPlus,
		token.TokenBitwiseNot, token.TokenStar, token.TokenBitwiseAnd,
		token.TokenPlusPlus, token.TokenMinusMinus:
		p.advance()
		x, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{Op: t.TokenType, X: x}, nil
	}
	return p.parsePostfix()
}

// parsePostfix parses a primary expression followed by any number of
// calls, member accesses, and postfix increment/decrement operators.
func (p *Parser) parsePostfix() (ast.Expr, error) {
	x, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for {
		switch p.peek().TokenType {
		case token.TokenLeftParen:
			x, err = p.parseCall(x)
			if err != nil {
				return nil, err
			}
		case token.TokenDot, token.TokenArrow:
			op := p.advance()
			nameTok, err := p.expect(token.TokenIdentifier)
			if err != nil {
				return nil, err
			}
			x = &ast.Member{Object: x, Name: nameTok.Lexeme, Deref: op.TokenType == token.TokenArrow}
		case token.TokenPlusPlus, token.TokenMinusMinus:
			op := p.advance()
			x = &ast.Unary{Op: op.TokenType, X: x, Post: true}
		default:
			return x, nil
		}
	}
}

// parseCall parses the parenthesized, comma-separated argument list of a
// call expression whose callee, fn, has already been parsed.
func (p *Parser) parseCall(fn ast.Expr) (ast.Expr, error) {
	p.advance()
	call := &ast.Call{Fn: fn}
	for !p.isAtEnd() && p.peek().TokenType != token.TokenRightParen {
		arg, err := p.parseExpression(1)
		if err != nil {
			return nil, err
		}
		call.Args = append(call.Args, arg)
		if !p.match(token.TokenComma) {
			break
		}
	}
	if _, err := p.expect(token.TokenRightParen); err != nil {
		return nil, err
	}
	return call, nil
}

// parsePrimary parses a literal, identifier, or parenthesized expression.
func (p *Parser) parsePrimary() (ast.Expr, error) {
	t := p.advance()
	switch t.TokenType {
	case token.TokenNumber:
		return &ast.NumberLit{Value: t.Lexeme}, nil
	case token.TokenFloat:
		return &ast.FloatLit{Value: t.Lexeme}, nil
	case token.TokenString:
		return &ast.StringLit{Value: t.Lexeme}, nil
	case token.TokenChar:
		return &ast.CharLit{Value: t.Lexeme}, nil
	case token.TokenIdentifier:
		switch t.Lexeme {
		case "true":
			return &ast.BoolLit{Value: t.Lexeme}, nil
		case "false":
			return &ast.BoolLit{Value: t.Lexeme}, nil
		case "nil":
			return &ast.NilLit{}, nil
		default:
			return &ast.Ident{Name: t.Lexeme}, nil
		}
	case token.TokenLeftParen:
		x, err := p.parseExpression(1)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.TokenRightParen); err != nil {
			return nil, err
		}
		return x, nil
	}
	return nil, p.errorf(t, "expected expression")
}
