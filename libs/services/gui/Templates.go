package gui

import (
	"text/template"
	_ "strings"
	"fmt"
	_ "github.com/dustin/go-humanize"
	_ "github.com/olekukonko/tablewriter"
	_ "os"
	"bitbucket.com/sharingmachine/kwkcli/libs/models"
	"bytes"
	"github.com/olekukonko/tablewriter"
	"strings"
	"github.com/dustin/go-humanize"
	"google.golang.org/grpc"
	"encoding/json"
)

var Templates = map[string]*template.Template{}

func init() {
	// Aliases
	add("alias:delete", "Are you sure you want to delete {{.FullKey}}? y/n", nil)
	add("alias:deleted", "{{.FullKey}} deleted.", nil)
	add("alias:notfound", "alias: {{.FullKey}} not found\n", nil)
	add("alias:new", "{{.FullKey}} created.\n", nil)
	add("alias:cat", "{{.Uri}}", nil)
	add("alias:ambiguouscat", "That alias is ambiguous please run it again with the extension:\n{{range .Items}}{{.FullKey}}\n{{ end }}", nil)
	add("alias:list", "{{. | listAliases}}", template.FuncMap{"listAliases" : listAliases})
	add("alias:chooseruntime", "{{. | listRuntimes}}", template.FuncMap{"listRuntimes" : listRuntimes})
	add("alias:edited", "{{.FullKey}} updated.", nil)
	add("alias:tag", "{{.FullKey}} tagged.", nil)
	add("alias:untag", "{{.FullKey}} untagged.", nil)
	add("alias:renamed", "{{.fullKey}} renamed to {{.newFullKey}}", nil)
	add("alias:patched", "{{.FullKey}} patched.", nil)
	add("alias:notdeleted", "{{.FullKey}} was pardoned.", nil)
	add("alias:inspect", "{{range .Items}}Alias: {{.Username}}/{{.FullKey}}\nRuntime: {{.Runtime}}\nURI: {{.Uri}}\nVersion: {{.Version}}\nTags: {{range $index, $element := .Tags}}{{if $index}}, {{end}}{{$element}}{{ end }}{{ end }}", nil)
	//add("alias:chooseruntime", "{{.}}", nil)

	// System
	add("system:upgraded", "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n   Successfully upgraded!  \n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n", nil)
	add("system:version", "kwk {{.version}}'\n", nil)

	// Account
	add("account:signedup", "Welcome to kwk {{.Username}}! You're signed in already.\n", nil)
	add("account:notloggedin", "You are not logged in please log in: kwk login <username> <password>\n", nil)
	add("account:usernamefield", "Your Kwk Username: ", nil)
	add("account:passwordfield", "Your Password: ", nil)
	add("account:signedin", "Welcome back {{.Username}}!\n", nil)
	add("account:signedout", "And you're signed out.\n", nil)
	add("account:profile", "You are: {{.Username}}!\n", nil)

	add("account:signup:email", "What's your email? ", nil)
	add("account:signup:username", "Choose a great username: ", nil)
	add("account:signup:password", "And enter a password (1 num, 1 cap, 8 chars): ", nil)


	add("search:alpha", "\n\033[7m  \"{{ .Term }}\" found in {{ .Total }} results in {{ .Took }} ms  \033[0m\n\n{{range .Results}} {{ .Username }}{{ \"/\" }}{{ .Key | blue }}.{{ .Extension | subdued }}\n\n{{ .Highlights | formatHighlight}}\n\n{{end}}", template.FuncMap{"formatHighlight" : formatHighlights, "blue" : formatBlue, "subdued" : formatSubdued })

	// General
	add("error", "{{. | printError }}\n", template.FuncMap{"printError" : printError})
}

