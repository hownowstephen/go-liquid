package liquid

import (
	"fmt"
	"regexp"
)

var (
	variableQuotedFragmentRegexp = regexp.MustCompile(fmt.Sprintf("%s(.*)", quotedFragmentRegexp.String())) // om
)

type Vars map[string]interface{}

type filter struct {
	filter string
	args   []string
	kwargs map[string]string
}

type Variable struct {
	name    string
	filters []filter
}

func CreateVariable(value string) (*Variable, error) {

	matches := variableQuotedFragmentRegexp.FindStringSubmatch(value)

	if len(matches) != 2 {
		return nil, fmt.Errorf("Bad match")
	}

	return StrictParser{}.Parse(value)
}

type VariableParser interface {
	Parse(markup string) (*Variable, error)
}

type StrictParser struct{}

func (l StrictParser) Parse(markup string) (*Variable, error) {
	var filters []filter
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
	p.consume(tEndOfString)

	return &Variable{
		// FIXME: don't Render()
		name:    name.Render(),
		filters: filters,
	}, nil
}

func ParseFilterExpressions(name string, unparsedArgs []string) filter {
	var args []string
	kwargs := make(map[string]string)

	for _, a := range unparsedArgs {
		//  if matches = a.match(/\A#{TagAttributes}\z/o)
		if submatches := tagAttributesRegexp.FindAllStringSubmatch(a, -1); len(submatches) > 0 {
			fmt.Println("SUBMATCH", submatches)

		} else {
			// FIXME: don't Render() this
			args = append(args, ParseExpression(a).Render())
		}
	}

	// null out the map if it is empty
	if len(kwargs) == 0 {
		kwargs = nil
	}

	return filter{
		filter: name,
		args:   args,
		kwargs: kwargs,
	}

}
