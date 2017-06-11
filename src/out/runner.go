package out

import (
	"fmt"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"io"
)

func NotExecutable(s *types.Snippet) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprint(w, "\n")
		fmt.Fprintf(w, "%s", s.Content)
		fmt.Fprint(w, "\n\n")
	})
}
