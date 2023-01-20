package ir

import (
	"fmt"
)

func Literal(v any) literal {
	return literal{
		val: v,
	}
}

type literal struct {
	val any
}

func (l literal) Render(Context) string {
	return fmt.Sprintf("%v", l.val)
}
