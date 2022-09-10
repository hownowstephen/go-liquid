package tokenize_test

import (
	"fmt"
	"testing"

	"github.com/hownowstephen/go-liquid/tokenize"
)

func TestLexer(t *testing.T) {

	l := tokenize.NewLexer(`100 200 300 400 -500 600.1 -70.00 | != "yolo" >= 'yooooo "yo" ooooo' "ahoy m'partner" potato[1].salad salad  `)

	for {
		token, err := l.Next()
		fmt.Println(token, err)
		if err != nil {
			break
		}
	}
}
