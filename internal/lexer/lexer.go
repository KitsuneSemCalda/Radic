// Package lexer implements a lexical scanner that converts Radic source
// code into a stream of tokens consumed by the parser.
package lexer

import (
	"fmt"
	"radic/internal/token"
)

// Lexer scans a Radic source string and produces tokens for the parser.
// It keeps track of the current read position along with the line and
// column of that position, which is used for error reporting.
type Lexer struct {
	source  string
	start   int
	current int
	line    int
	column  int
}

// New creates a Lexer ready to scan source, starting at line 1, column 1.
func New(source string) *Lexer {
	return &Lexer{
		source: source,
		line:   1,
		column: 1,
	}
}

// Tokenize scans the entire source and returns the resulting tokens.
// The terminating EOF token is consumed internally and is not included
// in the returned slice.
func (l *Lexer) Tokenize() []token.Token {
	var tokens []token.Token

	for {
		tok := l.nextToken()

		if tok.TokenType == token.TokenEOF {
			return tokens
		}

		tokens = append(tokens, tok)
	}
}

// nextToken scans and returns the next token in the source, skipping
// whitespace and comments along the way.
func (l *Lexer) nextToken() token.Token {
	for {
		l.skipWhitespace()
		l.start = l.current

		if l.isAtEnd() {
			return l.makeToken(token.TokenEOF)
		}

		c := l.advance()

		// NOTE: This sections is specially destined to handle comments, both single-line and multi-line.
		if c == '/' {
			if l.match('/') {
				for !l.isAtEnd() && l.peek() != '\n' {
					l.advance()
				}

				continue
			}

			if l.match('*') {
				if !l.skipBlockComment() {
					return l.errorToken("Unterminated block comment")
				}

				continue
			}
		}

		switch c {
		case '(':
			return l.makeToken(token.TokenLeftParen)
		case ')':
			return l.makeToken(token.TokenRightParen)
		case '[':
			return l.makeToken(token.TokenLeftBracket)
		case ']':
			return l.makeToken(token.TokenRightBracket)
		case '{':
			return l.makeToken(token.TokenLeftBrace)
		case '}':
			return l.makeToken(token.TokenRightBrace)
		case ',':
			return l.makeToken(token.TokenComma)
		case ';':
			return l.makeToken(token.TokenSemicolon)
		case '.':
			return l.makeToken(token.TokenDot)
		case '+':
			if l.match('=') {
				return l.makeToken(token.TokenPlusEquals)
			}
			if l.match('+') {
				return l.makeToken(token.TokenPlusPlus)
			}
			return l.makeToken(token.TokenPlus)
		case '-':
			if l.match('=') {
				return l.makeToken(token.TokenMinusEquals)
			}
			if l.match('-') {
				return l.makeToken(token.TokenMinusMinus)
			}
			if l.match('>') {
				return l.makeToken(token.TokenArrow)
			}
			return l.makeToken(token.TokenMinus)
		case '*':
			if l.match('=') {
				return l.makeToken(token.TokenStarEquals)
			}
			return l.makeToken(token.TokenStar)
		case '/':
			if l.match('=') {
				return l.makeToken(token.TokenSlashEquals)
			}
			return l.makeToken(token.TokenSlash)
		case '%':
			if l.match('=') {
				return l.makeToken(token.TokenPercentEquals)
			}
			return l.makeToken(token.TokenPercent)
		case '=':
			if l.match('=') {
				return l.makeToken(token.TokenEquals)
			}
			return l.makeToken(token.TokenAssign)
		case '!':
			if l.match('=') {
				return l.makeToken(token.TokenNotEquals)
			}
			return l.errorToken("unexpected character '!'")
		case '<':
			if l.match('=') {
				return l.makeToken(token.TokenLessThanEquals)
			}
			if l.match('<') {
				return l.makeToken(token.TokenBitwiseLeftShift)
			}
			return l.makeToken(token.TokenLessThan)
		case '>':
			if l.match('=') {
				return l.makeToken(token.TokenGreaterThanEquals)
			}
			if l.match('>') {
				return l.makeToken(token.TokenBitwiseRightShift)
			}
			return l.makeToken(token.TokenGreaterThan)
		case '&':
			return l.makeToken(token.TokenBitwiseAnd)
		case '|':
			return l.makeToken(token.TokenBitwiseOr)
		case '^':
			return l.makeToken(token.TokenBitwiseXor)
		case '~':
			return l.makeToken(token.TokenBitwiseNot)
		case '"':
			return l.scanString()
		case '\'':
			return l.scanChar()
		}

		switch {
		case isIdentifierStart(c):
			return l.scanIdentifier()
		case isDigit(c):
			return l.scanNumber()
		}

		return l.errorToken(fmt.Sprintf("unexpected character %q", c))
	}
}

// skipWhitespace advances past spaces, tabs, carriage returns, and newlines.
func (l *Lexer) skipWhitespace() {
	for !l.isAtEnd() {
		switch l.source[l.current] {
		case ' ', '\r', '\t', '\n':
			l.advance()
		default:
			return
		}
	}
}

// skipBlockComment consumes a /* ... */ comment, assuming the opening
// "/*" has already been consumed. It reports false if the source ends
// before the closing "*/" is found.
func (l *Lexer) skipBlockComment() bool {
	for !l.isAtEnd() {
		if l.peek() == '*' && l.peekNext() == '/' {
			l.advance()
			l.advance()
			return true
		}
		l.advance()
	}
	return false
}

