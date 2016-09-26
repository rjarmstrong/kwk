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
)

var Templates = map[string]*template.Template{}

func init() {
	// Aliases
	add("alias:delete", "Are you sure you want to delete {{.fullKey}} y/n?", nil)
	add("alias:notfound", "kwklink: {{.fullKey}} not found\n", nil)
	add("alias:new", "{{.fullKey}}", nil)
	add("alias:list", "{{. | listAliases}}", template.FuncMap{"listAliases" : listAliases})

	// System
	add("system:upgraded", "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n   Successfully upgraded!  \n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n", nil)
	add("system:version", "kwkcli {{.version}}", nil)

	// Account
	add("account:signedup", "Welcome to kwk {{.username}}! You're signed in already.", nil)
	add("account:notloggedin", "You are not logged in please log in: kwk login <username> <password>", nil)
	add("account:usernamefield", "Your Kwk Username: ", nil)
	add("account:passwordfield", "Your Password: ", nil)
	add("account:signedin", "Welcome back {{.Username}}!", nil)
	add("account:profile", "You are: {{.Username}}!", nil)

	// General
	add("error", "{{.}}", nil)
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

func listAliases(list *models.AliasList) string {
	buf := new(bytes.Buffer)

	fmt.Fprint(buf, Colour(LightBlue, "\nkwk.co/" + "rjarmstrong/"))
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
		v.Uri = strings.Replace(v.Uri, "https://", "", 1)
		v.Uri = strings.Replace(v.Uri, "http://", "", 1)
		v.Uri = strings.Replace(v.Uri, "www.", "", 1)
		v.Uri = strings.Replace(v.Uri, "\n", " ", -1)
		if len(v.Uri) >= 40 {
			v.Uri = v.Uri[0:10] + Colour(Subdued, "...") + v.Uri[len(v.Uri) - 30:len(v.Uri)]
		}

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
