package token

type TokenKind int

func (tk TokenKind) String() string {
	var names = map[TokenKind]string{
		TokenInvalid:           "INVALID",
		TokenEOF:               "EOF",
		TokenIdentifier:        "IDENTIFIER",
		TokenNumber:            "NUMBER",
		TokenString:            "STRING",
		TokenChar:              "CHAR",
		TokenFloat:             "FLOAT",
		TokenBool:              "BOOL",
		TokenVoid:              "VOID",
		TokenError:             "ERROR",
		TokenPlus:              "PLUS",
		TokenMinus:             "MINUS",
		TokenStar:              "STAR",
		TokenSlash:             "SLASH",
		TokenPercent:           "PERCENT",
		TokenEquals:            "EQUALS",
		TokenNotEquals:         "NOT_EQUALS",
		TokenLessThan:          "LESS_THAN",
		TokenLessThanEquals:    "LESS_THAN_EQUALS",
		TokenGreaterThan:       "GREATER_THAN",
		TokenGreaterThanEquals: "GREATER_THAN_EQUALS",
		TokenAnd:               "AND",
		TokenNot:               "NOT",
		TokenOr:                "OR",
		TokenAssign:            "ASSIGN",
		TokenLeftParen:         "LEFT_PAREN",
		TokenLeftBracket:       "LEFT_BRACKET",
		TokenLeftBrace:         "LEFT_BRACE",
		TokenRightParen:        "RIGHT_PAREN",
		TokenRightBrace:        "RIGHT_BRACE",
		TokenRightBracket:      "RIGHT_BRACKET",
		TokenComma:             "COMMA",
		TokenSemicolon:         "SEMICOLON",
		TokenDot:               "DOT",
		TokenArrow:             "ARROW",
		TokenPlusEquals:        "PLUS_EQUALS",
		TokenMinusEquals:       "MINUS_EQUALS",
		TokenStarEquals:        "STAR_EQUALS",
		TokenSlashEquals:       "SLASH_EQUALS",
		TokenPercentEquals:     "PERCENT_EQUALS",
		TokenPlusPlus:          "PLUS_PLUS",
		TokenMinusMinus:        "MINUS_MINUS",

		TokenBitwiseAnd:        "BITWISE_AND",
		TokenBitwiseOr:         "BITWISE_OR",
		TokenBitwiseXor:        "BITWISE_XOR",
		TokenBitwiseNot:        "BITWISE_NOT",
		TokenBitwiseLeftShift:  "BITWISE_LEFT_SHIFT",
		TokenBitwiseRightShift: "BITWISE_RIGHT_SHIFT",
		TokenIf:                "IF",
		TokenElse:              "ELSE",
		TokenSwitch:            "SWITCH",
		TokenCase:              "CASE",
		TokenDefault:           "DEFAULT",
		TokenBreak:             "BREAK",
		TokenContinue:          "CONTINUE",
		TokenWhile:             "WHILE",
		TokenFor:               "FOR",
		TokenFunc:              "FUNC",
		TokenReturn:            "RETURN",
		TokenStruct:            "STRUCT",
		TokenEnum:              "ENUM",
		TokenUnion:             "UNION",
		TokenInclude:           "INCLUDE",
	}

	if name, ok := names[tk]; ok {
		return name
	}
	return "UNKNOWN"
}

// IsOperator returns true if the token is an arithmetic, logical, bitwise, or comparison operator
func (tk TokenKind) IsOperator() bool {
	return tk.IsArithmeticOperator() || tk.IsLogicalOperator() || tk.IsComparisonOperator() || tk.IsBitwiseOperator()
}

// IsArithmeticOperator returns true if the token is an arithmetic operator
func (tk TokenKind) IsArithmeticOperator() bool {
	switch tk {
	case TokenPlus, TokenMinus, TokenStar, TokenSlash, TokenPercent:
		return true
	}
	return false
}

// IsComparisonOperator returns true if the token is a comparison operator
func (tk TokenKind) IsComparisonOperator() bool {
	switch tk {
	case TokenEquals, TokenNotEquals, TokenLessThan, TokenLessThanEquals,
		TokenGreaterThan, TokenGreaterThanEquals:
		return true
	}
	return false
}

// IsLogicalOperator returns true if the token is a logical operator
func (tk TokenKind) IsLogicalOperator() bool {
	switch tk {
	case TokenAnd, TokenOr, TokenNot:
		return true
	}
	return false
}

// IsBitwiseOperator returns true if the token is a bitwise operator
func (tk TokenKind) IsBitwiseOperator() bool {
	switch tk {
	case TokenBitwiseAnd, TokenBitwiseOr, TokenBitwiseXor, TokenBitwiseNot,
		TokenBitwiseLeftShift, TokenBitwiseRightShift:
		return true
	}
	return false
}

