package dashboard

import (
	"text/template"
	"text/tabwriter"
	"strings"
	"io"
	"os"
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
)

const logo  = `
      _              _
     | |            | |
     | | ____      _| | __
 \\  | |/ /\ \ /\ / / |/ /
 //  |   <  \ V  V /|   <
     |_|\_\  \_/\_/ |_|\_\

`

func Layout(out io.Writer, templ string, data interface{}) {
	out.Write([]byte(style.Colour(style.LightBlue, logo)))

	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	w := tabwriter.NewWriter(out, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		// If the writer is closed, t.Execute will fail, and there's nothing
		// we can do to recover.
		if os.Getenv("CLI_TEMPLATE_ERROR_DEBUG") != "" {
			//fmt.Fprintf(ErrWriter, "CLI TEMPLATE ERROR: %#v\n", err)
		}
		return
	}
	w.Flush()
}
