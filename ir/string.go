package ir

type String string

func (s String) Render(ctx Context) string {
	return string(s)
}
