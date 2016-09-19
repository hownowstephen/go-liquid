package liquid

import "io"

// Tokenizer allows iteration through a list of tokens
type Tokenizer struct {
	tokens []string
	index  int
}

// Next returns token, if available, and an EOF if the end has been reached
func (t *Tokenizer) Next() (string, error) {

	if t.index >= len(t.tokens) {
		return "", io.EOF
	}

	token := t.tokens[t.index]
	t.index++

	var err error
	if t.index >= len(t.tokens) {
		err = io.EOF
	}
	return token, err
}

// NewTokenizer creates a *Tokenizer instance specific to the supplied template
func NewTokenizer(template string) *Tokenizer {
	indices := templateParserRegexp.FindAllStringIndex(template, -1)

	var tokens []string
	var before int
	for _, loc := range indices {
		if loc[0] > before {
			tokens = append(tokens, template[before:loc[0]])
		}
		tokens = append(tokens, template[loc[0]:loc[1]])
		before = loc[1]
	}

	if before < len(template) {
		tokens = append(tokens, template[before:len(template)])
	}

	return &Tokenizer{
		tokens: tokens,
		index:  0,
	}
}
