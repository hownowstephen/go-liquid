package liquid

import (
	"reflect"
	"testing"
)

func checkTemplate(t *testing.T, tpl string, want []node) {
	template, err := ParseTemplate(tpl)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(template.nodes, want) {
		t.Errorf("Template parsed wrong, want: %v, got: %v", want, template.nodes)
	}
}

func TestBlankSpace(t *testing.T) {
	checkTemplate(t, "  ", []node{
		{value: "  "},
	})
}

func TestVariableBeginning(t *testing.T) {
	checkTemplate(t, "{{funk}}  ", []node{
		{value: "{{funk}}"},
		{value: "  "},
	})

	//     assert_equal Variable, template.root.nodelist[0].class
	//     assert_equal String, template.root.nodelist[1].class
}

func TestVariableEnd(t *testing.T) {
	checkTemplate(t, "  {{funk}}", []node{
		{value: "  "},
		{value: "{{funk}}"},
	})

	//     assert_equal String, template.root.nodelist[0].class
	//     assert_equal Variable, template.root.nodelist[1].class
}

func TestVariableMiddle(t *testing.T) {
	checkTemplate(t, "  {{funk}}  ", []node{
		{value: "  "},
		{value: "{{funk}}"},
		{value: "  "},
	})
	//     assert_equal String, template.root.nodelist[0].class
	//     assert_equal Variable, template.root.nodelist[1].class
	//     assert_equal String, template.root.nodelist[2].class
}

func TestVariableManyEmbeddedFragments(t *testing.T) {
	checkTemplate(t, "  {{funk}} {{so}} {{brother}} ", []node{
		{value: "  "},
		{value: "{{funk}}"},
		{value: " "},
		{value: "{{so}}"},
		{value: " "},
		{value: "{{brother}}"},
		{value: " "},
	})
	//     assert_equal [String, Variable, String, Variable, String, Variable, String],
	//       block_types(template.root.nodelist)
}

func TestWithBlock(t *testing.T) {
	checkTemplate(t, `  {% comment %} {% endcomment %} `, []node{
		{value: "  "},
		{nodelist: []node{
			{value: "{% comment %}"},
			{value: " "},
			{value: "{% endcomment %}"},
		}},
		{value: " "},
	})
}

func TestWithCustomTag(t *testing.T) {
	RegisterTag("testtag", &commentTag{})
	checkTemplate(t, `{% testtag %} {% endtesttag %}`, []node{
		node{
			nodelist: []node{
				node{value: "{% testtag %}"},
				node{value: " "},
				node{value: "{% endtesttag %}"},
			},
		},
	})
}
