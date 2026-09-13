package lexer

import "radic/internal/token"

// keywords maps reserved words recognized by the lexer to their
// corresponding token kind. Any identifier not found in this map is
// scanned as a plain TokenIdentifier.
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
	"i8":     token.TokenI8,
	"i16":    token.TokenI16,
	"i32":    token.TokenI32,
	"i64":    token.TokenI64,
	"u8":     token.TokenU8,
	"u16":    token.TokenU16,
	"u32":    token.TokenU32,
	"u64":    token.TokenU64,
}
