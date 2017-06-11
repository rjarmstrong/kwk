package out

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"io"
)

func Fatal(h vwrite.Handler) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s  ", style.Fire)
		h.Write(w)
	})
}

func Warn(h vwrite.Handler) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s  ", style.Warning)
		h.Write(w)
	})
}

func Success(h vwrite.Handler) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s%s  ", style.Margin, style.IconTick)
		h.Write(w)
	})
}

func Info(h vwrite.Handler) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s%s  ", style.Margin, style.Info)
		h.Write(w)
	})
}

func Question(h vwrite.Handler) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s%s  ", style.Margin, "CONFIRM|")
		h.Write(w)
	})
}
