package lexer

import (
	"strings"
	"testing"

	"radic/internal/token"
)

func TestNew(t *testing.T) {
	l := New("abc")

	if l.source != "abc" {
		t.Errorf("source = %q, want %q", l.source, "abc")
	}
	if l.start != 0 {
		t.Errorf("start = %d, want %d", l.start, 0)
	}
	if l.current != 0 {
		t.Errorf("current = %d, want %d", l.current, 0)
	}
	if l.line != 1 {
		t.Errorf("line = %d, want %d", l.line, 1)
	}
	if l.column != 1 {
		t.Errorf("column = %d, want %d", l.column, 1)
	}
}

func TestTokenizeEmptyInput(t *testing.T) {
	tokens := New("").Tokenize()
	if len(tokens) != 0 {
		t.Errorf("Tokenize() = %v, want empty slice", tokens)
	}
}

func TestTokenizeWhitespaceOnly(t *testing.T) {
	tokens := New("   \t\r\n\n  ").Tokenize()
	if len(tokens) != 0 {
		t.Errorf("Tokenize() = %v, want empty slice", tokens)
	}
}

func TestTokenizeDoesNotIncludeEOF(t *testing.T) {
	tokens := New("+").Tokenize()
	for _, tok := range tokens {
		if tok.TokenType == token.TokenEOF {
			t.Errorf("Tokenize() included an EOF token: %v", tokens)
		}
	}
}

func TestTokenizeSingleCharTokens(t *testing.T) {
	tests := []struct {
		source   string
		expected token.TokenKind
	}{
		{"(", token.TokenLeftParen},
		{")", token.TokenRightParen},
		{"[", token.TokenLeftBracket},
		{"]", token.TokenRightBracket},
		{"{", token.TokenLeftBrace},
		{"}", token.TokenRightBrace},
		{",", token.TokenComma},
		{";", token.TokenSemicolon},
		{".", token.TokenDot},
		{"+", token.TokenPlus},
		{"-", token.TokenMinus},
		{"*", token.TokenStar},
		{"/", token.TokenSlash},
		{"%", token.TokenPercent},
		{"=", token.TokenAssign},
		{"<", token.TokenLessThan},
		{">", token.TokenGreaterThan},
		{"&", token.TokenBitwiseAnd},
		{"|", token.TokenBitwiseOr},
		{"^", token.TokenBitwiseXor},
		{"~", token.TokenBitwiseNot},
	}

	for _, tt := range tests {
		t.Run(tt.expected.String(), func(t *testing.T) {
			tokens := New(tt.source).Tokenize()
			if len(tokens) != 1 {
				t.Fatalf("Tokenize(%q) produced %d tokens, want 1: %v", tt.source, len(tokens), tokens)
			}
			if tokens[0].TokenType != tt.expected {
				t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, tt.expected)
			}
			if tokens[0].Lexeme != tt.source {
				t.Errorf("Lexeme = %q, want %q", tokens[0].Lexeme, tt.source)
			}
		})
	}
}

func TestTokenizeCompoundOperators(t *testing.T) {
	tests := []struct {
		source   string
		expected token.TokenKind
	}{
		{"+=", token.TokenPlusEquals},
		{"++", token.TokenPlusPlus},
		{"-=", token.TokenMinusEquals},
		{"--", token.TokenMinusMinus},
		{"->", token.TokenArrow},
		{"*=", token.TokenStarEquals},
		{"/=", token.TokenSlashEquals},
		{"%=", token.TokenPercentEquals},
		{"==", token.TokenEquals},
		{"!=", token.TokenNotEquals},
		{"<=", token.TokenLessThanEquals},
		{"<<", token.TokenBitwiseLeftShift},
		{">=", token.TokenGreaterThanEquals},
		{">>", token.TokenBitwiseRightShift},
	}

	for _, tt := range tests {
		t.Run(tt.expected.String(), func(t *testing.T) {
			tokens := New(tt.source).Tokenize()
			if len(tokens) != 1 {
				t.Fatalf("Tokenize(%q) produced %d tokens, want 1: %v", tt.source, len(tokens), tokens)
			}
			if tokens[0].TokenType != tt.expected {
				t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, tt.expected)
			}
			if tokens[0].Lexeme != tt.source {
				t.Errorf("Lexeme = %q, want %q", tokens[0].Lexeme, tt.source)
			}
		})
	}
}

func TestTokenizeIdentifiers(t *testing.T) {
	tests := []string{"x", "_x", "foo_bar", "Foo123", "_", "a1b2c3"}

	for _, source := range tests {
		t.Run(source, func(t *testing.T) {
			tokens := New(source).Tokenize()
			if len(tokens) != 1 {
				t.Fatalf("Tokenize(%q) produced %d tokens, want 1: %v", source, len(tokens), tokens)
			}
			if tokens[0].TokenType != token.TokenIdentifier {
				t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, token.TokenIdentifier)
			}
			if tokens[0].Lexeme != source {
				t.Errorf("Lexeme = %q, want %q", tokens[0].Lexeme, source)
			}
		})
	}
}

