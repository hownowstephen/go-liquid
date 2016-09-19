package liquid

import "regexp"

// VariableLookup objects capture an expression as a series of lookup objects
// to be applied when rendering the variable
type VariableLookup struct {
	name         string
	lookups      []string
	commandFlags uint
}

func (v *VariableLookup) Render() string {
	// XXX: this is wrong, obviously
	return v.name
}

var (
	squareBracketedRegexp = regexp.MustCompile(`(?sm)\A\[(.*)\]\z`)
	// COMMAND_METHODS = ['size'.freeze, 'first'.freeze, 'last'.freeze]
	commandMethods = []string{"size", "first", "last"}
)

func ParseVariableLookup(markup string) *VariableLookup {

	var commandFlags uint

	lookups := variableParserRegexp.FindAllString(markup, -1)

	if len(lookups) == 0 {
		panic("OHNO WAT DO NOW")
	}

	name := lookups[0]

	if squareBracketedRegexp.MatchString(name) {
		// XXX: don't Render()
		name = ParseExpression(name).Render()
	}

	if len(lookups) == 1 {
		return &VariableLookup{name: name}
	}

	// Drop the name from the list
	lookups = lookups[1:]

	for i, lookup := range lookups {
		if squareBracketedRegexp.MatchString(lookup) {
			lookups[i] = ParseExpression(lookup).Render()
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
		lookups:      lookups,
		commandFlags: commandFlags,
	}
}
