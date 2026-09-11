package token

type Token struct {
	TokenType TokenKind
	Lexeme    string
	Line      int
	Column    int
}
