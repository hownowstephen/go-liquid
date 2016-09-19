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
	filterSeparatorRegexp       = regexp.MustCompile(`\|`)
	tagStartRegexp              = regexp.MustCompile(`\{\%`)
	tagEndRegexp                = regexp.MustCompile(`\%\}`)
	variableSignatureRegexp     = regexp.MustCompile(`\(?[\w\-\.\[\]]\)?`)
	variableSegmentRegexp       = regexp.MustCompile(`[\w\-]`)
	variableStartRegexp         = regexp.MustCompile(`\{\{`)
	variableEndRegexp           = regexp.MustCompile(`\}\}`)
	variableIncompleteEndRegexp = regexp.MustCompile(`\}\}?`)
	quotedStringRegexp          = regexp.MustCompile(`"[^"]*"|'[^']*'`)
	quotedFragmentRegexp        = regexp.MustCompile(fmt.Sprintf(`%v|(?:[^\s,\|'"]|%v)+`, quotedStringRegexp.String(), quotedStringRegexp.String())) // o
	tagAttributesRegexp         = regexp.MustCompile(fmt.Sprintf(`(\w+)\s*\:\s*(%v)`, quotedFragmentRegexp.String()))                                // o
	anyStartingTagRegexp        = regexp.MustCompile(`\{\{|\{\%`)
	partialTemplateParserRegexp = regexp.MustCompile(fmt.Sprintf(`(?m)%v.*?%v|%v.*?%v`, tagStartRegexp.String(), tagEndRegexp.String(), variableStartRegexp.String(), variableIncompleteEndRegexp.String())) // om
	templateParserRegexp        = regexp.MustCompile(fmt.Sprintf(`(?m)(%v|%v)`, partialTemplateParserRegexp.String(), anyStartingTagRegexp.String()))                                                        // om
	variableParserRegexp        = regexp.MustCompile(fmt.Sprintf(`\[[^\]]+\]|%v+\??`, variableSegmentRegexp.String()))                                                                                       // o
)
