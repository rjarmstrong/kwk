package routes

import (
	"github.com/kwk-super-snippets/cli/src/app/handlers"
	"github.com/kwk-super-snippets/cli/src/runtime"
	"github.com/kwk-super-snippets/types/errs"
	"os"
	"strings"
)

func RunNode(node *runtime.ProcessNode, snippets *handlers.Snippets) error {
	if len(os.Args) < 3 {
		return errs.New(errs.CodeInvalidArgument,
			"Invalid kwk call '%+v' in app.\n Invoke snippets as follows: kwk run <uri>",
			strings.Join(os.Args, " "))
	}
	if os.Args[1] != "run" {
		return errs.New(errs.CodeInvalidArgument, "'run' keyword required as first arg within an app.")
	}
	return snippets.RunNode(node, os.Args[2], os.Args[3:])
}
