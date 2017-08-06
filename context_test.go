package liquid

import (
	"fmt"
	"testing"
	"time"
)

func TestVariables(t *testing.T) {
	ctx := newContext()

	ctx.Assign("string", "string")
	v := ctx.Get("string")
	if v != "string" {
		t.Fatal(fmt.Sprintf(`Assigned "string" to Context but Get returns %+v`, v))
	}

	ctx.Assign("num", 5)
	v = ctx.Get("num")
	if v != 5 {
		t.Fatal(fmt.Sprintf(`Assigned 5 to Context but Get returns %+v`, v))
	}

	now := time.Now()
	ctx.Assign("time", now)
	v = ctx.Get("time")
	if v != now {
		t.Fatal(fmt.Sprintf(`Assigned a time.Time to Context but Get returns %+v`, v))
	}

	ctx.Assign("nil", nil)
	v = ctx.Get("nil")
	if v != nil {
		t.Fatal(fmt.Sprintf(`Assigned nil to Context but Get returns %+v`, v))
	}
}

func TestVariablesNotExisting(t *testing.T) {
	ctx := newContext()
	v := ctx.Get("wat")
	if v != nil {
		t.Fatal(`Non-existent variable fetched from Context was not nil`)
	}
}
