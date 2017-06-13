package liquid

import "testing"

// integration/document_test.rb

func TestUnexpectedOuterTag(t *testing.T) {
	_, err := ParseTemplate("{% else %}")
	if err == nil {
		t.Error("Expected syntax error parsing template")
		return
	}
	if _, ok := err.(ErrSyntax); !ok {
		t.Errorf("Expected syntax error, got: %v", err)
	}

	if err.Error() != "Liquid syntax error: Unexpected outer 'else' tag" {
		t.Errorf("Wrong error text, got: %v", err)
	}
}

func TestUnknownTag(t *testing.T) {
	_, err := ParseTemplate("{% foo %}")
	if err == nil {
		t.Error("Expected syntax error parsing template")
		return
	}
	if _, ok := err.(ErrSyntax); !ok {
		t.Errorf("Expected syntax error, got: %v", err)
	}

	if err.Error() != "Liquid syntax error: Unknown tag 'foo'" {
		t.Errorf("Wrong error text, got: %v", err)
	}
}
