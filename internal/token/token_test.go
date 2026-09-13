package token

import (
	"testing"
)

func TestTokenKindString(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected string
	}{
		{TokenPlus, "PLUS"},
		{TokenMinus, "MINUS"},
		{TokenIdentifier, "IDENTIFIER"},
		{TokenIf, "IF"},
		{TokenEOF, "EOF"},
		{TokenInvalid, "INVALID"},
		{TokenI8, "I8"},
		{TokenI16, "I16"},
		{TokenI32, "I32"},
		{TokenI64, "I64"},
		{TokenU8, "U8"},
		{TokenU16, "U16"},
		{TokenU32, "U32"},
		{TokenU64, "U64"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestIsOperator(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenPlus, true},
		{TokenMinus, true},
		{TokenAnd, true},
		{TokenLessThan, true},
		{TokenBitwiseAnd, true},
		{TokenBitwiseOr, true},
		{TokenBitwiseXor, true},
		{TokenPlusEquals, true},
		{TokenMinusMinus, true},
		{TokenIdentifier, false},
		{TokenIf, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsOperator(); got != tt.expected {
				t.Errorf("IsOperator() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsArithmeticOperator(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenPlus, true},
		{TokenMinus, true},
		{TokenStar, true},
		{TokenSlash, true},
		{TokenPercent, true},
		{TokenAnd, false},
		{TokenIdentifier, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsArithmeticOperator(); got != tt.expected {
				t.Errorf("IsArithmeticOperator() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsComparisonOperator(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenEquals, true},
		{TokenNotEquals, true},
		{TokenLessThan, true},
		{TokenGreaterThan, true},
		{TokenPlus, false},
		{TokenIdentifier, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsComparisonOperator(); got != tt.expected {
				t.Errorf("IsComparisonOperator() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsLogicalOperator(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenAnd, true},
		{TokenOr, true},
		{TokenNot, true},
		{TokenPlus, false},
		{TokenIdentifier, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsLogicalOperator(); got != tt.expected {
				t.Errorf("IsLogicalOperator() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsKeyword(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenIf, true},
		{TokenElse, true},
		{TokenSwitch, true},
		{TokenCase, true},
		{TokenDefault, true},
		{TokenBreak, true},
		{TokenContinue, true},
		{TokenFor, true},
		{TokenFunc, true},
		{TokenReturn, true},
		{TokenStruct, true},
		{TokenEnum, true},
		{TokenUnion, true},
		{TokenInclude, true},
		{TokenIdentifier, false},
		{TokenPlus, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsKeyword(); got != tt.expected {
				t.Errorf("IsKeyword() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsLiteral(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenNumber, true},
		{TokenString, true},
		{TokenChar, true},
		{TokenFloat, true},
		{TokenBool, true},
		{TokenVoid, true},
		{TokenError, true},
		{TokenIdentifier, false},
		{TokenIf, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsLiteral(); got != tt.expected {
				t.Errorf("IsLiteral() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsType(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenNumber, true},
		{TokenString, true},
		{TokenChar, true},
		{TokenFloat, true},
		{TokenBool, true},
		{TokenVoid, true},
		{TokenError, true},
		{TokenI8, true},
		{TokenI16, true},
		{TokenI32, true},
		{TokenI64, true},
		{TokenU8, true},
		{TokenU16, true},
		{TokenU32, true},
		{TokenU64, true},
		{TokenIdentifier, false},
		{TokenIf, false},
		{TokenPlus, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsType(); got != tt.expected {
				t.Errorf("IsType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsPunctuation(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenLeftParen, true},
		{TokenRightParen, true},
		{TokenComma, true},
		{TokenSemicolon, true},
		{TokenDot, true},
		{TokenPlus, false},
		{TokenIdentifier, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsPunctuation(); got != tt.expected {
				t.Errorf("IsPunctuation() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsCompoundOperator(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenPlusEquals, true},
		{TokenMinusEquals, true},
		{TokenStarEquals, true},
		{TokenSlashEquals, true},
		{TokenPercentEquals, true},
		{TokenPlus, false},
		{TokenIdentifier, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsCompoundOperator(); got != tt.expected {
				t.Errorf("IsCompoundOperator() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsIncrementDecrement(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenPlusPlus, true},
		{TokenMinusMinus, true},
		{TokenPlus, false},
		{TokenIdentifier, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsIncrementDecrement(); got != tt.expected {
				t.Errorf("IsIncrementDecrement() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsOpenDelimiter(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenLeftParen, true},
		{TokenLeftBracket, true},
		{TokenLeftBrace, true},
		{TokenRightParen, false},
		{TokenPlus, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsOpenDelimiter(); got != tt.expected {
				t.Errorf("IsOpenDelimiter() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsCloseDelimiter(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenRightParen, true},
		{TokenRightBracket, true},
		{TokenRightBrace, true},
		{TokenLeftParen, false},
		{TokenPlus, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsCloseDelimiter(); got != tt.expected {
				t.Errorf("IsCloseDelimiter() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsUnaryOperator(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenNot, true},
		{TokenMinus, true},
		{TokenPlus, true},
		{TokenStar, true},
		{TokenBitwiseAnd, true},
		{TokenBitwiseNot, true},
		{TokenSlash, false},
		{TokenIdentifier, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsUnaryOperator(); got != tt.expected {
				t.Errorf("IsUnaryOperator() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsBinaryOperator(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenPlus, true},
		{TokenAnd, true},
		{TokenLessThan, true},
		{TokenOr, true},
		{TokenBitwiseAnd, true},
		{TokenBitwiseOr, true},
		{TokenBitwiseXor, true},
		{TokenBitwiseLeftShift, true},
		{TokenBitwiseRightShift, true},
		{TokenBitwiseNot, false},
		{TokenNot, false},
		{TokenIdentifier, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsBinaryOperator(); got != tt.expected {
				t.Errorf("IsBinaryOperator() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsControlFlow(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenIf, true},
		{TokenElse, true},
		{TokenSwitch, true},
		{TokenBreak, true},
		{TokenContinue, true},
		{TokenFor, true},
		{TokenReturn, true},
		{TokenWhile, true},
		{TokenFunc, false},
		{TokenIdentifier, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsControlFlow(); got != tt.expected {
				t.Errorf("IsControlFlow() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsBitwiseOperator(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected bool
	}{
		{TokenBitwiseAnd, true},
		{TokenBitwiseOr, true},
		{TokenBitwiseXor, true},
		{TokenBitwiseNot, true},
		{TokenBitwiseLeftShift, true},
		{TokenBitwiseRightShift, true},
		{TokenPlus, false},
		{TokenIdentifier, false},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			if got := tt.kind.IsBitwiseOperator(); got != tt.expected {
				t.Errorf("IsBitwiseOperator() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTokenStringNewTokens(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected string
	}{
		{TokenBitwiseAnd, "BITWISE_AND"},
		{TokenBitwiseOr, "BITWISE_OR"},
		{TokenBitwiseXor, "BITWISE_XOR"},
		{TokenBitwiseNot, "BITWISE_NOT"},
		{TokenBitwiseLeftShift, "BITWISE_LEFT_SHIFT"},
		{TokenBitwiseRightShift, "BITWISE_RIGHT_SHIFT"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestTokenStringNewKeywords(t *testing.T) {
	tests := []struct {
		kind     TokenKind
		expected string
	}{
		{TokenSwitch, "SWITCH"},
		{TokenCase, "CASE"},
		{TokenDefault, "DEFAULT"},
		{TokenBreak, "BREAK"},
		{TokenContinue, "CONTINUE"},
		{TokenStruct, "STRUCT"},
		{TokenEnum, "ENUM"},
		{TokenUnion, "UNION"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestTokenCreation(t *testing.T) {
	tok := Token{
		TokenType: TokenPlus,
		Lexeme:    "+",
	}

	if tok.TokenType != TokenPlus {
		t.Errorf("TokenType = %v, want %v", tok.TokenType, TokenPlus)
	}

	if tok.Lexeme != "+" {
		t.Errorf("Lexeme = %q, want %q", tok.Lexeme, "+")
	}
}
