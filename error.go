package liquid

import (
	"errors"
	"fmt"
)

type Context interface {
	String() string
}

func LiquidError(err string, context Context) error {
	return errors.New(fmt.Sprintf("%v - %v", err, context.String()))
}
