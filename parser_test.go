package liquid

// https://github.com/Shopify/liquid/blob/master/test/unit/parser_unit_test.rb

import "testing"

// checkConsume tests the result of a consume call
func checkConsume(t *testing.T, p *Parser, tType, want string) {
	got, err := p.consume(tType)
	if err != nil {
		t.Error(err)
		return
	}
	if got != want {
		t.Errorf("Bad token, want: %v, got: %v", want, got)
	}
}

func TestConsume(t *testing.T) {
	p, err := NewParser("wat: 7")
	if err != nil {
		t.Error(err)
		return
	}
	checkConsume(t, p, tIdentifier, "wat")
	checkConsume(t, p, tColon, ":")
	checkConsume(t, p, tNumberLiteral, "7")
}

func TestJump(t *testing.T) {
	p, err := NewParser("wat: 7")
	if err != nil {
		t.Error(err)
		return
	}
	p.jump(2)
	checkConsume(t, p, tNumberLiteral, "7")
}

// checkTryConsume checks the result of an attempted consume call
func checkTryConsume(t *testing.T, p *Parser, tType, want string, expectSuccess bool) {
	got, err := p.consume(tType)
	if err != nil && expectSuccess {
		t.Error(err)
		return
	}
	if got != want {
		t.Errorf("Bad token, want: %v, got: %v", want, got)
	}
}

func TestTryConsume(t *testing.T) {
	p, err := NewParser("wat: 7")
	if err != nil {
		t.Error(err)
		return
	}
	checkTryConsume(t, p, tIdentifier, "wat", true)
	checkTryConsume(t, p, tDot, "", false)
	checkTryConsume(t, p, tColon, ":", true)
	checkTryConsume(t, p, tNumberLiteral, "7", true)
}

// checkTryID checks that consuming an identifier succeeds
// it is just a special case of checkConsume and should probably be refactored
func checkTryID(t *testing.T, p *Parser, want string, expectSuccess bool) {
	got, err := p.consume(tIdentifier)
	if err != nil && expectSuccess {
		t.Error(err)
		return
	}
	if got != want && expectSuccess {
		t.Errorf("Bad token, want: %v, got: %v", want, got)
	}
}

func TestTryID(t *testing.T) {
	p, err := NewParser("wat 6 Peter Hegemon")
	if err != nil {
		t.Error(err)
		return
	}
	checkTryID(t, p, "wat", true)
	checkTryID(t, p, "endgame", false)
	checkConsume(t, p, tNumberLiteral, "6")
	checkTryID(t, p, "Peter", true)
	checkTryID(t, p, "Achilles", false)
}

// checkLook checks that a lookahead matches our expected value
func checkLook(t *testing.T, p *Parser, offset uint64, tType string, want bool) {
	got, err := p.lookahead(tType, offset)
	if err != nil {
		t.Error(err)
		return
	}
	if got != want {
		t.Errorf("Lookahead for %v failed, want: %v, got: %v", tType, want, got)
	}
}

func TestLook(t *testing.T) {
	p, err := NewParser("wat 6 Peter Hegemon")
	if err != nil {
		t.Error(err)
		return
	}
	checkLook(t, p, 0, tIdentifier, true)
	checkConsume(t, p, tIdentifier, "wat")
	checkLook(t, p, 0, tComparisonOperator, false)
	checkLook(t, p, 0, tNumberLiteral, true)
	checkLook(t, p, 1, tIdentifier, true)
	checkLook(t, p, 1, tNumberLiteral, false)
}

func checkExpression(t *testing.T, p *Parser, want string) {
	got, err := p.expression()
	if err != nil {
		t.Error(err)
		return
	}
	if got != want {
		t.Errorf("Expression consume failed, want: %v, got: %v", want, got)
	}
}

func TestExpressions(t *testing.T) {
	p, err := NewParser(`hi.there hi?[5].there? hi.there.bob`)
	if err != nil {
		t.Error(err)
		return
	}
	checkExpression(t, p, "hi.there")
	checkExpression(t, p, "hi?[5].there?")
	checkExpression(t, p, "hi.there.bob")

	p, err = NewParser(`567 6.0 'lol' "wut"`)
	if err != nil {
		t.Error(err)
		return
	}
	checkExpression(t, p, `567`)
	checkExpression(t, p, `6.0`)
	checkExpression(t, p, `'lol'`)
	checkExpression(t, p, `"wut"`)
}

func TestRanges(t *testing.T) {
	p, err := NewParser(`(5..7) (1.5..9.6) (young..old) (hi[5].wat..old)`)
	if err != nil {
		t.Error(err)
		return
	}

	checkExpression(t, p, `(5..7)`)
	checkExpression(t, p, `(1.5..9.6)`)
	checkExpression(t, p, `(young..old)`)
	checkExpression(t, p, `(hi[5].wat..old)`)
}

func checkArgument(t *testing.T, p *Parser, want string) {
	got, err := p.argument()
	if err != nil {
		t.Error(err)
		return
	}
	if got != want {
		t.Errorf("Argument consume failed, want: %v, got: %v", want, got)
	}
}

func TestArguments(t *testing.T) {
	p, err := NewParser(`filter: hi.there[5], keyarg: 7`)
	if err != nil {
		t.Error(err)
		return
	}

	checkConsume(t, p, tIdentifier, "filter")
	checkConsume(t, p, tColon, ":")
	checkArgument(t, p, "hi.there[5]")
	checkConsume(t, p, tComma, ",")
	checkArgument(t, p, "keyarg: 7")
}

func TestInvalidExpression(t *testing.T) {
	p, err := NewParser(`==`)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = p.expression()
	if err == nil {
		t.Errorf("Expected error when parsing for expression")
	}
}