func printError(err error) string {
	e := new(models.Error)
	eString := strings.Replace(grpc.ErrorDesc(err), "\nError: ", "", -1)
	errBytes := []byte(eString)
	if e2 := json.Unmarshal(errBytes, e); e2 != nil {
		if eString == "Forbidden" {
			return render("account:notloggedin", nil)
		} else {
			return eString
		}
	}
	return e.Message
}

func render(name string, value interface{}) string {
	buf := new(bytes.Buffer)
	Templates[name].Execute(buf, value)
	return strings.Replace(buf.String(), "\n", "", -1)
}

func add(name string, templateText string, funcMap template.FuncMap) {
	t := template.New(name)
	if funcMap != nil {
		t.Funcs(funcMap)
	}
	Templates[name] = template.Must(t.Parse(templateText))
}

/*
`	alias:notdeleted
	messages := []string{"without a scratch", "uninjured", "intact", "unaffected", "unharmed",
			"unscathed", "out of danger", "safe and sound", "unblemished", "alive and well"}
		rnd := time.Now().Nanosecond() % (len(messages) - 1)
		fmt.Printf("'%s' is %s.\n", fullKey, messages[rnd])
	`,
	account:profile
	fmt.Println("~~~~~~ Your Profile ~~~~~~~~~")
		fmt.Println(gui.Build(2, gui.Space) + gui.Build(11, "~") + gui.Build(2, gui.Space))
		fmt.Println(gui.Build(6, gui.Space) + gui.Build(3, gui.UBlock) + gui.Build(6, gui.Space))
		fmt.Println(gui.Build(5, gui.Space) + gui.Build(5, gui.UBlock) + gui.Build(5, gui.Space))
		fmt.Println(gui.Build(6, gui.Space) + gui.Build(3, gui.UBlock) + gui.Build(6, gui.Space))
		fmt.Println(gui.Build(3, gui.Space) + gui.Build(9, gui.UBlock) + gui.Build(3, gui.Space))
		fmt.Println(gui.Build(3, gui.Space) + gui.Build(9, gui.UBlock) + gui.Build(3, gui.Space))
		fmt.Println(gui.Build(3, gui.Space) + gui.Build(9, gui.UBlock) + gui.Build(3, gui.Space))
		fmt.Println(gui.Build(2, gui.Space) + gui.Build(11, "~") + gui.Build(2, gui.Space))

		fmt.Printf("Email:      %v\n", u.Email)
		fmt.Printf("Username:   %v\n", u.Username)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")


		"github.com/olekukonko/tablewriter"
	"strings"
	"github.com/dustin/go-humanize"
		alias:list
		fmt.Print(gui.Colour(gui.LightBlue, "\nkwk.co/" + "rjarmstrong/"))
			fmt.Printf(gui.Build(102, " ") + "%d of %d records\n\n", len(list.Items), list.Total)

			tbl := tablewriter.NewWriter(os.Stdout)
			tbl.SetHeader([]string{"Alias", "Version", "URI", "Tags", ""})
			tbl.SetAutoWrapText(false)
			tbl.SetBorder(false)
			tbl.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
			tbl.SetCenterSeparator("")
			tbl.SetColumnSeparator("")
			tbl.SetAutoFormatHeaders(false)
			tbl.SetHeaderLine(true)
			tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

			for _, v := range list.Items {
				v.Uri = strings.Replace(v.Uri, "https://", "", 1)
				v.Uri = strings.Replace(v.Uri, "http://", "", 1)
				v.Uri = strings.Replace(v.Uri, "www.", "", 1)
				v.Uri = strings.Replace(v.Uri, "\n", " ", -1)
				if len(v.Uri) >= 40 {
					v.Uri = v.Uri[0:10] + gui.Colour(gui.Subdued, "...") + v.Uri[len(v.Uri) - 30:len(v.Uri)]
				}

				var tags = []string{}
				for _, v := range v.Tags {
					if v == "error" {
						tags = append(tags, gui.Colour(gui.Pink, v))
					} else {
						tags = append(tags, v)
					}

				}

				tbl.Append([]string{
					gui.Colour(gui.LightBlue, v.Key) + gui.Colour(gui.Subdued, "." + v.Extension),
					fmt.Sprintf("%d", v.Version),
					fmt.Sprintf("%s", v.Uri),
					strings.Join(tags, ", "),
					humanize.Time(v.Created),
				})

			}
			tbl.Render()

			if len(list.Items) == 0 {
				fmt.Println(gui.Colour(gui.Yellow, "No records on this page! Use a lower page number.\n"))
			} else {
				//gui.Colour(gui.Subdued, nextcmd)
				//nextcmd := fmt.Sprintf("For next page run: kwk list %v", 2)
			}
			if list.Size != 0 {
				fmt.Printf("\n %d of %d pages", list.Page, (list.Total / list.Size) + 1)
			}
			fmt.Print("\n\n")
 */

