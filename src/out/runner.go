package out

import (
	"fmt"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/vwrite"
	"io"
)

func NotExecutable(s *types.Snippet) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprint(w, "\n")
		fmt.Fprintf(w, "%s", s.Content)
		fmt.Fprint(w, "\n\n")
	})
}
