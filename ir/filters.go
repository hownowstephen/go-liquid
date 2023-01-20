package ir

type Filter func(Context, string) string

func TidyURL(ctx Context, url string) string {
	return "https://tidy.com"
}
