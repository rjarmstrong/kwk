package app

import (
	"fmt"
	"time"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"bufio"
	"os"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/api"
	"strings"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
)

var aliasResponses = map[string]gui.Template{
	"inspect" : func(input interface{}) interface{} {
		if input != nil {
			system.PrettyPrint(input)
		} else {
			fmt.Println("Invalid kwklink")
		}
		return nil
	},
	"new" : func(input interface{}) interface{}{
		k := input.(*api.Alias)
		fmt.Println(k.FullKey)
		return nil
	},
	"cat" : func(input interface{}) interface{}{
		k := input.(*api.AliasList)
		fmt.Println(k)
		return nil
	},
	"notfound" : func(input interface{}) interface{}{
		fmt.Printf(gui.Colour(gui.Yellow, "kwklink: '%s' not found\n"), input)
		return nil
	},
	"patch" : func(input interface{}) interface{}{
		if input != nil {
			k := input.(*api.Alias)
			fmt.Printf(gui.Colour(gui.LightBlue, "Patched %s"), k.FullKey)
		} else {
			fmt.Println("Invalid kwklink")
		}
		return nil
	},
	// delete returns a boolean indicating whether the user agreed to delete or not.
	"delete" : func(input interface{}) interface{} {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf(gui.Colour(gui.LightBlue, "Are you sure you want to delete %s y/n? "), input)
		yesNo, _, _ := reader.ReadRune()
		return string(yesNo) == "y"
	},
	"deleted" : func(input interface{}) interface{}{
		fmt.Println("Deleted")
		return nil
	},
	"notdeleted": func(input interface{}) interface{}{
		messages := []string{"without a scratch", "uninjured", "intact", "unaffected", "unharmed",
			"unscathed", "out of danger", "safe and sound", "unblemished", "alive and well"}
		rnd := time.Now().Nanosecond() % (len(messages) - 1)
		fmt.Printf("'%s' is %s.\n", input, messages[rnd])
		return nil
	},
	/*
	Move to serverside
		originalKey := k.FullKey
			uri := k.Uri
			if c.Args().Get(1) != "" && c.Args().Get(2) != "" {
				uri = strings.Replace(uri, c.Args().Get(1), c.Args().Get(2), -1)
			}
			kwklink := ""
			if c.Args().Get(3) != "" {
				kwklink = c.Args().Get(3)
			}
			k = apiClient.Create(uri, kwklink)
	 */

	"clone": func(input interface{}) interface{}{
		k := input.(*api.Alias)
		if input != nil {
			fmt.Printf(gui.Colour(gui.LightBlue, "Cloned as %s"), k.FullKey)
		} else {
			fmt.Println("Invalid kwklink")
		}
		return nil
	},
	"tag": func(input interface{}) interface{}{
		fmt.Println("Tagged")
		return nil
	},
	"untag": func(input interface{}) interface{}{
		fmt.Println("UnTagged")
		return nil
	},
	"list": func(input interface{}) interface{}{
		list := input.(*api.AliasList)
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
		return nil
	},
}
