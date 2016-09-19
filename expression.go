package liquid

import (
	"fmt"
	"regexp"
	"strconv"
)

var expressionLiterals = map[string]string{}

// LITERALS = {
//   nil => nil, 'nil'.freeze => nil, 'null'.freeze => nil, ''.freeze => nil,
//   'true'.freeze  => true,
//   'false'.freeze => false,
//   'blank'.freeze => MethodLiteral.new(:blank?, '').freeze,
//   'empty'.freeze => MethodLiteral.new(:empty?, '').freeze
// }

type Renderable interface {
	Render() string
}

type NilLiteral struct{}

func (n NilLiteral) Render() string { return "nil" }

type BoolLiteral bool

func (b BoolLiteral) Render() string {
	if b {
		return "true"
	}
	return "false"
}

type EmptyLiteral struct{}

func (e EmptyLiteral) Render() string { return "" }

type StringLiteral struct {
	value string
}

func (s StringLiteral) Render() string {
	return s.value
}

type RangeLookup struct {
	start int
	end   int
}

func (r RangeLookup) Render() string {
	return fmt.Sprintf("%v..%v", r.start, r.end)
}

type IntegerLiteral struct {
	value int
}

func (i IntegerLiteral) Render() string {
	return strconv.Itoa(i.value)
}

type FloatLiteral struct {
	value float64
}

func (i FloatLiteral) Render() string {
	return "this is a cool float"
}

var (
	singleQuotedStringRegex = regexp.MustCompile(`(?ms)\A'(.*)'\z`)
	doubleQuotedStringRegex = regexp.MustCompile(`(?ms)\A"(.*)"\z`)
	integerRegex            = regexp.MustCompile(`\A(-?\d+)\z`)
	rangeRegex              = regexp.MustCompile(`\A\((\S+)\.\.(\S+)\)\z`)
	floatRegex              = regexp.MustCompile(`\A(-?\d[\d\.]+)\z`)
	//     when //m # Single quoted strings
//       $1
//     when /\A"(.*)"\z/m # Double quoted strings
//       $1
//     when /\A(-?\d+)\z/ # Integer and floats
//       $1.to_i
//     when /\A\((\S+)\.\.(\S+)\)\z/ # Ranges
//       RangeLookup.parse($1, $2)
//     when /\A(-?\d[\d\.]+)\z/ # Floats
)

// ParseExpression takes an expression and converts it into something usable by Liquid
// XXX: not sure that renderable is the right approach. maybe just interface{} and then the render
// func will have to know how to deal with it? otherwise we have all these literals
func ParseExpression(markup string) Renderable {
	switch markup {
	case "nil", "null", "":
		return NilLiteral{}
	case "false":
		return BoolLiteral(false)
	case "true":
		return BoolLiteral(true)
	case "blank", "empty":
		return EmptyLiteral{}
	}

	if singleQuotedStringRegex.MatchString(markup) || doubleQuotedStringRegex.MatchString(markup) {
		return StringLiteral{markup[1 : len(markup)-1]}
	}

	if integerRegex.MatchString(markup) {
		value, err := strconv.Atoi(markup)
		if err != nil {
			// XXX: this needs a real handler
			panic(err)
		}
		return IntegerLiteral{value}
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
		return RangeLookup{start, end}
	}

	if floatRegex.MatchString(markup) {
		f, err := strconv.ParseFloat(markup, 64)
		if err != nil {
			panic(err)
		}
		return FloatLiteral{f}
	}

	return ParseVariableLookup(markup)
}

// def self.parse(markup)
//   if LITERALS.key?(markup)
//     LITERALS[markup]
//   else
//     case markup
//     when /\A'(.*)'\z/m # Single quoted strings
//       $1
//     when /\A"(.*)"\z/m # Double quoted strings
//       $1
//     when /\A(-?\d+)\z/ # Integer and floats
//       $1.to_i
//     when /\A\((\S+)\.\.(\S+)\)\z/ # Ranges
//       RangeLookup.parse($1, $2)
//     when // # Floats
//       $1.to_f
//     else
//       VariableLookup.parse(markup)
//     end
//   end
// end
