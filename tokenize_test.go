package liquid

import (
	"reflect"
	"testing"
)

func checkTokens(t *testing.T, tpl string, want []string) {
	tokenizer := NewTokenizer(tpl)
	var got []string
	for {
		token, err := tokenizer.Next()
		if token != "" {
			got = append(got, token)
		}

		if err != nil {
			break
		}
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Token list lengths didn't match, want: %v, got: %v", want, got)
	}

}

func TestTokenizeStrings(t *testing.T) {
	checkTokens(t, " ", []string{" "})
	checkTokens(t, "hello world", []string{"hello world"})
}

func TestTokenizeVariables(t *testing.T) {
	checkTokens(t, "{{funk}}", []string{"{{funk}}"})
	checkTokens(t, " {{funk}} ", []string{" ", "{{funk}}", " "})
	checkTokens(t, " {{funk}} {{so}} {{brother}} ", []string{" ", "{{funk}}", " ", "{{so}}", " ", "{{brother}}", " "})
	checkTokens(t, " {{  funk  }} ", []string{" ", "{{  funk  }}", " "})
}

func TestTokenizeBlocks(t *testing.T) {
	checkTokens(t, "{%comment%}", []string{"{%comment%}"})
	checkTokens(t, " {%comment%} ", []string{" ", "{%comment%}", " "})
	checkTokens(t, " {%comment%} {%endcomment%} ", []string{" ", "{%comment%}", " ", "{%endcomment%}", " "})
	checkTokens(t, " {% comment %} {% endcomment %} ", []string{" ", "{% comment %}", " ", "{% endcomment %}", " "})
}
func TestCalculateLineNumbersPerTokenWithProfiling(t *testing.T) {
	//     assert_equal [1],       tokenize_line_numbers("{{funk}}")
	//     assert_equal [1, 1, 1], tokenize_line_numbers(" {{funk}} ")
	//     assert_equal [1, 2, 2], tokenize_line_numbers("\n{{funk}}\n")
	//     assert_equal [1, 1, 3], tokenize_line_numbers(" {{\n funk \n}} ")
}
