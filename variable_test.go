package liquid

import (
	"reflect"
	"testing"
)

// Unit tests (unit/variable_unit_test.rb)

func checkVariable(t *testing.T, markup, wantName string, wantFilters []filter) {
	v, err := CreateVariable(markup)
	if err != nil {
		t.Error(err)
		return
	}

	vl := ParseVariableLookup(markup)
	if vl == nil {
		t.Errorf("Couldn't create VariableLookup for %v", markup)
		return
	}

	if v.name != vl.name {
		t.Errorf("%v lookup name mismatch, want: %v, got %v", markup, vl.name, v.name)
	}

	if v.name != wantName {
		t.Errorf("%v name mismatched, want: %v, got: %v", markup, wantName, v.name)
	}
	if !reflect.DeepEqual(v.filters, wantFilters) {
		t.Errorf("%v filters mismatched, want: %v, got: %v", markup, wantFilters, v.filters)
	}
}

func TestVariable(t *testing.T) {
	checkVariable(t, "hello", "hello", nil)
}

func TestFilters(t *testing.T) {
	checkVariable(t, "hello | textileze", "hello", []filter{
		{filter: "textileze"},
	})
	checkVariable(t, "hello | textileze | paragraph", "hello", []filter{
		{filter: "textileze"},
		{filter: "paragraph"},
	})
	checkVariable(t, `hello | strftime: '%Y'`, "hello", []filter{
		{filter: "strftime", args: []string{`%Y`}},
	})
	checkVariable(t, `'typo' | link_to: 'Typo', true`, "typo", []filter{
		{filter: "link_to", args: []string{`Typo`, "true"}},
	})
	checkVariable(t, `'typo' | link_to: 'Typo', false`, "typo", []filter{
		{filter: "link_to", args: []string{`Typo`, "false"}},
	})
	checkVariable(t, `'foo' | repeat: 3`, "foo", []filter{
		{filter: "repeat", args: []string{"3"}},
	})
	checkVariable(t, `'foo' | repeat: 3, 3`, "foo", []filter{
		{filter: "repeat", args: []string{"3", "3"}},
	})
	checkVariable(t, `'foo' | repeat: 3, 3, 3`, "foo", []filter{
		{filter: "repeat", args: []string{"3", "3", "3"}},
	})
	checkVariable(t, `hello | strftime: '%Y, okay?'`, "hello", []filter{
		{filter: "strftime", args: []string{`%Y, okay?`}},
	})
	checkVariable(t, `hello | things: "%Y, okay?", 'the other one'`, "hello", []filter{
		{filter: "things", args: []string{`%Y, okay?`, "the other one"}},
	})
}

func TestFilterWithDateParameter(t *testing.T) {
	checkVariable(t, `'2006-06-06' | date: "%m/%d/%y"`, "2006-06-06", []filter{
		{filter: "date", args: []string{`%m/%d/%y`}},
	})
}

func TestFiltersWithoutWhitespace(t *testing.T) {
	checkVariable(t, "hello | textileze | paragraph", "hello", []filter{
		{filter: "textileze"},
		{filter: "paragraph"},
	})
	checkVariable(t, "hello|textileze|paragraph", "hello", []filter{
		{filter: "textileze"},
		{filter: "paragraph"},
	})
	checkVariable(t, "hello|replace:'foo','bar'|textileze", "hello", []filter{
		{filter: "replace", args: []string{"foo", "bar"}},
		{filter: "textileze"},
	})
}

// XXX: This requires a lax parser, which we don't have
// func TestSymbol(t *testing.T) {
// 	checkVariable(t, "http://disney.com/logo.gif | image: 'med' ", "http://disney.com/logo.gif", []filter{
// 		{filter: "image", args: []string{"med"}},
// 	})
// }

func TestStringToFilter(t *testing.T) {
	checkVariable(t, "'http://disney.com/logo.gif' | image: 'med' ", "http://disney.com/logo.gif", []filter{
		{filter: "image", args: []string{"med"}},
	})
}

// merges test_string_single_quoted and test_string_double_quoted
func TestStringQuoting(t *testing.T) {
	checkVariable(t, `'hello'`, "hello", nil)
	checkVariable(t, `"hello"`, "hello", nil)
}

