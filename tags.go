package liquid

import "fmt"
import "regexp"

// Tag implements a parsing interface for generating liquid nodes
type Tag interface {
	Parse(name, markup string, tokenizer *Tokenizer, ctx *parseContext) Node
}

// RegisterTag registers a new tag (big surprise)
// XXX: mutex me
func RegisterTag(name string, tag Tag) {
	RegisteredTags[name] = tag
}

// RegisteredTags are all known tags
// XXX: do this better
var RegisteredTags = map[string]Tag{
	"comment": &commentTag{},
	"assign": &assignTag{
		syntax: regexp.MustCompile(fmt.Sprintf(`(?ms)(%v)\s*=\s*(.*)\s*`, variableSignatureRegexp.String())),
	},
}

// An example tag
type commentTag struct{}

func (t *commentTag) Parse(name, markup string, tokenizer *Tokenizer, ctx *parseContext) Node {

	subctx := &parseContext{
		line: ctx.line,
		end:  fmt.Sprintf("end%v", name),
	}

	nodelist, err := tokensToNodeList(tokenizer, subctx)
	if err != nil {
		panic(err)
	}

	ctx.line = subctx.line

	return blockNode{
		tag:   name,
		nodes: nodelist,
	}
}

type assignTag struct {
	syntax *regexp.Regexp
	to     string
	from   *VariableLookup
}

func (t *assignTag) Parse(name, markup string, tokenizer *Tokenizer, ctx *parseContext) Node {
	fmt.Println(name, markup)

	if submatches := t.syntax.FindAllStringSubmatch(markup, -1); len(submatches) > 0 {
		fmt.Println(submatches, len(submatches))

		t.to = submatches[0][1]
		t.from = ParseVariableLookup(submatches[0][2])

		// XXX: Tag shouldn't be an interface like this requiring a struct
		// it should just be a func that returns a Node
		return t
	}

	// localized syntax error
	panic(ErrSyntax("errors.syntax.assign"))
}

func (t *assignTag) Render(v Vars) (string, error) {
	fmt.Println("DOIN A RENDER")
	return "POTATO", nil
}

func (t *assignTag) Blank() bool {
	return false
}
