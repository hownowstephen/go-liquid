package liquid

import (
	"fmt"
	"regexp"
	"strings"
)

type tk struct {
	name  string
	value string
}

const (
	// tokens
	tPipe        = "pipe"
	tDot         = "dot"
	tColon       = "colon"
	tComma       = "comma"
	tOpenSquare  = "open_square"
	tCloseSquare = "close_square"
	tOpenRound   = "open_round"
	tCloseRound  = "close_round"
	tQuestion    = "question"
	tDash        = "dash"

	// sequences
	tIdentifier          = "id"
	tSingleStringLiteral = "string"
	tDoubleStringLiteral = "string"
	tNumberLiteral       = "number"
	tDotDot              = "dotdot"
	tComparisonOperator  = "comparison"

	// magic
	tEndOfString = "end_of_string"
)

// EndOfString is a
var EndOfString = tk{tEndOfString, ""}

var specialTokens = map[uint8]string{
	'|': tPipe,
	'.': tDot,
	':': tColon,
	',': tComma,
	'[': tOpenSquare,
	']': tCloseSquare,
	'(': tOpenRound,
	')': tCloseRound,
	'?': tQuestion,
	'-': tDash,
}

type sequence struct {
	name  string
	regex *regexp.Regexp
}

// Types of sequences to look for, in priority order
var sequenceTypes = []sequence{
	{tComparisonOperator, regexp.MustCompile(`^==|!=|<>|<=?|>=?|contains`)},
	{tSingleStringLiteral, regexp.MustCompile(`^'[^\']*'`)},
	{tDoubleStringLiteral, regexp.MustCompile(`^"[^\"]*"`)},
	{tNumberLiteral, regexp.MustCompile(`^-?\d+(\.\d+)?`)},
	{tIdentifier, regexp.MustCompile(`^[a-zA-Z_][\w-]*\??`)},
	{tDotDot, regexp.MustCompile(`^\.\.`)},
}

var whitespace = regexp.MustCompile(`\s`)

// Lexer converts liquid-y strings into lexographic tokens
func Lexer(s string) ([]tk, error) {

	s = strings.TrimSpace(s)
	var tokens []tk

TokenLoop:
	for i := 0; i < len(s); i++ {
		t := s[i]

		if whitespace.Match([]byte{t}) {
			continue
		}

		for _, seq := range sequenceTypes {
			if match := seq.regex.FindString(s[i:]); match != "" {
				tokens = append(tokens, tk{seq.name, match})
				i += len(match) - 1
				continue TokenLoop
			}
		}

		if name, ok := specialTokens[t]; ok {
			tokens = append(tokens, tk{name, string(t)})
			continue
		}

		return tokens, fmt.Errorf("Unexpected character: %v", string(t))

	}

	tokens = append(tokens, EndOfString)

	return tokens, nil
}
