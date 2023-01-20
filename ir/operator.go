package ir

type Op func(left, right Node) bool

func GT(left, right Node) bool {
	return true
}

func Equals(left, right Node) bool {
	return false
}
