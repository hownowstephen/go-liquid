package ir

type Condition []Conditional

func (c Condition) Render(ctx Context) string {
	for _, child := range c {
		if n, ok := child(ctx); ok {
			return n.Render(ctx)
		}
	}
	return ""
}

type Conditional func(Context) (Node, bool)

func If(left Node, op Op, right Node, body Node) Conditional {
	return func(ctx Context) (Node, bool) {
		if op(left, right) {
			return body, true
		}
		return nil, false
	}
}

func Else(body ...Node) Conditional {
	return func(ctx Context) (Node, bool) {
		return Block(body), true
	}
}
