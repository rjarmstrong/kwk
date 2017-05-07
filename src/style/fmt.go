package style

import (
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/types"
	"bytes"
	"fmt"
	"github.com/lunixbochs/vtclean"
	"strings"
)

// FStart starts an ansi escape sequence, should be terminated with style.End
func FStart(c types.AnsiCode, in interface{}) string {
	return fmt.Sprintf("\033[%dm%v", c, in)
}

func Fmt16(c types.AnsiCode, in interface{}) string {
	return fmtColor(c, in, "")
}

func Fmt256(c types.AnsiCode, in interface{}) string {
	return fmtColor(c, in, "38;5;")
}

// FmtFgBg formats both foreground and background. FStart is the recommended way to achieve the same effect.
func FmtFgBg(in string, fg types.AnsiCode, bg types.AnsiCode) string {
	r := fmt.Sprintf("%s38;5;%dm%s48;5;%dm%s%s", Start, fg, Start, bg, in, End)
	return r
}

func FPreview(in string, wrapAt int, lines int) string {
	if models.Prefs().DisablePreview {
		return ""
	}
	in = vtclean.Clean(in, false)
	return FBox(in, wrapAt, lines) + End
}

/*
 Creates a text box constrained by width (number of runes) and number of lines.
*/
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

/*
 Squeeze pushes text into a given width truncating the middle.
*/
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
		ansi := fmt.Sprintf("%s%s%dm%s%s", Start, ansiPattern, c, v, End)
		if models.Prefs() != nil && models.Prefs().PrintAnsi {
			a[i] = fmt.Sprintf("%q", ansi)
			continue
		}
		a[i] = ansi
	}
	return strings.Join(a, "\n")
}
