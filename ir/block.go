package ir

type Block []Node

func (b Block) Render(ctx Context) string {
	result := ""
	for _, n := range b {
		result += n.Render(ctx)
	}
	return result
}
