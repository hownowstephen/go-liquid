package ir

type Node interface {
	Render(Context) string
}
