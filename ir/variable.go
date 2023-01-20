package ir

func Variable(markup string, filters ...Filter) variable {
	return variable{
		markup:  markup,
		filters: filters,
	}
}

type variable struct {
	markup  string
	filters []Filter
}

func (v variable) Render(Context) string {
	return v.markup + " lol"
}
