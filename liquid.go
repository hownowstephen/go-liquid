package liquid

// based on revision 3146d5c of github.com:shopify/liquid.git

import (
	"fmt"
	"regexp"
)

const (
	ArgumentSeparator          = ','
	FilterArgumentSeparator    = ':'
	VariableAttributeSeparator = '.'
)

// Liquid regexes
var (
	FilterSeparator       = regexp.MustCompile(`\|`)
	TagStart              = regexp.MustCompile(`\{\%`)
	TagEnd                = regexp.MustCompile(`\%\}`)
	VariableSignature     = regexp.MustCompile(`\(?[\w\-\.\[\]]\)?`)
	VariableSegment       = regexp.MustCompile(`[\w\-]`)
	VariableStart         = regexp.MustCompile(`\{\{`)
	VariableEnd           = regexp.MustCompile(`\}\}`)
	VariableIncompleteEnd = regexp.MustCompile(`\}\}?`)
	QuotedString          = regexp.MustCompile(`"[^"]*"|'[^']*'`)
	QuotedFragment        = regexp.MustCompile(fmt.Sprintf(`%v|(?:[^\s,\|'"]|%v)+`, QuotedString.String(), QuotedString.String())) // o
	TagAttributes         = regexp.MustCompile(fmt.Sprintf(`(\w+)\s*\:\s*(%v)`, QuotedFragment.String()))                          // o
	AnyStartingTag        = regexp.MustCompile(`\{\{|\{\%`)
	PartialTemplateParser = regexp.MustCompile(fmt.Sprintf(`(?m)%v.*?%v|%v.*?%v`, TagStart.String(), TagEnd.String(), VariableStart.String(), VariableIncompleteEnd.String())) // om
	TemplateParser        = regexp.MustCompile(fmt.Sprintf(`(?m)(%v|%v)`, PartialTemplateParser.String(), AnyStartingTag.String()))                                            // om
	VariableParser        = regexp.MustCompile(fmt.Sprintf(`\[[^\]]+\]|%v+\??`, VariableSegment.String()))                                                                     // o
)
