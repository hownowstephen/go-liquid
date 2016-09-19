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
	testRegex(t, quotedFragmentRegexp, ``, []string{})
}

func TestQuote(t *testing.T) {
	testRegex(t, quotedFragmentRegexp, `"arg 1"`, []string{`"arg 1"`})
}

func TestWords(t *testing.T) {
	testRegex(t, quotedFragmentRegexp, `arg1, arg2`, []string{`arg1`, `arg2`})
}

func TestTags(t *testing.T) {
	testRegex(t, quotedFragmentRegexp, `<tr> </tr>`, []string{`<tr>`, `</tr>`})
	testRegex(t, quotedFragmentRegexp, `<tr></tr>`, []string{`<tr></tr>`})
	testRegex(t, quotedFragmentRegexp, `<style class="hello">' </style>`, []string{`<style`, `class="hello">`, `</style>`})

}

func TestDoubleQuotedWords(t *testing.T) {
	testRegex(t, quotedFragmentRegexp, `'arg1 arg2 "arg 3"`, []string{`arg1`, `arg2`, `"arg 3"`})
}

func TestSingleQuotedWords(t *testing.T) {
	testRegex(t, quotedFragmentRegexp, `arg1 arg2 'arg 3'`, []string{`arg1`, `arg2`, `'arg 3'`})
}

func TestQuotedWordsInTheMiddle(t *testing.T) {
	testRegex(t, quotedFragmentRegexp, `arg1 arg2 "arg 3" arg4   `, []string{`arg1`, `arg2`, `"arg 3"`, `arg4`})
}

func TestvariableParserRegexp(t *testing.T) {
	testRegex(t, variableParserRegexp, `var`, []string{`var`})
	testRegex(t, variableParserRegexp, `var.method`, []string{`var`, `method`})
	testRegex(t, variableParserRegexp, `var[method]`, []string{`var`, `[method]`})
	testRegex(t, variableParserRegexp, `var[method][0]`, []string{`var`, `[method]`, `[0]`})
	testRegex(t, variableParserRegexp, `var["method"][0]`, []string{`var`, `["method"]`, `[0]`})
	testRegex(t, variableParserRegexp, `var[method][0].method`, []string{`var`, `[method]`, `[0]`, `method`})
}
