package liquid

import (
	"fmt"
	"testing"
	"time"
)

func TestVariables(t *testing.T) {
	ctx := newContext()
	ctx.scopes.push()

	ctx.Assign("string", "string")
	v, err := ctx.Get("string")
	if err != nil {
		t.Fatal(err.Error())
	}
	if v != "string" {
		t.Fatal(fmt.Sprintf(`Assigned "string" to Context but Get returns %+v`, v))
	}

	ctx.Assign("num", 5)
	v, err = ctx.Get("num")
	if err != nil {
		t.Fatal(err.Error())
	}
	if v != 5 {
		t.Fatal(fmt.Sprintf(`Assigned 5 to Context but Get returns %+v`, v))
	}

	now := time.Now()
	ctx.Assign("time", now)
	v, err = ctx.Get("time")
	if err != nil {
		t.Fatal(err.Error())
	}
	if v != now {
		t.Fatal(fmt.Sprintf(`Assigned a time.Time to Context but Get returns %+v`, v))
	}

	ctx.Assign("nil", nil)
	v, err = ctx.Get("nil")
	if err != nil {
		t.Fatal(err.Error())
	}
	if v != nil {
		t.Fatal(fmt.Sprintf(`Assigned nil to Context but Get returns %+v`, v))
	}
}

func TestVariablesNotExisting(t *testing.T) {
	ctx := newContext()
	ctx.scopes.push()

	v, err := ctx.Get("wat")
	if err != nil {
		t.Fatal(err.Error())
	}
	if v != nil {
		t.Fatal(`Non-existent variable fetched from Context was not nil`)
	}
}

//scopeStack tests

func TestScopeStack(t *testing.T) {
	s := scopeStack{}

	vars, err := s.curr()
	if err == nil {
		t.Fatal(`empty scopeStack did not throw error on curr()`)
	}

	if vars != nil {
		t.Fatal(`empty scopeStack returned a Var`)
	}

	s.push()

	vars, err = s.curr()
	if err != nil {
		t.Fatal(err.Error())
	}
	if vars == nil {
		t.Fatal(`nil Var returned from scopeStack`)
	}

	if len(s) != 1 {
		t.Fatal(`unexpected scopeStack length`)
	}
}
