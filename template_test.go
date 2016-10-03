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
	if !reflect.DeepEqual(template.nodes, want) {
		t.Errorf("Template parsed wrong, want: %v, got: %v", want, template.nodes)
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

	//     assert_equal Variable, template.root.nodelist[0].class
	//     assert_equal String, template.root.nodelist[1].class
}

func TestVariableEnd(t *testing.T) {
	checkTemplate(t, "  {{funk}}", []Node{
		stringNode("  "),
		testVariableNode("funk"),
	})

	//     assert_equal String, template.root.nodelist[0].class
	//     assert_equal Variable, template.root.nodelist[1].class
}

func TestVariableMiddle(t *testing.T) {
	checkTemplate(t, "  {{funk}}  ", []Node{
		stringNode("  "),
		testVariableNode("funk"),
		stringNode("  "),
	})
	//     assert_equal String, template.root.nodelist[0].class
	//     assert_equal Variable, template.root.nodelist[1].class
	//     assert_equal String, template.root.nodelist[2].class
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
	//     assert_equal [String, Variable, String, Variable, String, Variable, String],
	//       block_types(template.root.nodelist)
}

func TestWithBlock(t *testing.T) {
	checkTemplate(t, `  {% comment %} {% endcomment %} `, []Node{
		stringNode("  "),
		blockNode{
			tag:   "comment",
			nodes: []Node{stringNode(" ")},
		},
		stringNode(" "),
	})
}

func TestWithCustomTag(t *testing.T) {
	RegisterTag("testtag", &commentTag{})
	checkTemplate(t, `{% testtag %} {% endtesttag %}`, []Node{
		blockNode{
			tag:   "testtag",
			nodes: []Node{stringNode(" ")},
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
