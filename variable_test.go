package liquid

import (
	"reflect"
	"testing"
)

// Unit tests (unit/variable_unit_test.rb)

func checkVariable(t *testing.T, markup, wantName string, wantFilters []Filter, checkLookup bool) *Variable {
	v, err := CreateVariable(markup)
	if err != nil {
		t.Error(err)
		return v
	}

	if checkLookup {
		vl := ParseVariableLookup(markup)
		if vl == nil {
			t.Errorf("Couldn't create VariableLookup for %v", markup)
			return v
		}

		if v.name.Name() != vl.Name() {
			t.Errorf("%v lookup name mismatch, want: %v, got %v", markup, vl.name, v.name.Name())
		}
	}

	if v.name.Name() != wantName {
		t.Errorf("%v name mismatched, want: %v, got: %v", markup, wantName, v.name.Name())
	}

	if !reflect.DeepEqual(v.filters, wantFilters) {
		t.Errorf("%v filters mismatched, want: %v, got: %v", markup, wantFilters, v.filters)
	}

	return v
}

func TestVariable(t *testing.T) {
	checkVariable(t, "hello", "hello", nil, true)
}

func TestFilters(t *testing.T) {
	checkVariable(t, "hello | textileze", "hello", []Filter{
		{name: "textileze"},
	}, true)
	checkVariable(t, "hello | textileze | paragraph", "hello", []Filter{
		{name: "textileze"},
		{name: "paragraph"},
	}, true)
	checkVariable(t, `hello | strftime: '%Y'`, "hello", []Filter{
		{name: "strftime", args: []Expression{stringExpr(`%Y`)}},
	}, true)
	checkVariable(t, `'typo' | link_to: 'Typo', true`, "typo", []Filter{
		{name: "link_to", args: []Expression{stringExpr(`Typo`), boolExpr(true)}},
	}, true)
	checkVariable(t, `'typo' | link_to: 'Typo', false`, "typo", []Filter{
		{name: "link_to", args: []Expression{stringExpr(`Typo`), boolExpr(false)}},
	}, true)
	checkVariable(t, `'foo' | repeat: 3`, "foo", []Filter{
		{name: "repeat", args: []Expression{integerExpr(3)}},
	}, true)
	checkVariable(t, `'foo' | repeat: 3, 3`, "foo", []Filter{
		{name: "repeat", args: []Expression{integerExpr(3), integerExpr(3)}},
	}, true)
	checkVariable(t, `'foo' | repeat: 3, 3, 3`, "foo", []Filter{
		{name: "repeat", args: []Expression{integerExpr(3), integerExpr(3), integerExpr(3)}},
	}, true)
	checkVariable(t, `hello | strftime: '%Y, okay?'`, "hello", []Filter{
		{name: "strftime", args: []Expression{stringExpr(`%Y, okay?`)}},
	}, true)
	checkVariable(t, `hello | things: "%Y, okay?", 'the other one'`, "hello", []Filter{
		{name: "things", args: []Expression{stringExpr(`%Y, okay?`), stringExpr("the other one")}},
	}, true)
}

func TestFilterWithDateParameter(t *testing.T) {
	checkVariable(t, `'2006-06-06' | date: "%m/%d/%y"`, "2006-06-06", []Filter{
		{name: "date", args: []Expression{stringExpr(`%m/%d/%y`)}},
	}, true)
}

func TestFiltersWithoutWhitespace(t *testing.T) {
	checkVariable(t, "hello | textileze | paragraph", "hello", []Filter{
		{name: "textileze"},
		{name: "paragraph"},
	}, true)
	checkVariable(t, "hello|textileze|paragraph", "hello", []Filter{
		{name: "textileze"},
		{name: "paragraph"},
	}, true)
	checkVariable(t, "hello|replace:'foo','bar'|textileze", "hello", []Filter{
		{name: "replace", args: []Expression{stringExpr("foo"), stringExpr("bar")}},
		{name: "textileze"},
	}, true)
}

// XXX: This requires a lax parser, which we don't have
// func TestSymbol(t *testing.T) {
// 	checkVariable(t, "http://disney.com/logo.gif | image: 'med' ", "http://disney.com/logo.gif", []Filter{
// 		{name: "image", args: []Expression{"med"}},
// 	}, true)
// }

