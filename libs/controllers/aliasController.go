package controllers

import (
	"bufio"
	"os"
	"fmt"
	"github.com/kwk-links/kwk-cli/libs/services/gui"
	"github.com/kwk-links/kwk-cli/libs/services/aliases"
	"time"
	"github.com/kwk-links/kwk-cli/libs/services/system"
	"github.com/olekukonko/tablewriter"
	"strings"
	"github.com/dustin/go-humanize"
	"strconv"
	"github.com/kwk-links/kwk-cli/libs/models"
	"github.com/kwk-links/kwk-cli/libs/services/openers"
)

type AliasController struct {
	service aliases.IAliases
	openers openers.IOpen
}

func NewAliasController(a aliases.IAliases, o openers.IOpen) *AliasController {
	return &AliasController{service:a, openers:o}
}

func (a *AliasController) Get(fullKey string){
	if list, err := a.service.Get(fullKey); err != nil {
		fmt.Println(err)
	} else {
		a.handleMultiResponse(fullKey, list, []string{})
	}
}

func (a *AliasController) New(uri string, fullKey string) {
	if alias, err := a.service.Create(uri, fullKey); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(alias.FullKey)
	}
}

func (a *AliasController) Edit(fullKey string) {
	if list, err := a.service.Get(fullKey); err != nil {
		fmt.Println(err)
	} else {
		// Currently not using edit on openers
		a.handleMultiResponse(fullKey, list, []string{})
	}
}

func (a *AliasController) Inspect(fullKey string) {
	if list, err := a.service.Get(fullKey); err != nil {
		fmt.Println(err)
	} else {
		system.PrettyPrint(list)
	}
}

func (a *AliasController) Delete(fullKey string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(gui.Colour(gui.LightBlue, "Are you sure you want to delete %s y/n? "), fullKey)
	yesNo, _, _ := reader.ReadRune()
	if string(yesNo) == "y" {
		if err := a.service.Delete(fullKey); err != nil {
			fmt.Println(err)
		}
		fmt.Println("Deleted")
	} else {
		messages := []string{"without a scratch", "uninjured", "intact", "unaffected", "unharmed",
			"unscathed", "out of danger", "safe and sound", "unblemished", "alive and well"}
		rnd := time.Now().Nanosecond() % (len(messages) - 1)
		fmt.Printf("'%s' is %s.\n", fullKey, messages[rnd])
	}
}

func (a *AliasController) Cat(fullKey string) {
	if list, err := a.service.Get(fullKey); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(list)
	}
}

func (a *AliasController) Patch(fullKey string, uri string){
	if _, err := a.service.Patch(fullKey, uri); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(gui.Colour(gui.LightBlue, "Patched %s"), fullKey)
	}
}

func (a *AliasController) Clone(fullKey string, newFullKey string){
	if alias, err := a.service.Clone(fullKey, newFullKey); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(gui.Colour(gui.LightBlue, "Cloned as %s"), alias.FullKey)
	}
}

func (a *AliasController) Rename(fullKey string, newKey string){
	if alias, err := a.service.Rename(fullKey, newKey); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(gui.Colour(gui.LightBlue, "Cloned as %s"), alias.FullKey)
	}
}


func (a *AliasController) Tag(fullKey string, tags ...string){
	if alias, err := a.service.Tag(fullKey, tags); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(gui.Colour(gui.LightBlue, "Tagged %s"), alias.FullKey)
	}
}

func (a *AliasController) UnTag(fullKey string, tags ...string){
	if alias, err := a.service.UnTag(fullKey, tags); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(gui.Colour(gui.LightBlue, "UnTagged %s"), alias.FullKey)
	}
}

func (a *AliasController) List(args ...string){
	var page, size int32
	var tags = []string{}
	for _, v := range args {
		if num, err := strconv.Atoi(v); err == nil {
			if page == 0 {
				page = int32(num)
			} else {
				size = int32(num)
			}
		} else {
			tags = append(tags, v)
		}
	}

	if list, err := a.service.List("richard", page, size, tags); err != nil {
		fmt.Println(err)
	} else {
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
	}
}


func (a *AliasController) handleMultiResponse(fullKey string, list *models.AliasList, args []string){
	if list.Total == 1 {
		a.openers.Open(&list.Items[0], args[1:])
	} else if list.Total > 1 {
		fmt.Println("Choose which is the correct key...")
		// print matches
		// read input
	} else {
		fmt.Println("Not found")
	}
}


//"notfound" : func(input interface{}) interface{}{
//	fmt.Printf(gui.Colour(gui.Yellow, "kwklink: '%s' not found\n"), input)
//	return nil
/*

