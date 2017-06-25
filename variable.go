package liquid

import (
	"fmt"
	"regexp"
)

var (
	variableQuotedFragmentRegexp = regexp.MustCompile(fmt.Sprintf("%s(.*)", quotedFragmentRegexp.String())) // om
)

// Vars get passed in the render step of a request
// and contain values and context information for the input
// variables of the request
type Vars map[string]interface{}

// Variable is a single liquid variable expression
// and any associated filters
type Variable struct {
	name    Expression
	filters []Filter
	markup  string
}

func (v *Variable) Render(vars Vars) (string, error) {
	return "VARIABLE_RENDER_UNIMPLEMENTED", nil
}

func (v *Variable) Blank() bool {
	panic("unimplemented")
}

// Filter is used to modify a Variable using the
// liquid pipe syntax "x | f1 | f2"
type Filter struct {
	name   string
	args   []Expression
	kwargs map[string]Expression
}

// CreateVariable performs a parse of the supplied markup
// and returns a Variable object
func CreateVariable(value string) (*Variable, error) {

	matches := variableQuotedFragmentRegexp.FindStringSubmatch(value)

	if len(matches) != 2 {
		return nil, fmt.Errorf("Bad match")
	}

	return ParseStrict(value)
}

// VariableParser is the signature of function required to perform
// a parse of a liquid variable expression
type VariableParser func(markup string) (*Variable, error)

// ParseStrict performs the strictest form of VariableParser parse
func ParseStrict(markup string) (*Variable, error) {
	var filters []Filter
	p, err := NewParser(markup)
	if err != nil {
		return nil, err
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	name := ParseExpression(expr)

	for p.tryConsume(tPipe) {

		filterName, err := p.consume(tIdentifier)
		if err != nil {
			return nil, err
		}

		var filterArgs []string
		// check if there are arguments
		if p.tryConsume(tColon) {
			// parse de args
			for {
				arg, err := p.argument()
				if err != nil {
					return nil, err
				}

				filterArgs = append(filterArgs, arg)

				if !p.tryConsume(tComma) {
					break
				}
			}
		}

		filters = append(filters, ParseFilterExpressions(filterName, filterArgs))
	}

	// consume an end of string, if there
	if _, err = p.consume(tEndOfString); err != nil {
		return nil, err
	}

	return &Variable{
		name:    name,
		filters: filters,
		markup:  markup,
	}, nil
}

// ParseFilterExpressions parses the filter args passed with a liquid variable
func ParseFilterExpressions(name string, unparsedArgs []string) Filter {
	var args []Expression
	kwargs := make(map[string]Expression)

	for _, a := range unparsedArgs {
		// Check for keyword arguments first, anything leftover will be treated like a regular argument
		if submatches := tagAttributesRegexp.FindStringSubmatch(a); len(submatches) > 0 {
			kwargs[submatches[1]] = ParseExpression(submatches[2])
		} else {
			args = append(args, ParseExpression(a))
		}
	}

	// null out the kwargs map if it is empty
	// this ensures that Filter{name: blah} will match a parsed filter
	if len(kwargs) == 0 {
		kwargs = nil
	}

	return Filter{
		name:   name,
		args:   args,
		kwargs: kwargs,
	}
}
