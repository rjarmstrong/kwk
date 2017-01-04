package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"github.com/dustin/go-humanize"
	_ "github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	_ "github.com/olekukonko/tablewriter"
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"google.golang.org/grpc"
	"text/template"
	"encoding/json"
	"strings"
	"bytes"
	"fmt"
)

var Templates = map[string]*template.Template{}

func init() {
	// Aliases
	add("snippet:delete", "Are you sure you want to delete {{.FullKey}}? [y/n] ", nil)
	add("snippet:deleted", "{{.FullKey}} deleted.", nil)
	add("snippet:updated", "Description updated:\n{{ .Description | blue }}", template.FuncMap{"blue": formatBlue})
	add("snippet:notfound", "Snippet: {{.FullKey}} not found\n", nil)
	add("snippet:cloned", "Cloned as {{.Username}}/{{.FullKey | blue}}\n", template.FuncMap{"blue": formatBlue})
	add("snippet:new", "{{.FullKey}} created "+style.OpenLock+"\n", nil)
	add("snippet:newprivate", "{{.FullKey}} created "+style.Lock+"\n", nil)
	add("snippet:cat", "{{.Snip | blue}}", template.FuncMap{"blue": formatBlue})
	add("snippet:edited", "Successfully updated {{ .FullKey | blue }}", template.FuncMap{"blue": formatBlue})
	add("snippet:editing", "{{ \"Editing file in default editor.\" | blue }}\nPlease save and close to continue. Or Ctrl+C to abort.\n", template.FuncMap{"blue": formatBlue})

	add("snippet:ambiguouscat", "That snippet is ambiguous please run it again with the extension:\n{{range .Items}}{{.FullKey}}\n{{ end }}", nil)
	add("snippet:list", "{{. | listSnippets }}", template.FuncMap{"listSnippets": listSnippets })
	add("snippet:choose", "{{. | listAmbiguous }}", template.FuncMap{"listAmbiguous": listAmbiguous })
	add("snippet:chooseruntime", "{{. | listRuntimes}}", template.FuncMap{"listRuntimes": listRuntimes})
	add("snippet:tag", "{{.FullKey}} tagged.", nil)
	add("snippet:untag", "{{.FullKey}} untagged.", nil)
	add("snippet:renamed", "{{.fullKey}} renamed to {{.newFullKey}}", nil)
	add("snippet:madeprivate", "{{.fullKey | blue }} made private "+style.Lock, template.FuncMap{"blue": formatBlue})
	add("snippet:patched", "{{.FullKey}} patched.", nil)
	add("snippet:notdeleted", "{{.FullKey}} was pardoned.", nil)
	add("snippet:inspect",
		"\n{{range .Items}}"+
			"Name: {{.Username}}/{{.FullKey}}"+
			"\nRuntime: {{.Runtime}}"+
			"\nSnippet: {{.Uri}}"+
			"\nVersion: {{.Version}}"+
			"\nTags: {{range $index, $element := .Tags}}{{if $index}}, {{end}} {{$element}}{{ end }}"+
			"\nWeb: \033[4mhttp://aus.kwk.co/{{.Username}}/{{.FullKey}}\033[0m"+
			"\nDescription: {{.Description}}"+
			"\nRun count: {{.RunCount}}"+
			"\nClone count: {{.ForkCount}}"+
			"\n{{ end }}\n\n", nil)

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

	add("search:alpha", "\n\033[7m  \"{{ .Term }}\" found in {{ .Total }} results in {{ .Took }} ms  \033[0m\n\n{{range .Results}}{{ .Username }}{{ \"/\" }}{{ .Key | blue }}.{{ .Extension | subdued }}\n{{ . | formatSearchResult}}\n{{end}}", template.FuncMap{"formatSearchResult": alphaSearchResult, "blue": formatBlue, "subdued": formatSubdued})

	add("search:alphaSuggest", "\n\033[7m Suggestions: \033[0m\n\n{{range .Results}}{{ .Username }}{{ \"/\" }}{{ .Key | blue }}.{{ .Extension | subdued }}\n{{end}}\n", template.FuncMap{"blue": formatBlue, "subdued": formatSubdued})

	// General
	add("error", "{{. | printError }}\n", template.FuncMap{"printError": printError})
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

func listRuntimes(list []interface{}) string {
	var options string
	for i, v := range list {
		m := v.(models.Match)
		options = options + fmt.Sprintf("%s %s   ", style.Colour(style.LightBlue, i+1), m.Runtime)
	}
	return options
}

func listAmbiguous(list []models.Snippet) string {
	var options string
	for i, v := range list {
		options = options + fmt.Sprintf("%s %s   ", style.Colour(style.LightBlue, i+1), v.FullKey)
	}
	return options
	return "hola"
}

func listSnippets(list *models.SnippetList) string {
	buf := new(bytes.Buffer)

	fmt.Fprint(buf, style.Colour(style.LightBlue, "\nkwk.co/"+list.Username+"\n\n"))

	tbl := tablewriter.NewWriter(buf)
	tbl.SetHeader([]string{"Name", "Version", "Snippet", "Tags", ""})
	tbl.SetAutoWrapText(false)
	tbl.SetBorder(false)
	tbl.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	tbl.SetCenterSeparator("")
	tbl.SetColumnSeparator("")
	tbl.SetAutoFormatHeaders(false)
	tbl.SetHeaderLine(true)
	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for _, v := range list.Items {
		var tags = []string{}
		for _, v := range v.Tags {
			if v == "error" {
				tags = append(tags, style.Colour(style.Pink, v))
			} else {
				tags = append(tags, v)
			}

		}

		var snip string
		var name string
		if v.Private {
			name = style.Colour(style.Subdued, "."+v.Key+"."+v.Extension)
			snip = style.Colour(style.Subdued, `<private>`)
		} else {
			name = style.Colour(style.LightBlue, v.Key) + style.Colour(style.Subdued, "."+v.Extension)
			snip = fmt.Sprintf("%s", formatUri(v.Snip))
		}

		tbl.Append([]string{
			name,
			fmt.Sprintf("%d", v.Version),
			snip,
			strings.Join(tags, ", "),
			humanize.Time(v.Created),
		})

	}
	tbl.Render()

	if len(list.Items) == 0 {
		fmt.Println(style.Colour(style.Yellow, "Create some snippets to fill this view!\n"))
	}
	fmt.Fprintf(buf, "\n%d of %d records\n\n", len(list.Items), list.Total)
	fmt.Fprint(buf, "\n\n")

	return buf.String()
}

func formatUri(uri string) string {
	uri = strings.Replace(uri, "https://", "", 1)
	uri = strings.Replace(uri, "http://", "", 1)
	uri = strings.Replace(uri, "www.", "", 1)
	uri = strings.Replace(uri, "\n", " ", -1)
	if len(uri) >= 40 {
		uri = uri[0:10] + style.Colour(style.Subdued, "...") + uri[len(uri)-30:]
	}
	if uri == "" {
		uri = "<empty>"
	}
	return uri
}

func alphaSearchResult(result models.SearchResult) string {
	if result.Highlights == nil {
		result.Highlights = map[string]string{}
	}
	if result.Highlights["uri"] == "" {
		result.Highlights["uri"] = result.Snip
	}
	lines := highlightsToLines(result.Highlights)
	f := ""
	for _, line := range lines {
		f = f + line.Key[:3] + "\u2847  " + line.Line + "\n"
	}
	f = style.Colour(style.Subdued, f)
	f = style.ColourSpan(40, f, "<em>", "</em>", style.Subdued)
	return f
}

func highlightsToLines(highlights map[string]string) []SearchResultLine {
	allLines := []SearchResultLine{}
	for k, v := range highlights {
		lines := strings.Split(v, "\n")
		for _, line := range lines {
			allLines = append(allLines, SearchResultLine{Key: k, Line: line})
		}
	}
	return allLines
}

type SearchResultLine struct {
	Key  string
	Line string
}

func formatSubdued(text string) string {
	return style.Colour(style.Subdued, text)
}

func formatBlue(text string) string {
	return style.Colour(style.LightBlue, text)
}

/*
`	snippet:notdeleted
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
		snippet:list
		fmt.Print(gui.Colour(gui.LightBlue, "\nkwk.co/" + "rjarmstrong/"))
			fmt.Printf(gui.Build(102, " ") + "%d of %d records\n\n", len(list.Items), list.Total)

			tbl := tablewriter.NewWriter(os.Stdout)
			tbl.SetHeader([]string{"Snippet", "Version", "URI", "Tags", ""})
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
