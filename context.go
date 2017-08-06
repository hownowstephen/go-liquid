package liquid

import (
	"errors"
	"fmt"
)

var (
	ErrNoScope     = errors.New(`no scopes to pop`)
	ErrVarNotFound = errors.New(`variable not found`)
)

type Context struct {
	scopes scopeStack
}

func newContext() Context {
	s := scopeStack{}
	return Context{s}
}

func (c *Context) Assign(k string, v interface{}) error {
	scope, err := c.scopes.curr()
	if err != nil {
		return err
	}
	scope[k] = v
	return nil
}

func (c *Context) Get(k string) (interface{}, error) {
	if len(c.scopes) < 1 {
		return nil, ErrNoScope
	}

	for i := len(c.scopes) - 1; i >= 0; i-- {
		if val, ok := c.scopes[i][k]; ok {
			return val, nil
		}
	}
	return nil, ErrVarNotFound
}

func interfaceToExpression(v interface{}) Expression {
	switch v.(type) {
	case string:
		return stringExpr(v.(string))
	case int:
		return integerExpr(v.(int))
	case float64:
		return floatExpr(v.(float64))
	case []interface{}:
		return arrayExpr(v.([]interface{}))
	}
	panic(fmt.Sprintf("DONT UNDERSTAND %v"))
}

func (c *Context) FindVariable(e Expression) (Expression, error) {

	var key string
	switch e.(type) {
	case stringExpr:
		key = string(e.(stringExpr))
	case literalExpr:
		key = string(e.(literalExpr))
	default:
		return nil, fmt.Errorf("DUNNO WHAT TO DO WITH %v OMG", e)
	}

	value, err := c.Get(key)
	if err != nil {
		if err == ErrVarNotFound {
			return nil, ErrNotFound(key)
		}
		return nil, err
	}

	return interfaceToExpression(value), nil
}

func (c *Context) lookupAndEvaluate() {
}

type scopeStack []Vars

// Adds a new scope to the scopeStack
func (s *scopeStack) push() {
	*s = append(*s, Vars{})
}

// Removes scope from the scopeStack
func (s *scopeStack) pop() error {
	l := len(*s)
	if l < 1 {
		return ErrNoScope
	}

	*s = (*s)[:l-1]
	return nil
}

func (s *scopeStack) curr() (Vars, error) {
	l := len(*s)
	if l < 1 {
		return nil, ErrNoScope
	}

	return (*s)[l-1], nil
}
