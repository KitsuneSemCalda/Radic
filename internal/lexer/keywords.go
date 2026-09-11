package lexer

import "radic/internal/token"

var keywords = map[string]token.TokenKind{
	"if":       token.TokenIf,
	"else":     token.TokenElse,
	"switch":   token.TokenSwitch,
	"case":     token.TokenCase,
	"default":  token.TokenDefault,
	"break":    token.TokenBreak,
	"continue": token.TokenContinue,
	"while":    token.TokenWhile,
	"for":      token.TokenFor,
	"func":     token.TokenFunc,
	"return":   token.TokenReturn,
	"struct":   token.TokenStruct,
	"enum":     token.TokenEnum,
	"union":    token.TokenUnion,
	"include":  token.TokenInclude,

	"and": token.TokenAnd,
	"or":  token.TokenOr,
	"not": token.TokenNot,

	"number": token.TokenNumber,
	"string": token.TokenString,
	"char":   token.TokenChar,
	"float":  token.TokenFloat,
	"bool":   token.TokenBool,
	"error":  token.TokenError,
	"void":   token.TokenVoid,
}
