package tokenize_test

import (
	"fmt"
	"testing"

	"github.com/hownowstephen/go-liquid/tokenize"
)

func TestTokenizer(t *testing.T) {

	z := tokenize.NewTokenizer(`hello world {{ this is liquid }}
	{% for a in "{{ b }}" %}
	would you believe me if I said {{ a | b | c | d e:f g:h }} is in {{ b }}
	{% endfor %}`)

	for i := 0; i < 100; i++ {
		tk, err := z.Next()
		if err != nil {
			fmt.Println("ERR", err)
			break
		}
		fmt.Printf("%d %#v\n", i, tk)
	}

}
