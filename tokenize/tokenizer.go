package tokenize

import (
	"fmt"
	"io"
)

type Tokenizer struct {
	data      string
	idx       int
	line      int
	lineStart int
}

func NewTokenizer(data string) *Tokenizer {
	return &Tokenizer{data: data}
}

func (t *Tokenizer) Next() (*Token, error) {

	typ := String

	if t.idx < len(t.data) {
		switch t.data[t.idx] {
		case '{':
			switch t.data[t.idx+1] {
			case '{':
				typ = Tag
			case '%':
				typ = Block
			}
		}
	}

	var line = t.line + 1
	for idx := t.idx; idx < len(t.data); idx++ {

		if t.data[idx] == '\n' {
			t.line++
			t.lineStart = idx
		}

		if quot := t.data[idx]; quot == '\'' || quot == '"' {
			start := idx
			for idx = idx + 1; idx < len(t.data)-1 && t.data[idx] != quot; idx++ {
			}
			if t.data[idx] != quot {
				return nil, fmt.Errorf("unterminated quote on line %d col %d", line, start-t.lineStart)
			}
		}

		if typ == Block && t.data[idx] == '%' && t.data[idx+1] == '}' {
			slice := t.data[t.idx : idx+2]
			t.idx = idx + 2
			return &Token{Type: typ, line: line, data: slice}, nil
		} else if typ == Tag && t.data[idx] == '}' && t.data[idx+1] == '}' {
			slice := t.data[t.idx : idx+2]
			t.idx = idx + 2
			return &Token{Type: typ, line: line, data: slice}, nil
		} else if typ == String && t.data[idx] == '{' && (t.data[idx+1] == '{' || t.data[idx+1] == '%') {
			slice := t.data[t.idx:idx]
			t.idx = idx
			return &Token{Type: typ, line: line, data: slice}, nil
		}
	}

	if t.idx < len(t.data) {
		if typ != String {
			return nil, fmt.Errorf("unterminated token on line %d col %d", line, t.idx-t.lineStart)
		}
		slice := t.data[t.idx:]
		t.idx = len(t.data)
		return &Token{Type: typ, line: t.line, data: slice}, nil
	}
	return nil, io.EOF
}
