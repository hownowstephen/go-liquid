package liquid

// https://github.com/Shopify/liquid/blob/master/lib/liquid/parser.rb

import (
	"errors"
	"fmt"
)

// Parser errors
var (
	ErrIndexOutOfBounds = errors.New("No tokens remaining")
)

// Parser objects parse the contents of liquid tags
type Parser struct {
	tokens []Token
	index  uint64
}

// NewParser generates a parser object for consuming tokens
func NewParser(input string) (*Parser, error) {
	tokens, err := Lexer(input)
	if err != nil {
		return nil, err
	}
	return &Parser{tokens: tokens}, nil
}

func (p *Parser) consume(tType string) (string, error) {
	if int(p.index) >= len(p.tokens) {
		return "", ErrIndexOutOfBounds
	}
	token := p.tokens[p.index]
	if tType != "" && token.name != tType {
		return "", fmt.Errorf("Expected %v but found %v:%v", tType, token.name, token.value)
	}
	p.index++

	return token.value, nil
}

// analog to `consume?` in the ruby
func (p *Parser) tryConsume(tType string) bool {
	if int(p.index) >= len(p.tokens) {
		return false
	}
	token := p.tokens[p.index]
	if tType != "" && token.name != tType {
		return false
	}
	p.index++
	return true
}

func (p *Parser) jump(count uint64) error {
	p.index += count
	if int(p.index) >= len(p.tokens) {
		return ErrIndexOutOfBounds
	}
	return nil
}

func (p *Parser) lookahead(tType string, offset uint64) (bool, error) {
	index := p.index + offset
	if int(index) >= len(p.tokens) {
		return false, ErrIndexOutOfBounds
	}
	token := p.tokens[index]
	return token.name == tType, nil
}

func (p *Parser) expression() (string, error) {
	token := p.tokens[p.index]
	if token.name == tIdentifier {
		return p.variableSignature()
	} else if token.name == tSingleStringLiteral || token.name == tDoubleStringLiteral || token.name == tNumberLiteral {
		return p.consume(token.name)
	} else if token.name == tOpenRound {
		p.consume(token.name)
		first, err := p.expression()
		if err != nil {
			return "", err
		}
		p.consume(tDotDot)
		last, err := p.expression()
		if err != nil {
			return "", err
		}
		_, err = p.consume(tCloseRound)
		return fmt.Sprintf("(%v..%v)", first, last), err
	}
	return "", fmt.Errorf("%v is not a valid expression", token)
}

func (p *Parser) variableSignature() (string, error) {

	result, err := p.consume(tIdentifier)
	if err != nil {
		return "", err
	}

	inSquare, err := p.lookahead(tOpenSquare, 0)

	for inSquare {
		open, err := p.consume("")
		if err != nil {
			return "", err
		}
		result += open
		expr, err := p.expression()
		if err != nil {
			return "", err
		}
		result += expr
		close, err := p.consume(tCloseSquare)
		if err != nil {
			fmt.Println("ALMOST, got:", result)
			return "", err
		}
		result += close

		inSquare, err = p.lookahead(tOpenSquare, 0)
		if err != nil {
			return "", err
		}
	}

	dotNext, err := p.lookahead(tDot, 0)
	if err != nil {
		return "", err
	}
	if dotNext {
		value, err := p.consume(tDot)
		if err != nil {
			return "", err
		}
		result += value
		vsig, err := p.variableSignature()
		if err != nil {
			return "", err
		}
		result += vsig
	}

	return result, nil
}

func (p *Parser) argument() (string, error) {

	var arg string

	idNext, err := p.lookahead(tIdentifier, 0)
	if err != nil {
		return "", err
	}
	colonAfter, err := p.lookahead(tColon, 1)
	if err != nil {
		return "", err
	}

	// check for a keyword argument (identifier: expression)
	if idNext && colonAfter {
		keyword, _ := p.consume(tIdentifier)
		arg += keyword
		colon, _ := p.consume(tColon)
		arg += colon
		arg += " "
	}

	expr, err := p.expression()
	return arg + expr, err

	// def argument
	//   str = ""
	//   # might be a keyword argument (identifier: expression)
	//   if look(:id) && look(:colon, 1)
	//     str << consume << consume << ' '.freeze
	//   end

	//   str << expression
	//   str
	// end
}
