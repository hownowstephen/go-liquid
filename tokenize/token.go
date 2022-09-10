package tokenize

type TokenType int

const (
	String TokenType = 0
	Tag    TokenType = 1
	Block  TokenType = 2
)

type Token struct {
	Type TokenType
	line int
	data string
}