func TestTokenizeKeywords(t *testing.T) {
	for word, kind := range keywords {
		t.Run(word, func(t *testing.T) {
			tokens := New(word).Tokenize()
			if len(tokens) != 1 {
				t.Fatalf("Tokenize(%q) produced %d tokens, want 1: %v", word, len(tokens), tokens)
			}
			if tokens[0].TokenType != kind {
				t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, kind)
			}
			if tokens[0].Lexeme != word {
				t.Errorf("Lexeme = %q, want %q", tokens[0].Lexeme, word)
			}
		})
	}
}

func TestTokenizeNumbers(t *testing.T) {
	tests := []struct {
		source   string
		expected token.TokenKind
	}{
		{"0", token.TokenNumber},
		{"123", token.TokenNumber},
		{"3.14", token.TokenFloat},
		{"0.5", token.TokenFloat},
		{"0x1A", token.TokenNumber},
		{"0X1a", token.TokenNumber},
		{"0b101", token.TokenNumber},
		{"0B110", token.TokenNumber},
		{"0o17", token.TokenNumber},
		{"0O17", token.TokenNumber},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			tokens := New(tt.source).Tokenize()
			if len(tokens) != 1 {
				t.Fatalf("Tokenize(%q) produced %d tokens, want 1: %v", tt.source, len(tokens), tokens)
			}
			if tokens[0].TokenType != tt.expected {
				t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, tt.expected)
			}
			if tokens[0].Lexeme != tt.source {
				t.Errorf("Lexeme = %q, want %q", tokens[0].Lexeme, tt.source)
			}
		})
	}
}

func TestTokenizeInvalidNumbers(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{"int followed by letters", "123abc"},
		{"float followed by letters", "3.14abc"},
		{"hex followed by letters", "0x1Ag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := New(tt.source).Tokenize()
			if len(tokens) != 1 {
				t.Fatalf("Tokenize(%q) produced %d tokens, want 1: %v", tt.source, len(tokens), tokens)
			}
			if tokens[0].TokenType != token.TokenInvalid {
				t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, token.TokenInvalid)
			}
		})
	}
}

func TestTokenizeMalformedRadixLiterals(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{"empty hex literal", "0x"},
		{"empty binary literal", "0b"},
		{"empty octal literal", "0o"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := New(tt.source).Tokenize()
			if len(tokens) != 1 {
				t.Fatalf("Tokenize(%q) produced %d tokens, want 1: %v", tt.source, len(tokens), tokens)
			}
			if tokens[0].TokenType != token.TokenError {
				t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, token.TokenError)
			}
		})
	}
}

func TestTokenizeStrings(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{"simple string", `"hello"`},
		{"empty string", `""`},
		{"string with escaped quote", `"a\"b"`},
		{"string with escaped backslash", `"a\\b"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := New(tt.source).Tokenize()
			if len(tokens) != 1 {
				t.Fatalf("Tokenize(%q) produced %d tokens, want 1: %v", tt.source, len(tokens), tokens)
			}
			if tokens[0].TokenType != token.TokenString {
				t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, token.TokenString)
			}
			if tokens[0].Lexeme != tt.source {
				t.Errorf("Lexeme = %q, want %q", tokens[0].Lexeme, tt.source)
			}
		})
	}
}

func TestTokenizeUnterminatedString(t *testing.T) {
	tests := []string{`"hello`, `"hello\`}

	for _, source := range tests {
		t.Run(source, func(t *testing.T) {
			tokens := New(source).Tokenize()
			if len(tokens) != 1 {
				t.Fatalf("Tokenize(%q) produced %d tokens, want 1: %v", source, len(tokens), tokens)
			}
			if tokens[0].TokenType != token.TokenError {
				t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, token.TokenError)
			}
			if !strings.Contains(tokens[0].Lexeme, "unterminated string") {
				t.Errorf("Lexeme = %q, want it to contain %q", tokens[0].Lexeme, "unterminated string")
			}
		})
	}
}

func TestTokenizeChars(t *testing.T) {
	tests := []string{`'a'`, `'\n'`, `'\''`}

	for _, source := range tests {
		t.Run(source, func(t *testing.T) {
			tokens := New(source).Tokenize()
			if len(tokens) != 1 {
				t.Fatalf("Tokenize(%q) produced %d tokens, want 1: %v", source, len(tokens), tokens)
			}
			if tokens[0].TokenType != token.TokenChar {
				t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, token.TokenChar)
			}
			if tokens[0].Lexeme != source {
				t.Errorf("Lexeme = %q, want %q", tokens[0].Lexeme, source)
			}
		})
	}
}

func TestTokenizeInvalidChars(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantMsg string
	}{
		{"empty char literal", `''`, "empty char literal"},
		{"multi-character literal", `'ab'`, "char literal must contain a single character"},
		{"unterminated char literal", `'a`, "unterminated char literal"},
		// The lexer does not consume the offending newline itself, so a
		// second, unterminated literal follows from the trailing quote.
		{"newline in char literal", "'\n'", "newline in char literal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := New(tt.source).Tokenize()
			if len(tokens) == 0 {
				t.Fatalf("Tokenize(%q) produced no tokens", tt.source)
			}
			if tokens[0].TokenType != token.TokenError {
				t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, token.TokenError)
			}
			if !strings.Contains(tokens[0].Lexeme, tt.wantMsg) {
				t.Errorf("Lexeme = %q, want it to contain %q", tokens[0].Lexeme, tt.wantMsg)
			}
		})
	}
}

