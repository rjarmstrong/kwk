package routes

import (
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk/src/style"
	"github.com/urfave/cli"
	"os"
)

func FirstIs(name string) bool {
	return len(os.Args) > 1 && os.Args[1] == name
}

func ReplaceArg(match string, replacement string) {
	for i, v := range os.Args {
		if match == v {
			os.Args[i] = replacement
		}
	}
}

func SetupFlags(prefs *out.Prefs, ap *cli.App) *cli.App {
	ap.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "covert, x",
			Usage:       "Open browser in covert mode.",
			Destination: &prefs.Covert,
		},
		cli.BoolFlag{
			Name:        "naked, n",
			Usage:       "List without styles",
			Destination: &prefs.Naked,
		},
		cli.BoolFlag{
			Name:        "ansi",
			Usage:       "Prints ansi escape sequences for debugging purposes",
			Destination: &style.PrintAnsi,
		},
		cli.BoolFlag{
			Name:        "quiet, q",
			Usage:       "List names only",
			Destination: &prefs.Quiet,
		},
	}
	return ap
}