func listRuntimes(list []interface{}) string {
	var options string
	for i, v := range list {
		m := v.(models.Match)
		options = options + fmt.Sprintf("%d) %s  %d%%\n", i, m.Runtime, m.Score)
	}
	return options
}

func listAliases(list *models.AliasList) string {
	buf := new(bytes.Buffer)

	fmt.Fprint(buf, Colour(LightBlue, "\nkwk.co/" + list.Username))
	fmt.Fprintf(buf, Build(102, " ") + "%d of %d records\n\n", len(list.Items), list.Total)

	tbl := tablewriter.NewWriter(buf)
	tbl.SetHeader([]string{"Alias", "Version", "URI", "Tags", ""})
	tbl.SetAutoWrapText(false)
	tbl.SetBorder(false)
	tbl.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	tbl.SetCenterSeparator("")
	tbl.SetColumnSeparator("")
	tbl.SetAutoFormatHeaders(false)
	tbl.SetHeaderLine(true)
	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for _, v := range list.Items {
		v.Uri = formatUri(v.Uri)

		var tags = []string{}
		for _, v := range v.Tags {
			if v == "error" {
				tags = append(tags, Colour(Pink, v))
			} else {
				tags = append(tags, v)
			}

		}

		tbl.Append([]string{
			Colour(LightBlue, v.Key) + Colour(Subdued, "." + v.Extension),
			fmt.Sprintf("%d", v.Version),
			fmt.Sprintf("%s", v.Uri),
			strings.Join(tags, ", "),
			humanize.Time(v.Created),
		})

	}
	tbl.Render()

	if len(list.Items) == 0 {
		fmt.Println(Colour(Yellow, "No records on this page! Use a lower page number.\n"))
	} else {
		//Colour(Subdued, nextcmd)
		//nextcmd := fmt.Sprintf("For next page run: kwk list %v", 2)
	}
	if list.Size != 0 {
		fmt.Fprintf(buf, "\n %d of %d pages", list.Page, (list.Total / list.Size) + 1)
	}
	fmt.Fprint(buf, "\n\n")

	return buf.String()
}


func formatUri(uri string) string {
	uri = strings.Replace(uri, "https://", "", 1)
	uri = strings.Replace(uri, "http://", "", 1)
	uri = strings.Replace(uri, "www.", "", 1)
	uri = strings.Replace(uri, "\n", " ", -1)
	if len(uri) >= 40 {
		uri = uri[0:10] + Colour(Subdued, "...") + uri[len(uri) - 30:]
	}
	if uri == ""{
		uri = "<empty>"
	}
	return uri
}

func formatHighlights(highlights map[string]string) string {
	r := ""
	for _, v := range highlights {
		lines := strings.Split(v, "\n")
		for _, line := range lines {
			line =  "\u28FF    " + line
			line = ColourSpan(40, line, "<em>", "</em>", Subdued)
			r = r + line + "\n"
		}
	}
	return Colour(Subdued, r)
}

func formatSubdued(text string) string {
	return Colour(Subdued, text)
}

func formatBlue(text string) string {
	return Colour(LightBlue, text)
}