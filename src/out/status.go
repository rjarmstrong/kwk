package out

import (
	"fmt"
	"github.com/rjarmstrong/kwk/src/style"
	"io"
)

func Fatal(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, "%s%s  %s", style.Margin, style.Fire, fmt.Sprintf(format, args...))
}

func Warn(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, "\n%s%s  %s", style.Margin, style.Fmt256(style.ColorBrightRed, style.Warning), fmt.Sprintf(format, args...))
}

func Success(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, "%s%s  %s", style.Margin, style.IconTick, fmt.Sprintf(format, args...))
}

func Info(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, "\n%s%s  %s\n", style.Margin, style.Fmt256(style.ColorPouchCyan, style.Info), fmt.Sprintf(format, args...))
}

func Prompt(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, "\n%s%s %s\n%s%s ",
		style.Margin,
		style.Fmt256(style.ColorWeekGrey, style.Prompt),
		fmt.Sprintf(format, args...),
		style.Margin,
		style.Fmt256(style.ColorPouchCyan, style.Prompt),
	)
}
