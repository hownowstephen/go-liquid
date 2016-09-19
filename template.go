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

type Template struct {
	nodes []node
}

type Node interface {
	Render(Vars) (string, error)
}

type node struct {
	value    string
	nodelist []node
}

func (n node) Blank() bool {
	return false
}

func (n node) Render(vars Vars) (string, error) {
	return n.value, nil
}

// Tag implements a parsing interface for generating liquid nodes
type Tag interface {
	Parse(name, markup string, tokenizer *Tokenizer, ctx *parseContext) node
}

type commentTag struct{}

func (t *commentTag) Parse(name, markup string, tokenizer *Tokenizer, ctx *parseContext) node {

	subctx := &parseContext{
		line: ctx.line,
		end:  fmt.Sprintf("end%v", name),
	}

	nodelist, err := consume(tokenizer, subctx)
	if err != nil {
		panic(err)
	}

	ctx.line = subctx.line

	// body.parse(tokens, parse_context) do |end_tag_name, end_tag_params|
	//     @blank &&= body.blank?

	//     return false if end_tag_name == block_delimiter
	//     unless end_tag_name
	//       raise SyntaxError.new(parse_context.locale.t("errors.syntax.tag_never_closed".freeze, block_name: block_name))
	//     end

	//     # this tag is not registered with the system
	//     # pass it to the current block for special handling or error reporting
	//     unknown_tag(end_tag_name, end_tag_params, tokens)
	//   end
	return node{"", append([]node{{value: markup}}, nodelist...)}
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

func consume(tokenizer *Tokenizer, ctx *parseContext) ([]node, error) {
	var nodeList []node

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
					return append(nodeList, node{value: markup}), err
				} else if tag := RegisteredTags[tagName]; tag != nil {
					newTag := tag.Parse(tagName, markup, tokenizer, ctx)
					blank = blank && newTag.Blank()
					nodeList = append(nodeList, newTag)
				} else {
					// @TODO: Liquid returns the value instead. why?
					// return yield tag_name, markup
					return nil, LiquidError(fmt.Sprintf("Unknown tag: %v, %v", tagName, markup), ctx)
				}
			} else {
				return nil, LiquidError(fmt.Sprintf("Missing tag terminator: %v", token), ctx)
			}
		case strings.HasPrefix(token, varStartToken):
			nodeList = append(nodeList, createVariable(token, ctx))
			blank = false
		default:
			nodeList = append(nodeList, node{value: token})
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
	// TODO: this strips out the values being split on, but we need those!
	ctx := &parseContext{line: 0}

	nodeList, err := consume(tokenizer, ctx)

	return &Template{nodeList}, err
}

func (t *Template) Render(vars Vars) (string, error) {
	if len(t.nodes) == 0 || t.nodes[0].Blank() {
		return "", nil
	}

	return t.nodes[0].Render(vars)
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

func createVariable(token string, ctx *parseContext) node {
	parsed := contentOfVariableRegexp.FindStringSubmatch(token)
	return node{value: parsed[0]}
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
