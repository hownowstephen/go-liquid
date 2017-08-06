package liquid

import (
	"errors"
	"fmt"
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
		return nil, errors.New(`No scopes to pop`)
	}

	var val interface{}
	for i := len(c.scopes) - 1; i >= 0; i-- {
		scope := c.scopes[i]
		val = scope[k]
		if val != nil {
			break
		}
	}
	return val, nil
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

	scope, err := c.scopes.curr()
	if err != nil {
		return nil, err
	}

	if value, ok := scope[key]; ok {
		// XXX: assumes flat variable structure. wrong
		return interfaceToExpression(value), nil
	}

	return nil, ErrNotFound(key)
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
		return errors.New(`No scopes to pop`)
	}

	*s = (*s)[:l-1]
	return nil
}

func (s *scopeStack) curr() (Vars, error) {
	l := len(*s)
	if l < 1 {
		return nil, errors.New(`No scopes to pop`)
	}

	return (*s)[l-1], nil
}
