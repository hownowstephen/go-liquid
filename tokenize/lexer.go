package tokenize

import (
	"errors"
	"fmt"
	"io"
	"unicode"
)

type lexemeType int

const (
	Int lexemeType = iota
	Float
	Quoted
	Pipe
	Assign
	Equal
	Greater
	GTE
	Less
	LTE
	NotEqual
	Identifier
	Minus
	Plus
	Multiply
	Divide
	OpenSquareBracket
	CloseSquareBracket
	Dot
	Ellipsis
)

type lexeme struct {
	Type  lexemeType
	Value string
}

type Lexer struct {
	data string
	idx  int
}

func NewLexer(data string) *Lexer {
	return &Lexer{data: data}
}

func isRune(expect rune) func(r rune) bool {
	return func(r rune) bool {
		fmt.Println("checking", string(r), "is", string(expect), r == expect)
		return r == expect
	}
}

func (l *Lexer) consumeOne(typ lexemeType) (*lexeme, error) {
	return l.consume(typ, 1)
}

func (l *Lexer) consume(typ lexemeType, count int) (*lexeme, error) {
	if l.idx >= len(l.data) {
		return nil, io.EOF
	}
	l.idx += count
	return &lexeme{Type: typ}, nil
}

func (l *Lexer) Next() (*lexeme, error) {

	if l.idx >= len(l.data) {
		return nil, io.EOF
	}

	r := rune(l.data[l.idx])

	switch r {
	case '|':
		return l.consumeOne(Pipe)
	case '-':
		// XXX: not strictly correct, check how other lexers handle this
		if l.checkNext(unicode.IsSpace) {
			return l.consumeOne(Minus)
		}
	case '+':
		return l.consumeOne(Plus)
	case '*':
		return l.consumeOne(Multiply)
	case '/':
		return l.consumeOne(Divide)
	case '=':
		if l.checkNext(isRune('=')) {
			return l.consume(Equal, 2)
		}
		return l.consumeOne(Assign)
	case '>':
		if l.checkNext(isRune('=')) {
			return l.consume(GTE, 2)
		}
		return l.consumeOne(Greater)
	case '<':
		if l.checkNext(isRune('=')) {
			return l.consume(LTE, 2)
		}
		return l.consumeOne(Less)
	case '[':
		return l.consumeOne(OpenSquareBracket)
	case ']':
		return l.consumeOne(CloseSquareBracket)
	case '!':
		if l.checkNext(isRune('=')) {
			return l.consume(NotEqual, 2)
		}
		return nil, errors.New("invalid token '!'")
	case '.':
		if l.checkNext(isRune('.')) {
			return l.consume(Ellipsis, 2)
		}
		return l.consumeOne(Dot)
	}

	switch {
	case unicode.IsNumber(r) || r == '-' && l.checkNext(unicode.IsNumber):
		return l.consumeNumeric()
	case unicode.IsSpace(r):
		for ; l.idx < len(l.data); l.idx++ {
			if !unicode.IsSpace(rune(l.data[l.idx])) {
				break
			}
		}
		return l.Next()
	case r == '"' || r == '\'':
		return l.consumeQuoted(r)
	default:
		return l.consumeExpression()
	}
}

func (l *Lexer) checkNext(checker func(r rune) bool) bool {
	return l.checkIdx(l.idx+1, checker)
}

func (l *Lexer) checkIdx(idx int, checker func(r rune) bool) bool {
	if idx >= len(l.data) {
		return false
	}
	return checker(rune(l.data[idx]))
}

func (l *Lexer) consumeQuoted(quot rune) (*lexeme, error) {
	l.idx++
	for idx, b := range l.data[l.idx:] {
		if b == quot {
			slice := l.data[l.idx : l.idx+idx]
			l.idx += idx + 1
			return &lexeme{Type: Quoted, Value: slice}, nil
		}
	}
	return nil, errors.New("unterminated quote")
}

func (l *Lexer) consumeNumeric() (*lexeme, error) {
	var hasDecimal bool
	for idx, b := range l.data[l.idx:] {
		if unicode.IsNumber(rune(b)) {
			continue
		} else if idx == 0 && b == '-' {
			continue
		} else if !hasDecimal && b == '.' && l.checkIdx(idx+1, unicode.IsNumber) {
			hasDecimal = true
			continue
		}

		slice := l.data[l.idx : idx+l.idx]
		l.idx += idx

		typ := Int
		if hasDecimal {
			typ = Float
		}

		return &lexeme{Type: typ, Value: slice}, nil
	}
	return nil, io.EOF
}

func (l *Lexer) consumeExpression() (*lexeme, error) {
	for idx, b := range l.data[l.idx:] {
		if unicode.IsLetter(rune(b)) || b == '_' {
			continue
		}
		slice := l.data[l.idx : idx+l.idx]
		l.idx += idx
		return &lexeme{Type: Identifier, Value: slice}, nil
	}

	slice := l.data[l.idx:]
	l.idx = len(l.data)

	return &lexeme{Type: Identifier, Value: slice}, nil
}

// func LexTag(data string) (*LiquidTag, error) {

// }
