package gui

import "text/template"

var Templates = map[string]*template.Template{
	"alias:delete" : t("alias:delete", "Are you sure you want to delete {{.fullKey}} y/n?"),
	"error" : t("error", "{{.error}}"),
	"account:signedup" : t("account:signedup", "Welcome to kwk {{.username}}! You're signed in already."),
	"account:notloggedin" : t("account:notloggedin", "You are not logged in please log in: kwk login <username> <password>"),
	"account:usernamefield" : t("account:usernamefield", "Your Kwk Username: "),
	"account:passwordfield" : t("account:passwordfield", "Your Password: "),
	"account:signedin" : t("account:signedin", "Welcome back {{.username}}!"),
	"system:upgraded" : t("system:upgraded",
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n   Successfully upgraded!  \n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n"),
	"alias:notfound": t("alias:notfound", "kwklink: {{.fullKey}} not found\n"),
	"alias:new": t("alias:new", "{{.fullKey}}"),
}

func t(name string, templ string) *template.Template {
	return template.Must(template.New(name).Parse(templ))
}

/*
`	alias:notdeleted
	messages := []string{"without a scratch", "uninjured", "intact", "unaffected", "unharmed",
			"unscathed", "out of danger", "safe and sound", "unblemished", "alive and well"}
		rnd := time.Now().Nanosecond() % (len(messages) - 1)
		fmt.Printf("'%s' is %s.\n", fullKey, messages[rnd])
	`,

	system:upgraded
	 fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	     fmt.Println("      Successfully upgraded!")
	     fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

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
