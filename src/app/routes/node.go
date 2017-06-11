package routes

import (
	"github.com/rjarmstrong/kwk/src/app/handlers"
	"github.com/rjarmstrong/kwk/src/cli"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk/src/runtime"
	"github.com/rjarmstrong/kwk-types/errs"
	"os"
	"strings"
)

func RunNode(pr cli.UserWithToken, prefs *out.Prefs, node *runtime.ProcessNode, snippets *handlers.Snippets) error {
	if len(os.Args) < 3 {
		return errs.New(errs.CodeInvalidArgument,
			"Invalid kwk call '%+v' in app.\n Invoke snippets as follows: kwk run <uri>",
			strings.Join(os.Args, " "))
	}
	if os.Args[1] != "run" {
		return errs.New(errs.CodeInvalidArgument, "'run' keyword required as first arg within an app.")
	}
	return snippets.RunNode(pr, prefs, node, os.Args[2], os.Args[3:])
}