//   def test_integer
//     var = create_variable(%( 1000 ))
//     assert_equal 1000, var.name
//   end

//   def test_float
//     var = create_variable(%( 1000.01 ))
//     assert_equal 1000.01, var.name
//   end

//   def test_dashes
//     assert_equal VariableLookup.new('foo-bar'), create_variable('foo-bar').name
//     assert_equal VariableLookup.new('foo-bar-2'), create_variable('foo-bar-2').name

//     with_error_mode :strict do
//       assert_raises(Liquid::SyntaxError) { create_variable('foo - bar') }
//       assert_raises(Liquid::SyntaxError) { create_variable('-foo') }
//       assert_raises(Liquid::SyntaxError) { create_variable('2foo') }
//     end
//   end

//   def test_string_with_special_chars
//     var = create_variable(%( 'hello! $!@.;"ddasd" ' ))
//     assert_equal 'hello! $!@.;"ddasd" ', var.name
//   end

//   def test_string_dot
//     var = create_variable(%( test.test ))
//     assert_equal VariableLookup.new('test.test'), var.name
//   end

//   def test_filter_with_keyword_arguments
//     var = create_variable(%( hello | things: greeting: "world", farewell: 'goodbye'))
//     assert_equal VariableLookup.new('hello'), var.name
//     assert_equal [['things', [], { 'greeting' => 'world', 'farewell' => 'goodbye' }]], var.filters
//   end

//   def test_lax_filter_argument_parsing
//     var = create_variable(%( number_of_comments | pluralize: 'comment': 'comments' ), error_mode: :lax)
//     assert_equal VariableLookup.new('number_of_comments'), var.name
//     assert_equal [['pluralize', ['comment', 'comments']]], var.filters
//   end

//   def test_strict_filter_argument_parsing
//     with_error_mode(:strict) do
//       assert_raises(SyntaxError) do
//         create_variable(%( number_of_comments | pluralize: 'comment': 'comments' ))
//       end
//     end
//   end

//   def test_output_raw_source_of_variable
//     var = create_variable(%( name_of_variable | upcase ))
//     assert_equal " name_of_variable | upcase ", var.raw
//   end

//   def test_variable_lookup_interface
//     lookup = VariableLookup.new('a.b.c')
//     assert_equal 'a', lookup.name
//     assert_equal ['b', 'c'], lookup.lookups
//   end

//   private

//   def create_variable(markup, options = {})
//     Variable.new(markup, ParseContext.new(options))
//   end
// end

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

func TestSimpleVariable(t *testing.T) {
	checkTemplateRender(t, "{{test}}", map[string]interface{}{"test": "worked"}, "worked")
}

// def test_simple_variable
//     template = Template.parse(%({{test}}))
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
//     template = Template.parse(%({{ test }}))
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
//     template = Template.parse(%({{ test.test }}))
//     assert_equal 'worked', template.render!('test' => { 'test' => 'worked' })
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
//     template = Template.parse(%({{ test }}))
//     template.assigns['test'] = 'worked'
//     assert_equal 'worked', template.render!
//   end

//   def test_reuse_parsed_template
//     template = Template.parse(%({{ greeting }} {{ name }}))
//     template.assigns['greeting'] = 'Goodbye'
//     assert_equal 'Hello Tobi', template.render!('greeting' => 'Hello', 'name' => 'Tobi')
//     assert_equal 'Hello ', template.render!('greeting' => 'Hello', 'unknown' => 'Tobi')
//     assert_equal 'Hello Brian', template.render!('greeting' => 'Hello', 'name' => 'Brian')
//     assert_equal 'Goodbye Brian', template.render!('name' => 'Brian')
//     assert_equal({ 'greeting' => 'Goodbye' }, template.assigns)
//   end

//   def test_assigns_not_polluted_from_template
//     template = Template.parse(%({{ test }}{% assign test = 'bar' %}{{ test }}))
//     template.assigns['test'] = 'baz'
//     assert_equal 'bazbar', template.render!
//     assert_equal 'bazbar', template.render!
//     assert_equal 'foobar', template.render!('test' => 'foo')
//     assert_equal 'bazbar', template.render!
//   end

//   def test_hash_with_default_proc
//     template = Template.parse(%(Hello {{ test }}))
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
