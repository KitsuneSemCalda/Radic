package parser

import (
	"radic/internal/token"
)

var binaryPrecedence = map[token.TokenKind]int{
	token.TokenOr: 1,

	token.TokenAnd: 2,

	token.TokenEquals:            3,
	token.TokenNotEquals:         3,
	token.TokenLessThan:          3,
	token.TokenLessThanEquals:    3,
	token.TokenGreaterThan:       3,
	token.TokenGreaterThanEquals: 3,

	token.TokenPlus:       4,
	token.TokenMinus:      4,
	token.TokenBitwiseOr:  4,
	token.TokenBitwiseXor: 4,

	token.TokenStar:              5,
	token.TokenSlash:             5,
	token.TokenPercent:           5,
	token.TokenBitwiseAnd:        5,
	token.TokenBitwiseLeftShift:  5,
	token.TokenBitwiseRightShift: 5,
}
