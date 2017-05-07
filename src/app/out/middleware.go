package out

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/kwk-super-snippets/types/vwrite"
	"io"
)

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
		fmt.Fprintf(w, "%s%s  ", style.Margin, style.InfoDeskPerson)
		h.Write(w)
	})
}

func Question(h vwrite.Handler) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "%s%s  ", style.Margin, "CONFIRM|")
		h.Write(w)
	})
}
