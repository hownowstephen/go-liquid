package liquid

import "fmt"

type Context struct {
	scopes []Vars
}

func newContext() Context {
	return Context{Vars{}}
}

func (c *Context) Assign(k string, v interface{}) {
	//TODO: multiple stack handling
	c.scopes[][k] = v
}

func (c *Context) Get(k string) interface{} {
	return c.vars[k]
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

	if value, ok := c.vars[key]; ok {
		// XXX: assumes flat variable structure. wrong
		return interfaceToExpression(value), nil
	}

	return nil, ErrNotFound(key)
	// XXX: this doesn't handle scoping, everything is global for now :yay:
}

func (c *Context) lookupAndEvaluate() {
}


