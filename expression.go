package liquid

import (
	"fmt"
	"regexp"
	"strconv"
)

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

// Expression objects contain specific types of usable data
// that can be rendered to a string value at runtime.
type Expression interface {
	Render(Vars) string
	// Name might be better as a helper method that switches on type
	Name() string
}

// Base expression types

type nilExpr struct{}

func (e nilExpr) Render(v Vars) string { return "nil" }

func (e nilExpr) Name() string {
	return e.Render(nil)
}

type boolExpr bool

func (e boolExpr) Render(v Vars) string {
	if e {
		return "true"
	}
	return "false"
}

func (e boolExpr) Name() string {
	return e.Render(nil)
}

type stringExpr string

func (e stringExpr) Render(v Vars) string {
	return string(e)
}

func (e stringExpr) Name() string {
	return e.Render(nil)
}

type integerExpr int

func (e integerExpr) Render(v Vars) string {
	return strconv.Itoa(int(e))
}

func (e integerExpr) Name() string {
	return e.Render(nil)
}

type floatExpr float64

// XXX: Implement properly
func (f floatExpr) Render(v Vars) string {
	return "this is a cool float"
}

func (f floatExpr) Name() string {
	return f.Render(nil)
}

type rangeExpr struct {
	start int
	end   int
}

func (e rangeExpr) Render(v Vars) string {
	return fmt.Sprintf("%v..%v", e.start, e.end)
}

func (e rangeExpr) Name() string {
	return e.Render(nil)
}

// literalExpr acts like an atom
type literalExpr string

func (e literalExpr) Render(v Vars) string {
	return string(e)
}

func (e literalExpr) Name() string {
	return string(e)
}