func TestStringToFilter(t *testing.T) {
	checkVariable(t, "'http://disney.com/logo.gif' | image: 'med' ", "http://disney.com/logo.gif", []Filter{
		{name: "image", args: []Expression{stringExpr("med")}},
	}, false)
}

// merges test_string_single_quoted and test_string_double_quoted
func TestStringQuoting(t *testing.T) {
	checkVariable(t, `'hello'`, "hello", nil, true)
	checkVariable(t, `"hello"`, "hello", nil, true)
}

func TestIntegerVariable(t *testing.T) {
	v := checkVariable(t, `1000`, `1000`, nil, true)
	if v != nil && v.name != integerExpr(1000) {
		t.Errorf("Expected integer, got %v", reflect.TypeOf(v.name))
	}
}

func TestFloatVariable(t *testing.T) {
	v := checkVariable(t, `1000.01`, `1000.01`, nil, false)
	if v != nil && v.name != floatExpr(1000.01) {
		t.Errorf("Expected float, got %v", reflect.TypeOf(v.name))
	}
}

func TestDashes(t *testing.T) {
	for _, expr := range []string{"foo-bar", "foo-bar-2"} {
		vl := ParseVariableLookup(expr)
		v, err := CreateVariable(expr)
		if err != nil {
			t.Errorf("%v couldn't create variable: %v", expr, err)
			continue
		}
		if v.name.Name() != vl.Name() {
			t.Errorf(`mismatch! Lookup: "%v", Variable: "%v"`, vl.Name(), v.name.Name())
		}
	}

	for _, badExpr := range []string{"foo - bar", "-foo", "2foo"} {
		v, err := CreateVariable(badExpr)
		if err == nil {
			t.Errorf(`expression "%v" should not be a valid variable, got: %v`, badExpr, v.name)
		}
	}
}

func TestStringWithSpecialChars(t *testing.T) {
	v := checkVariable(t, `'hello! $!@.;"ddasd" '`, `hello! $!@.;"ddasd" `, nil, false)
	if v != nil && v.name != stringExpr(`hello! $!@.;"ddasd" `) {
		t.Errorf("wrong type, want stringExpr, got: %v", reflect.TypeOf(v.name))
	}
}

func TestStringDot(t *testing.T) {
	checkVariable(t, "test.test", "test", nil, true)
}

func TestFilterWithKeywordArguments(t *testing.T) {
	checkVariable(t, `hello | things: greeting: "world", farewell: 'goodbye'`, "hello", []Filter{
		Filter{
			name: "things",
			kwargs: map[string]Expression{
				"greeting": stringExpr("world"),
				"farewell": stringExpr("goodbye"),
			},
		},
	}, true)
}

// XXX: lax parsing is not implemented
//   def test_lax_filter_argument_parsing
//     var = create_variable(%( number_of_comments | pluralize: 'comment': 'comments' ), error_mode: :lax)
//     assert_equal VariableLookup.new('number_of_comments'), var.name
//     assert_equal [['pluralize', ['comment', 'comments']]], var.filters
//   end

func TestStringFilterArgumentParsing(t *testing.T) {
	_, err := CreateVariable("number_of_comments | pluralize: 'comment': 'comments'")
	if err == nil {
		t.Error("CreateVariable should have failed due to invalid filters")
	}
}

func TestOutputRawSourceOfVariable(t *testing.T) {
	source := " name_of_variable | upcase "
	v, err := CreateVariable(source)
	if err != nil {
		t.Error(err)
		return
	}
	if v.markup != source {
		t.Errorf("variable markup mismatched - want: %v, got: %v", source, v.markup)
	}
}

func TestVariableLookupInterface(t *testing.T) {
	lookup := ParseVariableLookup("a.b.c")
	if lookup.name != literalExpr("a") {
		t.Errorf(`bad name, want: literalExpr("a"), got: %v("%v")`, reflect.TypeOf(lookup.name), lookup.name.Name())
	}

	if len(lookup.lookups) != 2 {
		t.Errorf("expected 2 lookups, got %v", len(lookup.lookups))
		return
	}

	if !reflect.DeepEqual(lookup.lookups, []Expression{literalExpr("b"), literalExpr("c")}) {
		t.Errorf(`bad lookups, want: [literalExpr("a"), literalExpr("b")], got: [%v("%v"), %v("%v")]`, reflect.TypeOf(lookup.lookups[0]), lookup.lookups[0].Name(), reflect.TypeOf(lookup.lookups[1]), lookup.lookups[1].Name())
	}
}

