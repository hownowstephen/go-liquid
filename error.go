package liquid

import (
	"errors"
	"fmt"
)

type liquidContext interface {
	String() string
}

func LiquidError(err string, context liquidContext) error {
	return errors.New(fmt.Sprintf("%v - %v", err, context))
}

func ErrNotFound(variable string) error {
	return LiquidError(fmt.Sprintf("Liquid::ErrorNotFound %v", variable), nil)
}