func TestTokenizeLineComment(t *testing.T) {
	tokens := New("1 // this is a comment\n2").Tokenize()

	if len(tokens) != 2 {
		t.Fatalf("Tokenize() produced %d tokens, want 2: %v", len(tokens), tokens)
	}
	if tokens[0].Lexeme != "1" || tokens[1].Lexeme != "2" {
		t.Errorf("Tokenize() = %v, want lexemes [1 2]", tokens)
	}
}

func TestTokenizeBlockComment(t *testing.T) {
	tokens := New("1 /* this\nis a comment */ 2").Tokenize()

	if len(tokens) != 2 {
		t.Fatalf("Tokenize() produced %d tokens, want 2: %v", len(tokens), tokens)
	}
	if tokens[0].Lexeme != "1" || tokens[1].Lexeme != "2" {
		t.Errorf("Tokenize() = %v, want lexemes [1 2]", tokens)
	}
}

func TestTokenizeUnterminatedBlockComment(t *testing.T) {
	tokens := New("/* never closed").Tokenize()

	if len(tokens) != 1 {
		t.Fatalf("Tokenize() produced %d tokens, want 1: %v", len(tokens), tokens)
	}
	if tokens[0].TokenType != token.TokenError {
		t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, token.TokenError)
	}
	if !strings.Contains(tokens[0].Lexeme, "Unterminated block comment") {
		t.Errorf("Lexeme = %q, want it to contain %q", tokens[0].Lexeme, "Unterminated block comment")
	}
}

func TestTokenizeBangWithoutEquals(t *testing.T) {
	tokens := New("!").Tokenize()

	if len(tokens) != 1 {
		t.Fatalf("Tokenize() produced %d tokens, want 1: %v", len(tokens), tokens)
	}
	if tokens[0].TokenType != token.TokenError {
		t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, token.TokenError)
	}
	if !strings.Contains(tokens[0].Lexeme, "unexpected character '!'") {
		t.Errorf("Lexeme = %q, want it to contain %q", tokens[0].Lexeme, "unexpected character '!'")
	}
}

func TestTokenizeUnexpectedCharacter(t *testing.T) {
	tokens := New("@").Tokenize()

	if len(tokens) != 1 {
		t.Fatalf("Tokenize() produced %d tokens, want 1: %v", len(tokens), tokens)
	}
	if tokens[0].TokenType != token.TokenError {
		t.Errorf("TokenType = %v, want %v", tokens[0].TokenType, token.TokenError)
	}
	if !strings.Contains(tokens[0].Lexeme, "unexpected character") {
		t.Errorf("Lexeme = %q, want it to contain %q", tokens[0].Lexeme, "unexpected character")
	}
}

func TestErrorTokenReportsLine(t *testing.T) {
	tokens := New("1\n2\n@").Tokenize()

	if len(tokens) != 3 {
		t.Fatalf("Tokenize() produced %d tokens, want 3: %v", len(tokens), tokens)
	}
	if tokens[2].TokenType != token.TokenError {
		t.Fatalf("TokenType = %v, want %v", tokens[2].TokenType, token.TokenError)
	}
	if !strings.Contains(tokens[2].Lexeme, "line 3") {
		t.Errorf("Lexeme = %q, want it to contain %q", tokens[2].Lexeme, "line 3")
	}
}

func TestTokenizeFullProgram(t *testing.T) {
	source := `func add(a, b) {
		return a + b; // sum
	}`

	expected := []token.TokenKind{
		token.TokenFunc,
		token.TokenIdentifier,
		token.TokenLeftParen,
		token.TokenIdentifier,
		token.TokenComma,
		token.TokenIdentifier,
		token.TokenRightParen,
		token.TokenLeftBrace,
		token.TokenReturn,
		token.TokenIdentifier,
		token.TokenPlus,
		token.TokenIdentifier,
		token.TokenSemicolon,
		token.TokenRightBrace,
	}

	tokens := New(source).Tokenize()

	if len(tokens) != len(expected) {
		t.Fatalf("Tokenize() produced %d tokens, want %d: %v", len(tokens), len(expected), tokens)
	}

	for i, want := range expected {
		if tokens[i].TokenType != want {
			t.Errorf("tokens[%d].TokenType = %v, want %v", i, tokens[i].TokenType, want)
		}
	}
}
