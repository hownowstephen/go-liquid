package liquid

import (
	"reflect"
	"regexp"
	"testing"
)

func testRegex(t *testing.T, r *regexp.Regexp, raw string, want []string) {
	got := r.FindAllString(raw, -1)
	if len(got) == 0 {
		got = make([]string, 0)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Regex %v failed, want: %v, got: %v", r.String(), want, got)
	}
}

func TestEmpty(t *testing.T) {
	testRegex(t, QuotedFragment, ``, []string{})
}

func TestQuote(t *testing.T) {
	testRegex(t, QuotedFragment, `"arg 1"`, []string{`"arg 1"`})
}

func TestWords(t *testing.T) {
	testRegex(t, QuotedFragment, `arg1, arg2`, []string{`arg1`, `arg2`})
}

func TestTags(t *testing.T) {
	testRegex(t, QuotedFragment, `<tr> </tr>`, []string{`<tr>`, `</tr>`})
	testRegex(t, QuotedFragment, `<tr></tr>`, []string{`<tr></tr>`})
	testRegex(t, QuotedFragment, `<style class="hello">' </style>`, []string{`<style`, `class="hello">`, `</style>`})

}

func TestDoubleQuotedWords(t *testing.T) {
	testRegex(t, QuotedFragment, `'arg1 arg2 "arg 3"`, []string{`arg1`, `arg2`, `"arg 3"`})
}

func TestSingleQuotedWords(t *testing.T) {
	testRegex(t, QuotedFragment, `arg1 arg2 'arg 3'`, []string{`arg1`, `arg2`, `'arg 3'`})
}

func TestQuotedWordsInTheMiddle(t *testing.T) {
	testRegex(t, QuotedFragment, `arg1 arg2 "arg 3" arg4   `, []string{`arg1`, `arg2`, `"arg 3"`, `arg4`})
}

func TestVariableParser(t *testing.T) {
	testRegex(t, VariableParser, `var`, []string{`var`})
	testRegex(t, VariableParser, `var.method`, []string{`var`, `method`})
	testRegex(t, VariableParser, `var[method]`, []string{`var`, `[method]`})
	testRegex(t, VariableParser, `var[method][0]`, []string{`var`, `[method]`, `[0]`})
	testRegex(t, VariableParser, `var["method"][0]`, []string{`var`, `["method"]`, `[0]`})
	testRegex(t, VariableParser, `var[method][0].method`, []string{`var`, `[method]`, `[0]`, `method`})
}
