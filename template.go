package liquid

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	FullToken         = regexp.MustCompile(fmt.Sprintf(`(?m)\A%v\s*(\w+)\s*(.*)?%v\z`, TagStart, TagEnd)) // om
	ContentOfVariable = regexp.MustCompile(fmt.Sprintf(`(?m)\A%v(.*)%v\z`, VariableStart, VariableEnd))   // om
	TokenIsBlank      = regexp.MustCompile(`\A\s*\z`)
)

const (
	TagStartToken = "{%"
	VarStartToken = "{{"
)

type Template struct {
	nodes []node
}

type node struct {
	value string
}

func (n node) Blank() bool {
	return false
}

type Tag struct{}

func (t *Tag) Parse(name, markup string, ctx *parseContext) node {
	return node{}
}

var RegisteredTags = map[string]*Tag{}

type parseContext struct {
	line int
}

func (c *parseContext) String() string {
	return fmt.Sprintf("Line: %v", c.line)
}

// ParseTemplate performs the parsing step from Liquid::BlockBody.parse
func ParseTemplate(template string) (*Template, error) {

	// tokenize the source
	tokens := TemplateParser.Split(template, -1)
	// TODO: this strips out the values being split on, but we need those!
	ctx := &parseContext{line: 0}

	var nodeList []node

	blank := true

	for _, token := range tokens {
		if token == "" {
			continue
		}

		switch {
		case strings.HasPrefix(token, TagStartToken):
			if matched := FullToken.FindStringSubmatch(token); len(matched) > 0 {
				tagName := matched[0]
				markup := matched[1]
				if tag := RegisteredTags[tagName]; tag != nil {
					newTag := tag.Parse(tagName, markup, ctx)
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
		case strings.HasPrefix(token, VarStartToken):
			nodeList = append(nodeList, createVariable(token, ctx))
			blank = false
		default:
			nodeList = append(nodeList, node{token})
			blank = blank && TokenIsBlank.MatchString(token)
		}
		ctx.line += strings.Count(token, "\n")
	}

	return &Template{nodeList}, nil
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
	parsed := ContentOfVariable.FindStringSubmatch(token)
	fmt.Println("M'VAR", parsed)
	return node{}
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
