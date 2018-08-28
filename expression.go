package liquid

import (
	"fmt"
	"regexp"
	"strconv"
)

// Expression objects contain specific types of usable data
// that can be rendered to a string value at runtime.
type Expression interface {
	// Evaluate with the supplied Context
	// XXX: not so certain this is actually necessary, will need to review once
	// more non-primitives have been implemented
	Evaluate(Context) Expression
	// Name might be better as a helper method that switches on type
	Name() string
}

// Define some constant literals
var (
	Nil   = nilExpr{}
	True  = boolExpr(true)
	False = boolExpr(false)
	// XXX: MethodLiteral isn't implemented
	//   'blank'.freeze => MethodLiteral.new(:blank?, '').freeze,
	//   'empty'.freeze => MethodLiteral.new(:empty?, '').freeze

	// Regexes for parsing different types of literals
	singleQuotedStringRegex = regexp.MustCompile(`(?ms)\A'(.*)'\z`)
	doubleQuotedStringRegex = regexp.MustCompile(`(?ms)\A"(.*)"\z`)
	integerRegex            = regexp.MustCompile(`\A(-?\d+)\z`)
	rangeRegex              = regexp.MustCompile(`\A\((\S+)\.\.(\S+)\)\z`)
	floatRegex              = regexp.MustCompile(`\A(-?\d[\d\.]+)\z`)
)

// ParseExpression takes an expression and converts it into something usable by Liquid
// XXX: not sure that renderable is the right approach. maybe just interface{} and then the render
// func will have to know how to deal with it? otherwise we have all these literals
func ParseExpression(markup string) Expression {
	switch markup {
	case "nil", "null", "":
		return Nil
	case "false":
		return False
	case "true":
		return True
		// XXX: implement
		// case "blank", "empty":
		// 	return EmptyLiteral{}
	}

	if singleQuotedStringRegex.MatchString(markup) || doubleQuotedStringRegex.MatchString(markup) {
		return stringExpr(markup[1 : len(markup)-1])
	}

	if integerRegex.MatchString(markup) {
		value, err := strconv.Atoi(markup)
		if err != nil {
			// XXX: this needs a real handler
			panic(err)
		}
		return integerExpr(value)
	}

	if submatch := rangeRegex.FindAllStringSubmatch(markup, 1); len(submatch) > 0 {
		start, err := strconv.Atoi(submatch[0][0])
		if err != nil {
			panic(err)
		}

		end, err := strconv.Atoi(submatch[0][1])
		if err != nil {
			panic(err)
		}
		return rangeExpr{start, end}
	}

	if floatRegex.MatchString(markup) {
		f, err := strconv.ParseFloat(markup, 64)
		if err != nil {
			panic(err)
		}
		return floatExpr(f)
	}

	return ParseVariableLookup(markup)
}

// Base expression types

type nilExpr struct{}

func (e nilExpr) Evaluate(c Context) Expression { return e }

func (e nilExpr) Name() string {
	return "nil"
}

type boolExpr bool

func (e boolExpr) Evaluate(c Context) Expression {
	return e
}

func (e boolExpr) Name() string {
	return fmt.Sprintf("%v", e)
}

type stringExpr string

func (e stringExpr) Evaluate(c Context) Expression {
	return e
}

func (e stringExpr) Name() string {
	return string(e)
}

type integerExpr int

func (e integerExpr) Evaluate(c Context) Expression {
	return e
}

func (e integerExpr) Name() string {
	return strconv.Itoa(int(e))
}

type floatExpr float64

func (e floatExpr) Evaluate(c Context) Expression {
	return e
}

func (e floatExpr) Name() string {
	return strconv.FormatFloat(float64(e), 'f', 2, 64)
}

type rangeExpr struct {
	start int
	end   int
}

func (e rangeExpr) Evaluate(c Context) Expression {
	panic("nope")
}

func (e rangeExpr) Name() string {
	return fmt.Sprintf("%v..%v", e.start, e.end)
}

// literalExpr acts like an atom
type literalExpr string

func (e literalExpr) Evaluate(c Context) Expression {
	return e
}

func (e literalExpr) Name() string {
	return string(e)
}

type arrayExpr []interface{}

func (e arrayExpr) Evaluate(c Context) Expression {
	return e
}

func (e arrayExpr) Name() string {
	return "some array, dunno lol"
}
