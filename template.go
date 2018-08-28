package liquid

import (
	"fmt"
	"regexp"
	"strings"
)

// Combines template.rb, document.rb and block_body.rb

var (
	fullTokenRegexp         = regexp.MustCompile(fmt.Sprintf(`(?m)\A%v\s*(\w+)\s*(.*)?%v\z`, tagStartRegexp, tagEndRegexp)) // om
	contentOfVariableRegexp = regexp.MustCompile(fmt.Sprintf(`(?m)\A%v(.*)%v\z`, variableStartRegexp, variableEndRegexp))   // om
	tokenIsBlankRegexp      = regexp.MustCompile(`\A\s*\z`)
)

const (
	tagStartToken = "{%"
	varStartToken = "{{"
)

// Template is a parsed liquid string containing a list
// of Nodes that can be used to render an output
type Template struct {
	Nodes []Node
}

// Node must be implemented by all parts of a template, and
// provides the necessary rendering handlers to allow generating
// a final output
type Node interface {
	Render(Vars) (string, error)
	Blank() bool
}

type stringNode string

func (n stringNode) Render(v Vars) (string, error) {
	return string(n), nil
}

func (n stringNode) Blank() bool {
	return n == ""
}

// Tag implements a parsing interface for generating liquid Nodes
type Tag interface {
	Parse(name, markup string, tokenizer *Tokenizer, ctx *parseContext) Node
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
		Nodes: nodelist,
	}
}

// RegisterTag registers a new tag (big surprise)
// and probably needs a mutex?
func RegisterTag(name string, tag Tag) {
	RegisteredTags[name] = tag
}

// RegisteredTags are all known tags
var RegisteredTags = map[string]Tag{
	"comment": &commentTag{},
}

type parseContext struct {
	line int
	end  string
}

func (c *parseContext) String() string {
	return fmt.Sprintf("Line: %v, End: %v", c.line, c.end)
}

func tokensToNodeList(tokenizer *Tokenizer, ctx *parseContext) ([]Node, error) {
	var nodeList []Node

	blank := true

	var token string
	var done error

	for done == nil {
		token, done = tokenizer.Next()

		if token == "" {
			continue
		}

		switch {
		case strings.HasPrefix(token, tagStartToken):
			if matched := fullTokenRegexp.FindStringSubmatch(token); len(matched) > 0 {
				markup, tagName := matched[0], matched[1]
				// Check for end tag
				if strings.HasPrefix(tagName, "end") {
					var err error
					if tagName != ctx.end {
						err = LiquidError(fmt.Sprintf("Unexpected end tag: %v, %v", tagName, markup), ctx)
					}
					return nodeList, err
				} else if tag, ok := RegisteredTags[tagName]; ok {
					newTag := tag.Parse(tagName, markup, tokenizer, ctx)
					blank = blank && newTag.Blank()
					nodeList = append(nodeList, newTag)
				} else if tagName == "else" || tagName == "end" {
					return nil, ErrSyntax("Unexpected outer 'else' tag")
				} else {
					return nil, ErrSyntax(fmt.Sprintf("Unknown tag '%v'", tagName))
				}
			} else {
				return nil, LiquidError(fmt.Sprintf("Missing tag terminator: %v", token), ctx)
			}

		case strings.HasPrefix(token, varStartToken):
			nodeList = append(nodeList, createVariable(token, ctx))
			blank = false
		default:
			nodeList = append(nodeList, stringNode(token))
			blank = blank && tokenIsBlankRegexp.MatchString(token)
		}
		ctx.line += strings.Count(token, "\n")
	}

	return nodeList, nil
}

// ParseTemplate performs the parsing step from Liquid::BlockBody.parse
func ParseTemplate(template string) (*Template, error) {

	// tokenize the source
	tokenizer := NewTokenizer(template)
	ctx := &parseContext{line: 0}
	nodeList, err := tokensToNodeList(tokenizer, ctx)

	return &Template{nodeList}, err
}

// Render the template with the supplied variables
func (t *Template) Render(vars Vars) (string, error) {
	if len(t.Nodes) == 0 || t.Nodes[0].Blank() {
		return "", nil
	}

	// Obviously we need to actually render the rest of the Nodes.
	return t.Nodes[0].Render(vars)
}

//     def render_node(node, context)
//       node_output = (node.respond_to?(:render) ? node.render(context) : node)
//       node_output = node_output.is_a?(Array) ? node_output.join : node_output.to_s

//       context.resource_limits.render_length += node_output.length
//       if context.resource_limits.reached?
//         raise MemoryError.new("Memory limits exceeded".freeze)
//       end
//       node_output
//     end

//     def create_variable(token, parse_context)
//       token.scan(ContentOfVariable) do |content|
//         markup = content.first
//         return Variable.new(markup, parse_context)
//       end
//       raise_missing_variable_terminator(token, parse_context)
//     end

func createVariable(token string, ctx *parseContext) Node {
	parsed := contentOfVariableRegexp.FindStringSubmatch(token)

	v, err := CreateVariable(parsed[1])
	if err != nil {
		panic(err)
	}
	return v
}

type blockNode struct {
	tag   string
	Nodes []Node
}

func (n blockNode) Render(v Vars) (string, error) {
	panic("unimplemented")
}

func (n blockNode) Blank() bool {
	return len(n.Nodes) == 0
}

//     def raise_missing_tag_terminator(token, parse_context)
//       raise SyntaxError.new(parse_context.locale.t("errors.syntax.tag_termination".freeze, token: token, tag_end: TagEnd.inspect))
//     end

//     def raise_missing_variable_terminator(token, parse_context)
//       raise SyntaxError.new(parse_context.locale.t("errors.syntax.variable_termination".freeze, token: token, tag_end: VariableEnd.inspect))
//     end

//     def registered_tags
//       Template.tags
//     end
//   end
// end
