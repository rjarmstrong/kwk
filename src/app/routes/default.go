package routes

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/app/handlers"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/urfave/cli"
	"strings"
)

func DefaultRoute(info *types.AppInfo, snippets *handlers.Snippets, eh errs.Handler) func(*cli.Context, string) {
	return func(c *cli.Context, firstArg string) {
		i := c.Args().Get(1)
		if strings.HasPrefix(firstArg, "@") {
			snippets.GetEra(firstArg)
			return
		}
		var err error
		switch i {
		case "version":
			fmt.Println(info.String())
		case "run":
			err = snippets.Run(c.Args().First(), []string(c.Args())[2:])
		case "r":
			err = snippets.Run(c.Args().First(), []string(c.Args())[2:])
		case "edit":
			err = snippets.Edit(c.Args().First())
		case "e":
			err = snippets.Edit(c.Args().First())
		default:
			err = snippets.InspectListOrRun(c.Args().First(), false, []string(c.Args())[1:]...)
		}
		if err != nil {
			eh.Handle(err)
		}
	}
}
