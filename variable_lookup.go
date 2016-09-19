package liquid

import "regexp"

// VariableLookup objects capture an expression as a series of lookup objects
// to be applied when rendering the variable
type VariableLookup struct {
	name         Expression
	lookups      []Expression
	commandFlags uint
}

func (v *VariableLookup) Render(vars Vars) string {
	// XXX: this is wrong, obviously
	return v.name.Render(vars)
}

func (v *VariableLookup) Name() string {
	return v.name.Render(nil)
}

var (
	squareBracketedRegexp = regexp.MustCompile(`(?sm)\A\[(.*)\]\z`)
	// COMMAND_METHODS = ['size'.freeze, 'first'.freeze, 'last'.freeze]
	commandMethods = []string{"size", "first", "last"}
)

func ParseVariableLookup(markup string) *VariableLookup {

	var name Expression
	var commandFlags uint

	lookups := variableParserRegexp.FindAllString(markup, -1)

	if len(lookups) == 0 {
		panic("OHNO WAT DO NOW")
	}

	name = literalExpr(lookups[0])

	if squareBracketedRegexp.MatchString(lookups[0]) {
		// XXX: don't Render()
		name = ParseExpression(lookups[0])
	}

	if len(lookups) == 1 {
		return &VariableLookup{name: name}
	}

	// Drop the name from the list
	lookupExpressions := make([]Expression, len(lookups)-1)

	for i, lookup := range lookups[1:] {
		if squareBracketedRegexp.MatchString(lookup) {
			lookupExpressions[i] = ParseExpression(lookup)
		} else {
			lookupExpressions[i] = literalExpr(lookup)
		}

		// can be optimized
		for _, command := range commandMethods {
			if lookup == command {
				commandFlags |= 1 << uint(i)
			}
		}
	}

	return &VariableLookup{
		name:         name,
		lookups:      lookupExpressions,
		commandFlags: commandFlags,
	}
}
