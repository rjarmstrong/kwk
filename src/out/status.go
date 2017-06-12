package out

import (
	"fmt"
	"github.com/rjarmstrong/kwk/src/style"
	"io"
)

func Fatal(w io.Writer, format string, args ...interface{}) {
	styl := append([]interface{}{style.Margin, style.Fire}, args)
	fmt.Fprintf(w, "%s%s"+format, styl)
}

func Warn(w io.Writer, format string, args ...interface{}) {
	styl := append([]interface{}{style.Margin, style.Warning}, args)
	fmt.Fprintf(w, "%s%s"+format, styl)
}

func Success(w io.Writer, format string, args ...interface{}) {
	styl := append([]interface{}{style.Margin, style.IconTick}, args)
	fmt.Fprintf(w, "%s%s"+format, styl)
}

func Info(w io.Writer, format string, args ...interface{}) {
	styl := append([]interface{}{style.Margin, style.Info}, args)
	fmt.Fprintf(w, "%s%s"+format, styl)
}
