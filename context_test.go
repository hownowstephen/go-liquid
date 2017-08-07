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
	if err != ErrVarNotFound {
		t.Fatal(`ErrVarNotFound not returned after calling Context.Get() with a non-existant variable`)
	}
	if v != nil {
		t.Fatal(`Non-existent variable fetched from Context was not nil`)
	}
}

func TestHyphenatedVariable(t *testing.T) {
	ctx := newContext()
	ctx.scopes.push()

	err := ctx.Assign("oh-my", "godz")
	if err != nil {
		t.Fatal(err.Error())
	}

	v, err := ctx.Get("oh-my")
	if err != nil {
		t.Fatal(err.Error())
	}

	if v != "godz" {
		t.Fatal(fmt.Sprintf(`expected "godz", got %+v`, v))
	}
}

// Copied the name from the ruby implementation but it's badly named.
// This actually tests if you can _access_ a variable set in a higher scope.
func TestAddItemInOuterScope(t *testing.T) {
	ctx := newContext()
	ctx.scopes.push()

	err := ctx.Assign("test", "test")
	if err != nil {
		t.Fatal(err.Error())
	}

	ctx.scopes.push()

	v, err := ctx.Get("test")
	if err != nil {
		t.Fatal(err.Error())
	}

	if v != "test" {
		t.Fatal(fmt.Sprintf(`Expected "test" but got %+v when looking up variable defined in higher scope"`, v))
	}

	ctx.scopes.pop()
	v, err = ctx.Get("test")
	if err != nil {
		t.Fatal(err.Error())
	}
	if v != "test" {
		t.Fatal(fmt.Sprintf(`Expected "test" but got %+v when looking up variable`, v))
	}
}

func TestAddItemInInnerScope(t *testing.T) {
	ctx := newContext()
	ctx.scopes.push() // Pristine
	ctx.scopes.push() // To be defiled with our variables

	err := ctx.Assign("test", "test")
	if err != nil {
		t.Fatal(err.Error())
	}

	v, err := ctx.Get("test")
	if err != nil {
		t.Fatal(err.Error())
	}
	if v != "test" {
		t.Fatal(fmt.Sprintf(`Expected "test" but got %+v when looking up variable`, v))
	}

	ctx.scopes.pop()

	v, err = ctx.Get("test")
	if err != ErrVarNotFound {
		t.Fatal(fmt.Sprintf(`ErrorVarNotFound not thrown after popping scope that contained it: %v`, err.Error()))
	}
	if v != nil {
		t.Fatal(fmt.Sprintf(`Got %+v from scope higher scope than it was defined in!`, v))
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