func checkTemplateRender(t *testing.T, template string, vars map[string]interface{}, want string) {
	tpl, err := ParseTemplate(template)
	if err != nil {
		t.Errorf("Couldn't parse the template: %v", err)
		return
	}

	if got, err := tpl.Render(vars); err == nil && got != want {
		t.Errorf(`Template didn't render properly, want: "%v" got: "%v"`, want, got)
	} else if err != nil {
		t.Errorf(`Template rendering error: %v`, err)
	}
}

// Integration Tests

// def test_simple_variable
//     template = Template.parse(%({{test}}, true))
//     assert_equal 'worked', template.render!('test' => 'worked')
//     assert_equal 'worked wonderfully', template.render!('test' => 'worked wonderfully')
//   end

//   def test_variable_render_calls_to_liquid
//     assert_template_result 'foobar', '{{ foo }}', 'foo' => ThingWithToLiquid.new
//   end

//   def test_simple_with_whitespaces
//     template = Template.parse(%(  {{ test }}  ))
//     assert_equal '  worked  ', template.render!('test' => 'worked')
//     assert_equal '  worked wonderfully  ', template.render!('test' => 'worked wonderfully')
//   end

//   def test_ignore_unknown
//     template = Template.parse(%({{ test }}, true))
//     assert_equal '', template.render!
//   end

//   def test_using_blank_as_variable_name
//     template = Template.parse("{% assign foo = blank %}{{ foo }}")
//     assert_equal '', template.render!
//   end

//   def test_using_empty_as_variable_name
//     template = Template.parse("{% assign foo = empty %}{{ foo }}")
//     assert_equal '', template.render!
//   end

//   def test_hash_scoping
//     template = Template.parse(%({{ test.test }}, true))
//     assert_equal 'worked', template.render!('test' => { 'test' => 'worked' }, true)
//   end

//   def test_false_renders_as_false
//     assert_equal 'false', Template.parse("{{ foo }}").render!('foo' => false)
//     assert_equal 'false', Template.parse("{{ false }}").render!
//   end

//   def test_nil_renders_as_empty_string
//     assert_equal '', Template.parse("{{ nil }}").render!
//     assert_equal 'cat', Template.parse("{{ nil | append: 'cat' }}").render!
//   end

//   def test_preset_assigns
//     template = Template.parse(%({{ test }}, true))
//     template.assigns['test'] = 'worked'
//     assert_equal 'worked', template.render!
//   end

//   def test_reuse_parsed_template
//     template = Template.parse(%({{ greeting }} {{ name }}, true))
//     template.assigns['greeting'] = 'Goodbye'
//     assert_equal 'Hello Tobi', template.render!('greeting' => 'Hello', 'name' => 'Tobi')
//     assert_equal 'Hello ', template.render!('greeting' => 'Hello', 'unknown' => 'Tobi')
//     assert_equal 'Hello Brian', template.render!('greeting' => 'Hello', 'name' => 'Brian')
//     assert_equal 'Goodbye Brian', template.render!('name' => 'Brian')
//     assert_equal({ 'greeting' => 'Goodbye' }, template.assigns)
//   end

//   def test_assigns_not_polluted_from_template
//     template = Template.parse(%({{ test }}{% assign test = 'bar' %}{{ test }}, true))
//     template.assigns['test'] = 'baz'
//     assert_equal 'bazbar', template.render!
//     assert_equal 'bazbar', template.render!
//     assert_equal 'foobar', template.render!('test' => 'foo')
//     assert_equal 'bazbar', template.render!
//   end

//   def test_hash_with_default_proc
//     template = Template.parse(%(Hello {{ test }}, true))
//     assigns = Hash.new { |h, k| raise "Unknown variable '#{k}'" }
//     assigns['test'] = 'Tobi'
//     assert_equal 'Hello Tobi', template.render!(assigns)
//     assigns.delete('test')
//     e = assert_raises(RuntimeError) do
//       template.render!(assigns)
//     end
//     assert_equal "Unknown variable 'test'", e.message
//   end

//   def test_multiline_variable
//     assert_equal 'worked', Template.parse("{{\ntest\n}}").render!('test' => 'worked')
//   end
