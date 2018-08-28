package liquid

import (
	"fmt"
	"regexp"
)

var (
	squareBracketedRegexp = regexp.MustCompile(`(?sm)\A\[(.*)\]\z`)
	commandMethods        = []string{"size", "first", "last"}
)

// VariableLookup objects capture an expression as a series of lookup objects
// to be applied when rendering the variable
type VariableLookup struct {
	name         Expression
	lookups      []Expression
	commandFlags uint
}

func (v *VariableLookup) Evaluate(c Context) Expression {

	name := v.name.Evaluate(c)
	object, err := c.FindVariable(name)
	if err != nil {
		return nil
	}

	return object
	//         def evaluate(context)
	//       name = context.evaluate(@name)
	//       object = context.find_variable(name)

	//       @lookups.each_index do |i|
	//         key = context.evaluate(@lookups[i])

	//         # If object is a hash- or array-like object we look for the
	//         # presence of the key and if its available we return it
	//         if object.respond_to?(:[]) &&
	//             ((object.respond_to?(:key?) && object.key?(key)) ||
	//              (object.respond_to?(:fetch) && key.is_a?(Integer)))

	//           # if its a proc we will replace the entry with the proc
	//           res = context.lookup_and_evaluate(object, key)
	//           object = res.to_liquid

	//           # Some special cases. If the part wasn't in square brackets and
	//           # no key with the same name was found we interpret following calls
	//           # as commands and call them on the current object
	//         elsif @command_flags & (1 << i) != 0 && object.respond_to?(key)
	//           object = object.send(key).to_liquid

	//           # No key was present with the desired value and it wasn't one of the directly supported
	//           # keywords either. The only thing we got left is to return nil or
	//           # raise an exception if `strict_variables` option is set to true
	//         else
	//           return nil unless context.strict_variables
	//           raise Liquid::UndefinedVariable, "undefined variable #{key}"
	//         end

	//         # If we are dealing with a drop here we have to
	//         object.context = context if object.respond_to?(:context=)
	//       end

	//       object
	// end
}

func (v *VariableLookup) Name() string {
	return v.name.Name()
}

func ParseVariableLookup(markup string) *VariableLookup {

	var name Expression
	var commandFlags uint

	lookups := variableParserRegexp.FindAllString(markup, -1)

	if len(lookups) == 0 {
		fmt.Println("No lookups found, returning nil")
		return nil
	}

	name = literalExpr(lookups[0])

	if m := squareBracketedRegexp.FindStringSubmatch(lookups[0]); len(m) == 2 {
		name = ParseExpression(m[1])
	}

	if len(lookups) == 1 {
		return &VariableLookup{name: name}
	}

	// Drop the name from the list
	lookupExpressions := make([]Expression, len(lookups)-1)

	for i, lookup := range lookups[1:] {
		if m := squareBracketedRegexp.FindStringSubmatch(lookups[0]); len(m) == 2 {
			name = ParseExpression(m[1])
		} else {
			lookupExpressions[i] = literalExpr(lookup)
		}

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