// IsKeyword returns true if the token is a keyword
func (tk TokenKind) IsKeyword() bool {
	switch tk {
	case TokenIf, TokenElse, TokenSwitch, TokenCase, TokenDefault, TokenBreak, TokenContinue,
		TokenWhile, TokenFor, TokenFunc, TokenReturn, TokenStruct, TokenEnum, TokenUnion, TokenInclude:
		return true
	}
	return false
}

// IsLiteral returns true if the token is a literal value
func (tk TokenKind) IsLiteral() bool {
	switch tk {
	case TokenNumber, TokenString, TokenChar, TokenFloat, TokenBool, TokenVoid:
		return true
	}
	return false
}

// IsPunctuation returns true if the token is a punctuation mark or delimiter
func (tk TokenKind) IsPunctuation() bool {
	switch tk {
	case TokenLeftParen, TokenRightParen, TokenLeftBracket, TokenRightBracket,
		TokenLeftBrace, TokenRightBrace, TokenComma, TokenSemicolon, TokenDot, TokenArrow:
		return true
	}
	return false
}

// IsCompoundOperator returns true if the token is a compound assignment operator
func (tk TokenKind) IsCompoundOperator() bool {
	switch tk {
	case TokenPlusEquals, TokenMinusEquals, TokenStarEquals, TokenSlashEquals, TokenPercentEquals:
		return true
	}
	return false
}

// IsIncrementDecrement returns true if the token is ++ or --
func (tk TokenKind) IsIncrementDecrement() bool {
	return tk == TokenPlusPlus || tk == TokenMinusMinus
}

// IsOpenDelimiter returns true if the token is an opening bracket/parenthesis/brace
func (tk TokenKind) IsOpenDelimiter() bool {
	switch tk {
	case TokenLeftParen, TokenLeftBracket, TokenLeftBrace:
		return true
	}
	return false
}

// IsCloseDelimiter returns true if the token is a closing bracket/parenthesis/brace
func (tk TokenKind) IsCloseDelimiter() bool {
	switch tk {
	case TokenRightParen, TokenRightBracket, TokenRightBrace:
		return true
	}
	return false
}

// IsUnaryOperator returns true if the token can be used as a unary operator
func (tk TokenKind) IsUnaryOperator() bool {
	switch tk {
	case TokenNot, TokenMinus, TokenPlus, TokenStar, TokenBitwiseAnd, TokenBitwiseNot:
		return true
	}
	return false
}

// IsBinaryOperator returns true if the token is a binary operator
func (tk TokenKind) IsBinaryOperator() bool {
	switch tk {
	case TokenPlus, TokenMinus, TokenStar, TokenSlash, TokenPercent: // arithmetic
		return true
	case TokenEquals, TokenNotEquals, TokenLessThan, TokenLessThanEquals,
		TokenGreaterThan, TokenGreaterThanEquals: // comparison
		return true
	case TokenAnd, TokenOr: // logical (binary only - NOT is unary)
		return true
	case TokenBitwiseAnd, TokenBitwiseOr, TokenBitwiseXor,
		TokenBitwiseLeftShift, TokenBitwiseRightShift: // bitwise (binary only - NOT is unary)
		return true
	}
	return false
}

// IsControlFlow returns true if the token is a control flow keyword
func (tk TokenKind) IsControlFlow() bool {
	switch tk {
	case TokenIf, TokenElse, TokenSwitch, TokenBreak, TokenContinue, TokenWhile, TokenFor, TokenReturn:
		return true
	}
	return false
}

const (
	TokenInvalid TokenKind = iota
	TokenEOF

	TokenIdentifier

	// Language Types
	TokenNumber
	TokenString
	TokenChar
	TokenFloat
	TokenBool
	TokenError
	TokenVoid

	// Algebraic Expression
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenPercent

	// Equality & Comparison
	TokenEquals
	TokenNotEquals
	TokenLessThan
	TokenLessThanEquals
	TokenGreaterThan
	TokenGreaterThanEquals

	// Logical Operators
	TokenAnd
	TokenNot
	TokenOr

	// Compound Operators
	TokenPlusEquals
	TokenMinusEquals
	TokenStarEquals
	TokenSlashEquals
	TokenPercentEquals
	TokenPlusPlus
	TokenMinusMinus

	// Assignment
	TokenAssign

	// Punctuation
	TokenLeftParen
	TokenLeftBracket
	TokenLeftBrace
	TokenRightParen
	TokenRightBrace
	TokenRightBracket
	TokenComma
	TokenSemicolon
	TokenDot
	TokenArrow

	// Bitwise Operators
	TokenBitwiseAnd
	TokenBitwiseOr
	TokenBitwiseXor
	TokenBitwiseNot
	TokenBitwiseLeftShift
	TokenBitwiseRightShift

	// Keywords
	TokenIf
	TokenElse
	TokenSwitch
	TokenCase
	TokenDefault
	TokenBreak
	TokenContinue
	TokenWhile
	TokenFor
	TokenFunc
	TokenReturn
	TokenStruct
	TokenEnum
	TokenUnion
	TokenInclude
)
