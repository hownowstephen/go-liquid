package liquid

// based on revision 3146d5c of github.com:shopify/liquid.git

import (
	"fmt"
	"regexp"
)

// Regular Expressions for parsing liquid tags
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
	quotedFragmentRegexp        = regexp.MustCompile(fmt.Sprintf(`%v|(?:[^\s,\|'"]|%v)+`, quotedStringRegexp, quotedStringRegexp))
	tagAttributesRegexp         = regexp.MustCompile(fmt.Sprintf(`(\w+)\s*\:\s*(%v)`, quotedFragmentRegexp))
	anyStartingTagRegexp        = regexp.MustCompile(`\{\{|\{\%`)
	partialTemplateParserRegexp = regexp.MustCompile(fmt.Sprintf(`(?ms)%v.*?%v|%v.*?%v`, tagStartRegexp, tagEndRegexp, variableStartRegexp, variableIncompleteEndRegexp))
	templateParserRegexp        = regexp.MustCompile(fmt.Sprintf(`(?ms)(%v|%v)`, partialTemplateParserRegexp, anyStartingTagRegexp))
	variableParserRegexp        = regexp.MustCompile(fmt.Sprintf(`\[[^\]]+\]|%v+\??`, variableSegmentRegexp))
)
