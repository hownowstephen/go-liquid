package liquid

import (
	"reflect"
	"testing"
)

func checkLexerTokens(t *testing.T, raw string, want []Token) {
	got, err := Lexer(raw)
	if err != nil {
		t.Error("Got an error", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Tokens did not match, want: %v, got: %v", want, got)
	}
}

func TestStrings(t *testing.T) {
	checkLexerTokens(t, ` 'this is a test""' "wat 'lol'"`, []Token{
		{tSingleStringLiteral, `'this is a test""'`},
		{tDoubleStringLiteral, `"wat 'lol'"`},
		EndOfString,
	})
}

func TestInteger(t *testing.T) {
	checkLexerTokens(t, "hi 50", []Token{
		{tIdentifier, "hi"},
		{tNumberLiteral, "50"},
		EndOfString,
	})
}

func TestFloat(t *testing.T) {
	checkLexerTokens(t, "hi 5.0", []Token{
		{tIdentifier, "hi"},
		{tNumberLiteral, "5.0"},
		EndOfString,
	})
}

func TestComparison(t *testing.T) {
	checkLexerTokens(t, "== <> contains", []Token{
		{tComparisonOperator, "=="},
		{tComparisonOperator, "<>"},
		{tComparisonOperator, "contains"},
		EndOfString,
	})
}

func TestSpecials(t *testing.T) {
	checkLexerTokens(t, "| .:", []Token{
		{tPipe, "|"},
		{tDot, "."},
		{tColon, ":"},
		EndOfString,
	})

	checkLexerTokens(t, "[,]", []Token{
		{tOpenSquare, "["},
		{tComma, ","},
		{tCloseSquare, "]"},
		EndOfString,
	})
}

func TestFancyIdentifiers(t *testing.T) {
	checkLexerTokens(t, "hi five?", []Token{
		{tIdentifier, "hi"},
		{tIdentifier, "five?"},
		EndOfString,
	})
	checkLexerTokens(t, "2foo", []Token{
		{tNumberLiteral, "2"},
		{tIdentifier, "foo"},
		EndOfString,
	})
}

func TestWhitespace(t *testing.T) {
	checkLexerTokens(t, "five|\n\t ==", []Token{
		{tIdentifier, "five"},
		{tPipe, "|"},
		{tComparisonOperator, "=="},
		EndOfString,
	})
}

func TestUnexpectedCharacter(t *testing.T) {
	if _, err := Lexer("%"); err == nil {
		t.Error(`Should raise an error for '%'`)
	}
}
