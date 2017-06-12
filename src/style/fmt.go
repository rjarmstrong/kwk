package style

import (
	"bytes"
	"fmt"
	"github.com/rjarmstrong/kwk-types"
	"strings"
)

// PrintAnsi is a setting to allow printing of ansi escape sequences for debugging purposes.
var PrintAnsi bool

// FStart starts an ansi escape sequence, should be terminated with style.End
func FStart(c types.AnsiCode, in interface{}) string {
	return fmt.Sprintf("\033[%dm%v", c, in)
}

// Fmt16 formats the input with 16 color terminal colors.
func Fmt16(c types.AnsiCode, in interface{}) string {
	return fmtColor(c, in, "")
}

// Fmt256 formats the input with 256 color terminal colors.
func Fmt256(c types.AnsiCode, in interface{}) string {
	return fmtColor(c, in, "38;5;")
}

// FmtFgBg formats both foreground and background. FStart is the recommended way to achieve the same effect.
func FmtFgBg(in string, fg types.AnsiCode, bg types.AnsiCode) string {
	r := fmt.Sprintf("%s38;5;%dm%s48;5;%dm%s%s", Esc, fg, Esc, bg, in, End)
	return r
}

// FBox creates a text box constrained by width (number of runes) and number of lines.
func FBox(in string, wrapAt int, lines int) string {
	in = strings.Replace(in, "\n", "  ", -1)
	in = strings.TrimSpace(in)
	var numRunes = 0
	b := bytes.Buffer{}
	var trim bool
	lineCount := 1
	for _, r := range in {
		numRunes++
		if trim && r == ' ' {
			continue
		}
		trim = false
		b.WriteRune(r)
		if numRunes%wrapAt == 0 {
			if lineCount >= lines {
				return strings.TrimSpace(b.String())
			}
			b.WriteString("\n")
			lineCount++
			trim = true
		}
	}
	return strings.TrimSpace(b.String())
}

// Squeeze pushes text into a given width truncating the middle.
func Squeeze(text string) string {
	if len(text) >= 40 {
		text = text[0:10] + "..." + text[len(text)-30:]
	}
	return " " + text
}

// fmtColor 'in' is the item to be formatted, ansiPattern is the short pattern which
// denotes the type of formatting. e.g. 256 colors is: 38;5;
// If the preference 'printansi' is set to false this method will have no effect.
func fmtColor(c types.AnsiCode, in interface{}, ansiPattern string) string {
	a := strings.Split(fmt.Sprintf("%v", in), "\n")
	for i, v := range a {
		ansi := fmt.Sprintf("%s%s%dm%s%s", Esc, ansiPattern, c, v, End)
		if PrintAnsi {
			a[i] = fmt.Sprintf("%q", ansi)
			continue
		}
		a[i] = ansi
	}
	return strings.Join(a, "\n")
}
