package liquid

import "fmt"
import "regexp"

// Tag implements a parsing interface for generating liquid nodes
type Tag func(name, markup string, tokenizer *Tokenizer, ctx *parseContext) Node

// RegisterTag registers a new tag (big surprise)
// XXX: mutex me
func RegisterTag(name string, tag Tag) {
	registeredTags[name] = tag
}

func GetTag(name string) (Tag, bool) {
	if t, ok := registeredTags[name]; ok {
		return t, true
	}
	return nil, false
}

// registeredTags are all known tags
// XXX: do this better
var registeredTags = map[string]Tag{
	"comment": NewCommentTag,
	"assign":  NewAssignTag,
}

// NewCommentTag handles {% comment %} [..] {% endcomment %} blocks
func NewCommentTag(name, markup string, tokenizer *Tokenizer, ctx *parseContext) Node {

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

var assignSyntax = regexp.MustCompile(fmt.Sprintf(`(?ms)(%v)\s*=\s*(.*)\s*`, variableSignatureRegexp.String()))

func NewAssignTag(name, markup string, tokenizer *Tokenizer, ctx *parseContext) Node {

	if submatches := assignSyntax.FindAllStringSubmatch(markup[2:len(markup)-2], -1); len(submatches) > 0 {

		v, err := CreateVariable(submatches[0][2])
		if err != nil {
			panic(err)
		}

		return AssignTag{
			to:   submatches[0][1],
			from: v,
		}
	}

	// localized syntax error
	panic(ErrSyntax("errors.syntax.assign"))
}

type AssignTag struct {
	to   string
	from *Variable
}

func (t AssignTag) Render(v *Vars) (string, error) {
	expr, err := t.from.Render(v)
	if err != nil {
		return "", err
	}
	v.v[t.to] = expr

	// assign tags leave no trace
	return "", nil
}

func (t AssignTag) Blank() bool {
	return true
}
