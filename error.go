package liquid

import (
	"fmt"
)

// ErrSyntax acts as a generic liquid syntax error
type ErrSyntax string

func (e ErrSyntax) Error() string {
	return fmt.Sprintf("Liquid syntax error: %v", string(e))
}

type liquidContext interface {
	String() string
}

// LiquidError prettifies an error with a liquid Context
func LiquidError(err string, context liquidContext) error {
	return fmt.Errorf("%v - %v", err, context)
}

// ErrNotFound wraps a missing variable error
func ErrNotFound(variable string) error {
	return LiquidError(fmt.Sprintf("Liquid::ErrorNotFound %v", variable), nil)
}
