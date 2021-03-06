package liquid

import (
	"reflect"
	"testing"
)

func checkTemplate(t *testing.T, tpl string, want []Node) {
	template, err := ParseTemplate(tpl)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(template.Nodes, want) {
		t.Errorf("Template parsed wrong, want: %v, got: %v", want, template.Nodes)
	}
}

func TestBlankSpace(t *testing.T) {
	checkTemplate(t, "  ", []Node{
		stringNode("  "),
	})
}

func TestVariableBeginning(t *testing.T) {
	checkTemplate(t, "{{funk}}  ", []Node{
		testVariableNode("funk"),
		stringNode("  "),
	})
}

func TestVariableEnd(t *testing.T) {
	checkTemplate(t, "  {{funk}}", []Node{
		stringNode("  "),
		testVariableNode("funk"),
	})
}

func TestVariableMiddle(t *testing.T) {
	checkTemplate(t, "  {{funk}}  ", []Node{
		stringNode("  "),
		testVariableNode("funk"),
		stringNode("  "),
	})
}

func TestVariableManyEmbeddedFragments(t *testing.T) {
	checkTemplate(t, "  {{funk}} {{so}} {{brother}} ", []Node{
		stringNode("  "),
		testVariableNode("funk"),
		stringNode(" "),
		testVariableNode("so"),
		stringNode(" "),
		testVariableNode("brother"),
		stringNode(" "),
	})
}

func TestWithBlock(t *testing.T) {
	checkTemplate(t, `  {% comment %} {% endcomment %} `, []Node{
		stringNode("  "),
		BlockNode{
			Tag:    "comment",
			markup: "{% comment %}",
			Nodes:  []Node{stringNode(" ")},
		},
		stringNode(" "),
	})
}

func TestWithCustomTag(t *testing.T) {
	RegisterTag("testtag", &commentTag{})
	checkTemplate(t, `{% testtag %} {% endtesttag %}`, []Node{
		BlockNode{
			Tag:    "testtag",
			markup: "{% testtag %}",
			Nodes:  []Node{stringNode(" ")},
		},
	})
}

func TestVariableLookup(t *testing.T) {
	checkTemplate(t, `{{a.b.first}}<br>{{c.d[0]}}<br>{{e.f | first}}`, []Node{
		testVariableNode("a.b.first"),
		stringNode("<br>"),
		testVariableNode("c.d[0]"),
		stringNode("<br>"),
		testVariableNode("e.f | first"),
	})
}

func TestParseIfBlock(t *testing.T) {
	checkTemplate(t, `{% if x > 100 %}Huge{% elsif x > 10 %}Big{% else %}Normal{% endif %}`, []Node{
		BlockNode{
			Tag:    "if",
			markup: "{% if x > 100 %}",
			Nodes: []Node{
				stringNode("Huge"),
				elseNode{tag: "elsif", markup: "{% elsif x > 10 %}"},
				stringNode("Big"),
				elseNode{tag: "else", markup: "{% else %}"},
				stringNode("Normal"),
			},
		},
	})
}

func testVariableNode(v string) Node {
	variable, err := CreateVariable(v)
	if err != nil {
		return nil
	}
	return variable
}