// scanIdentifier consumes an identifier, returning the matching keyword
// token if it is a reserved word or TokenIdentifier otherwise.
func (l *Lexer) scanIdentifier() token.Token {
	for !l.isAtEnd() && isIdentifierPart(l.peek()) {
		l.advance()
	}
	if kind, ok := keywords[l.source[l.start:l.current]]; ok {
		return l.makeToken(kind)
	}
	return l.makeToken(token.TokenIdentifier)
}

// scanNumber consumes a numeric literal in decimal, hexadecimal (0x),
// binary (0b), or octal (0o) form, producing TokenFloat when a decimal
// point is present or TokenNumber otherwise.
func (l *Lexer) scanNumber() token.Token {
	if l.source[l.start] == '0' {
		switch l.peek() {
		case 'x', 'X':
			l.advance()
			if !l.skipDigits(isHexDigit) {
				return l.errorToken("invalid hex literal")
			}
			return l.finishNumber(token.TokenNumber)
		case 'b', 'B':
			l.advance()
			if !l.skipDigits(isBinaryDigit) {
				return l.errorToken("invalid binary literal")
			}
			return l.finishNumber(token.TokenNumber)
		case 'o', 'O':
			l.advance()
			if !l.skipDigits(isOctalDigit) {
				return l.errorToken("invalid octal literal")
			}
			return l.finishNumber(token.TokenNumber)
		}
	}

	for isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance()
		for isDigit(l.peek()) {
			l.advance()
		}
		return l.finishNumber(token.TokenFloat)
	}

	return l.finishNumber(token.TokenNumber)
}

// finishNumber guards against a number being immediately followed by
// identifier characters (e.g. "123abc"), returning TokenInvalid in that
// case instead of the given kind.
func (l *Lexer) finishNumber(kind token.TokenKind) token.Token {
	if !l.isAtEnd() && isIdentifierStart(l.peek()) {
		for !l.isAtEnd() && isIdentifierPart(l.peek()) {
			l.advance()
		}
		return l.makeToken(token.TokenInvalid)
	}
	return l.makeToken(kind)
}

// scanString consumes a double-quoted string literal, assuming the
// opening quote has already been consumed. Escaped characters are
// skipped without interpretation.
func (l *Lexer) scanString() token.Token {
	for {
		if l.isAtEnd() {
			return l.errorToken("unterminated string")
		}

		switch l.peek() {
		case '"':
			l.advance()
			return l.makeToken(token.TokenString)
		case '\\':
			l.advance()
			if l.isAtEnd() {
				return l.errorToken("unterminated string")
			}
			l.advance()
		default:
			l.advance()
		}
	}
}

// scanChar consumes a single-quoted character literal, assuming the
// opening quote has already been consumed. It reports an error for
// empty literals, literals containing more than one character, an
// unterminated literal, or a raw newline inside the literal.
func (l *Lexer) scanChar() token.Token {
	count := 0

	for {
		if l.isAtEnd() {
			return l.errorToken("unterminated char literal")
		}

		switch l.peek() {
		case '\'':
			l.advance()
			switch count {
			case 0:
				return l.errorToken("empty char literal")
			case 1:
				return l.makeToken(token.TokenChar)
			default:
				return l.errorToken("char literal must contain a single character")
			}
		case '\\':
			l.advance()
			if l.isAtEnd() {
				return l.errorToken("unterminated char literal")
			}
			l.advance()
			count++
		case '\n':
			return l.errorToken("newline in char literal")
		default:
			l.advance()
			count++
		}
	}
}

// advance consumes and returns the current byte, updating the tracked
// line and column.
func (l *Lexer) advance() byte {
	c := l.source[l.current]
	l.current++
	if c == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return c
}

// peek returns the current byte without consuming it, or 0 at EOF.
func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

// peekNext returns the byte after the current one without consuming
// it, or 0 if that would be out of range.
func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

// match consumes the current byte and returns true if it equals
// expected, otherwise it leaves the position unchanged and returns false.
func (l *Lexer) match(expected byte) bool {
	if l.isAtEnd() || l.source[l.current] != expected {
		return false
	}
	l.advance()
	return true
}

// isAtEnd reports whether the lexer has reached the end of the source.
func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

// skipDigits consumes bytes for which valid returns true and reports
// whether at least one byte was consumed.
func (l *Lexer) skipDigits(valid func(byte) bool) bool {
	n := 0
	for valid(l.peek()) {
		l.advance()
		n++
	}
	return n > 0
}

// makeToken builds a token of the given kind whose lexeme is the source
// text between start and the current position.
func (l *Lexer) makeToken(kind token.TokenKind) token.Token {
	return token.Token{TokenType: kind, Lexeme: l.source[l.start:l.current]}
}

// errorToken builds a TokenError token whose lexeme carries msg
// prefixed with the current line number.
func (l *Lexer) errorToken(msg string) token.Token {
	return token.Token{
		TokenType: token.TokenError,
		Lexeme:    fmt.Sprintf("line %d: %s", l.line, msg),
	}
}

// Character classification predicates used while scanning.
func isDigit(c byte) bool       { return c >= '0' && c <= '9' }
func isHexDigit(c byte) bool    { return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') }
func isBinaryDigit(c byte) bool { return c == '0' || c == '1' }
func isOctalDigit(c byte) bool  { return c >= '0' && c <= '7' }
func isAsciiLetter(c byte) bool { return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') }
func isIdentifierStart(c byte) bool {
	return c == '_' || isAsciiLetter(c)
}
func isIdentifierPart(c byte) bool {
	return isIdentifierStart(c) || isDigit(c)
}
