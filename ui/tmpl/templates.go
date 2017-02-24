package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"github.com/dustin/go-humanize"
	_ "github.com/dustin/go-humanize"
	"github.com/rjarmstrong/tablewriter"
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"encoding/json"
	"text/template"
	"strings"
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"
)

var Templates = map[string]*template.Template{}

//const logo = `
//  ▋                         ▋
//  ▋                         ▋
//  ▋   ◢  ◥◣           ◢◤   ▋   ◢◤
//  ▋ ◤      ◥◣    ◢◤   ◢◤    ▋ ◤
//  ▋ ◣       ◥◣ ◢◤ ◥◣ ◢◤     ▋ ◣
//  ▋   ◥     ◢◤     ◢◤      ▋   ◥◣
//`

const logo = `
`

func init() {
	// Aliases
	add("dashboard", style.Fmt(style.Cyan, logo)+"{{. | listRoot }}", template.FuncMap{"listRoot": listRoot })

	add("snippet:updated", "Description updated:\n{{ .Description | blue }}", template.FuncMap{"blue": blue})
	add("api:not-found", "{{. | yellow }} Not found\n", template.FuncMap{"yellow": yellow})
	add("snippet:cloned", "Cloned as {{.Username}}/{{.FullName | blue}}\n", template.FuncMap{"blue": blue})
	add("snippet:new", "{{. | blue }} created "+style.OpenLock+"\n", template.FuncMap{"blue": blue})
	add("snippet:newprivate", "{{.FullName | blue }} created "+style.Lock+"\n", template.FuncMap{"blue": blue})
	add("snippet:cat", "{{.Snip | blue}}", template.FuncMap{"blue": blue})
	add("snippet:edited", "Successfully updated {{ .String | blue }}", template.FuncMap{"blue": blue})
	add("snippet:editing", "{{ \"Editing... \" | blue }}\nPlease edit file and save.\n - NB: After saving, no changes will be saved until running kwk edit <name> again.\n - Ctrl+C to abort.\n", template.FuncMap{"blue": blue})
	add("snippet:edit-prompt", "{{ .String | blue }} doesn't exist - would you like create it? [y/n] \n", template.FuncMap{"blue": blue})

	add("snippet:ambiguouscat", "That snippet is ambiguous please run it again with the extension:\n{{range .Items}}{{.FullName | blue }}\n{{ end }}", template.FuncMap{"blue": blue})
	add("snippet:list", "{{. | listSnippets }}", template.FuncMap{"listSnippets": listSnippets })
	add("pouch:list-root", "{{. | listRoot }}", template.FuncMap{"listRoot": listRoot })

	add("snippet:tag", "{{.FullName | blue }} tagged.", template.FuncMap{"blue": blue})
	add("snippet:untag", "{{.FullName | blue }} untagged.", template.FuncMap{"blue": blue})
	add("snippet:renamed", "{{.originalName | blue }} renamed to {{.newName | blue }}", template.FuncMap{"blue": blue})
	add("snippet:madeprivate", "{{.fullName | blue }} made private "+style.Lock, template.FuncMap{"blue": blue})
	add("snippet:patched", "{{.FullName | blue }} patched.", template.FuncMap{"blue": blue})

	add("snippet:check-delete", "Are you sure you want to delete snippet {{. | yellow }}? [y/n] ", template.FuncMap{"yellow": yellow})
	add("snippet:deleted", "Snippets {{. | blue }} deleted.", template.FuncMap{"blue": blue})
	add("snippet:not-deleted", "Snippets {{. | blue }} NOT deleted.", template.FuncMap{"blue": blue})

	add("snippet:moved-root", "{{ .Quant | blue }} snippet(s) moved to root.", template.FuncMap{"blue": blue})
	add("snippet:moved-pouch", "{{ .Quant | blue }} snippet(s) moved to pouch {{ .Pouch | blue }}", template.FuncMap{"blue": blue})
	add("snippet:create-pouch", "{{ \"Would you like to create the snippet in a new pouch? [y/n] \" | yellow }}?  ", template.FuncMap{"yellow": yellow})

	add("snippet:inspect",
		"\n{{range .Items}}"+"Name: {{.String}}"+"\nCreated: {{.Created}}"+"\nTags: {{range $index, $element := .Tags}}{{if $index}}, {{end}} {{$element}}{{ end }}"+"\nDescription: {{.Description}}"+"\nRun count: {{.RunCount}}"+"\nClone count: {{.CloneCount}}"+"\n{{ end }}\n\n", nil)

	add("pouch:not-deleted", "{{. | blue }} was NOT deleted.", template.FuncMap{"blue": blue})
	add("pouch:deleted", "{{. | blue }} was deleted.", template.FuncMap{"blue": blue})

	add("pouch:check-delete", "Are you sure you want to delete pouch {{. | yellow }}? [y/n] ", template.FuncMap{"yellow": yellow})
	add("pouch:created", "Pouch: {{. | blue }} created.", template.FuncMap{"blue": blue})
	add("pouch:renamed", "Pouch: {{. | blue }} renamed.", template.FuncMap{"blue": blue})
	add("pouch:locked", "🔒  pouch {{. | blue }} locked.", template.FuncMap{"blue": blue})
	add("pouch:unlocked", "🔓  pouch {{. | blue }} unlocked and public.", template.FuncMap{"blue": blue})
	add("pouch:not-locked", "Pouch: {{. | blue }} NOT locked.", template.FuncMap{"blue": blue})
	add("pouch:check-unlock", "Are you sure you want pouch 👝  {{. | blue }} public ? [y/n]", template.FuncMap{"blue": blue})

	// System
	add("system:upgraded", "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n   Successfully upgraded!  \n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n", nil)
	add("system:version", "kwk version:\n CLI: {{ .String | blue }} released {{ .Time | time }}\n API: {{ .Api.String | blue}}\n", template.FuncMap{"blue": blue, "time": humanTime })

	// Account
	add("account:signedup", "Welcome to kwk {{.Username | blue }}!\n You're signed in already.\n", template.FuncMap{"blue": blue})
	addColor("account:usernamefield", "Your Kwk Username: ", blue)
	addColor("account:passwordfield", "Your Password: ", blue)
	add("account:signedin", "Welcome back {{.Username | blue }}!\n", template.FuncMap{"blue": blue})
	addColor("account:signedout", "And you're signed out.\n", blue)
	add("account:profile", "You are: {{.Username | blue}}!\n", template.FuncMap{"blue": blue})

	add("dialog:choose", "{{. | multiChoice }}\n", template.FuncMap{"multiChoice": multiChoice})
	add("dialog:header", "{{.| blue }}\n", template.FuncMap{"blue": blue})

	add("env:changed", style.InfoDeskPerson+"  {{ \"env.yml\" | blue }} set to: {{.| blue }}\n", template.FuncMap{"blue": blue})

	addColor("account:signup:email", "What's your email? ", blue)
	addColor("account:signup:username", "Choose a great username: ", blue)
	addColor("account:signup:password", "And enter a password (1 num, 1 cap, 8 chars): ", blue)

	add("search:alpha", "\n\033[7m  \"{{ .Term }}\" found in {{ .Total }} results in {{ .Took }} ms  \033[0m\n\n{{range .Results}}{{ .Username }}{{ \"/\" }}{{ .Name | blue }}.{{ .Extension | subdued }}\n{{ . | formatSearchResult}}\n{{end}}", template.FuncMap{"formatSearchResult": alphaSearchResult, "blue": blue, "subdued": subdued})
	add("search:alphaSuggest", "\n\033[7m Suggestions: \033[0m\n\n{{range .Results}}{{ .Username }}{{ \"/\" }}{{ .Name | blue }}.{{ .Extension | subdued }}\n{{end}}\n", template.FuncMap{"blue": blue, "subdued": subdued})
	add("search:typeahead", "{{range .Results}}{{ .String }}\n{{end}}", nil)

	// errors
	add("validation:title", "{{. | yellow }}\n", template.FuncMap{"yellow": yellow})
	add("validation:multi-line", " - {{ .Desc | yellow }}\n", template.FuncMap{"yellow": yellow})
	add("validation:one-line", style.Warning+"  {{ .Desc | yellow }}\n", template.FuncMap{"yellow": yellow})

	add("api:not-authenticated", "{{ \"Please login to continue: kwk login\" | yellow }}\n", template.FuncMap{"yellow": yellow})
	add("api:not-implemented", "{{ \"The kwk cli is a greater version than supported by kwk API.\" | yellow }}\n", template.FuncMap{"yellow": yellow})
	add("api:denied", "{{ \"Permission denied\" | yellow }}\n", template.FuncMap{"yellow": yellow})
	addColor("api:error", "\n"+style.Fire+"  We have a code RED error. \n- To report type: kwk upload-errors \n- You can also try to upgrade: npm update kwkcli -g\n", red)
	addColor("api:not-available", style.Ambulance+"  Kwk is DOWN! Please try again in a bit.\n", yellow)
	add("api:exists", "{{ \"That item already exists.\" | yellow }}\n", template.FuncMap{"yellow": yellow})
	add("free-text", "{{.}}", nil)
}

func humanTime(t int64) string {
	return humanize.Time(time.Unix(t, 0))
}

func add(name string, templateText string, funcMap template.FuncMap) {
	t := template.New(name)
	if funcMap != nil {
		t.Funcs(funcMap)
	}
	Templates[name] = template.Must(t.Parse(templateText))
}

func addColor(name string, text string, color ColorFunc) {
	add(name, fmt.Sprintf("{{ %q | color }}", text), template.FuncMap{"color": color})
}

func multiChoice(list []models.Snippet) string {
	var options string
	for i, v := range list {
		options = options + fmt.Sprintf("[%s] %s   ", style.Fmt(style.Cyan, i+1), v.SnipName.String())
	}
	return options
}

func listHorizontal(l []interface{}) []byte {
	var buff bytes.Buffer
	w := tabwriter.NewWriter(&buff, 20, 3, 2, ' ', tabwriter.DiscardEmptyColumns)
	var item bytes.Buffer
	for i, v := range l {
		if i%5 == 0 {
			item.WriteString("   ")
		}
		if sn, ok := v.(*models.Snippet); ok {
			item.WriteString(statusString(sn))
			item.WriteString("  ")
			item.WriteString(style.Fmt(style.Cyan, sn.SnipName.Name))
			item.WriteString(style.Fmt(style.Subdued, "."+sn.SnipName.Ext))
			item.WriteString(" ")
		}
		if pch, ok := v.(*models.Pouch); ok {
			if models.Prefs().ListAll || !pch.MakePrivate {
				if pch.MakePrivate {
					item.WriteString("🔒")
				} else {
					item.WriteString("👝")
				}
				item.WriteString("  ")
				item.WriteString(pch.Name)
				item.WriteString(style.Fmt(style.Subdued, fmt.Sprintf(" (%d)", pch.SnipCount)))
			}
		}

		item.WriteString(" \t")
		x := i + 1
		if x%5 == 0 {
			item.WriteString("\n")
		}
		if x%20 == 0 {
			item.WriteString("\n")
		}
		fmt.Fprint(w, fmt.Sprintf("%s", item.String()))
		item.Reset()
	}
	w.Flush()
	return buff.Bytes()
}

var mainMarkers = map[string]string{
	"go": "func main() {",
}

type CodeLine struct {
	Margin string
	Body string
}

func chunkFormatSnippet(s *models.Snippet, expand int) string {
	if s.Snip == "" {
		s.Snip = "<empty>"
	}
	chunks := strings.Split(s.Snip, "\n")

	// Return any non code previews
	if s.Role == models.SnipRolePreferences {
		return `(Global prefs) 'kwk edit prefs'`
	} else if s.Role == models.SnipRoleEnvironment {
		return `(Local environment) 'kwk edit env'`
	} else if s.Ext == "url" {
		return uri(s.Snip)
	}

	code := []CodeLine{}
	// Add line numbers and pad
	for i, v := range chunks {
		code = append(code, CodeLine{
			Margin: style.FmtStart(style.Subdued, fmt.Sprintf("%3d ", i+1)),
			Body:   fmt.Sprintf("    %s", strings.Replace(v, "\t", "  ", -1)),
		})
	}
	lastLine := code[len(chunks)-1]

	// Add to preview starting from most important line
	marker := mainMarkers[s.Ext]
	if marker != "" {
		var clipped []CodeLine
		var from int
		for i, v := range code {
			if strings.Contains(v.Body, marker) {
				from = i
			}
			if from > 0 {
				clipped = append(clipped, v)
			}
		}
		code = clipped
	}

	crop := len(code) >= expand

	// crop width
	var preview []CodeLine
	if crop {
		preview = code[0:expand]
	} else {
		preview = code
	}

	// Add page tear and last line
	if crop && expand < len(code) {
		preview = append(preview, CodeLine{"----", strings.Repeat("-", 70)})
		preview = append(preview, lastLine)
	}

	buff := bytes.Buffer{}
	for _, v := range preview {
		// Pad
		var width = 60
		pad := width - len(v.Body)
		if pad > 0 {
			v.Body = v.Body + strings.Repeat(" ", pad)
		} else {
			v.Body = v.Body[0:width]
		}
		// Style
		m := style.Fmt256(style.GreyBg238, true, v.Margin)
		b := style.Fmt256(style.GreyBg236, true, v.Body)
		buff.WriteString(m)
		buff.WriteString(b)
		buff.WriteString("\n")
	}

	return buff.String()
}

func statusString(s *models.Snippet) string {
	if s.Ext == "url" {
		return "🌎" //"🌐"
	}
	if s.RunStatus == models.RunStatusSuccess {
		return "⚡" //style.Fmt(style.Green, "●") //"✓"//
	} else if s.RunStatus == models.RunStatusFail {
		return "🔥" //style.Fmt(style.Red, "●") //
	}
	return "📄" //"🔸"
}

//func listLong(l []interface{}) []byte {
//	var buff bytes.Buffer
//	w := tabwriter.NewWriter(&buff, 7, 1, 3, ' ', tabwriter.TabIndent)
//	var item bytes.Buffer
//	item.WriteString("   ")
//	item.WriteString(style.Fmt(style.Subdued, "   Name"))
//	item.WriteString(" ")
//	item.WriteString("\t")
//	item.WriteString(style.Fmt(style.Subdued, "Snippet"))
//	item.WriteString("\n")
//	for _, v := range l {
//		item.WriteString("   ")
//		if sn, ok := v.(*models.Snippet); ok {
//			item.WriteString(statusString(sn.RunStatus))
//			item.WriteString("  ")
//			item.WriteString(style.Fmt(style.Cyan, sn.SnipName.String()))
//			item.WriteString("\t")
//			item.WriteString(formatSnippet(sn))
//			item.WriteString("\t")
//		}
//		if pch, ok := v.(*models.Pouch); ok {
//			if models.Prefs().ListAll || !pch.MakePrivate {
//				if pch.MakePrivate {
//					item.WriteString("🔒")
//				} else {
//					item.WriteString("👝")
//				}
//				item.WriteString("  ")
//				item.WriteString(pch.Name)
//				item.WriteString(" ")
//				//item.WriteString(style.Fmt(style.DarkGrey, fmt.Sprintf("%d", pch.SnipCount)))
//				item.WriteString("\t")
//				item.WriteString("")
//			}
//		}
//		fmt.Fprint(w, fmt.Sprintf("%s", item.String()))
//		item.Reset()
//		fmt.Fprint(w, "\t\n")
//	}
//	w.Flush()
//	return buff.Bytes()
//}

func listRoot(r *models.Root) string {
	var buff bytes.Buffer

	buff.WriteString("\n")
	buff.WriteString(style.Fmt(style.Cyan, "   kwk.co/") + r.Username + "/\n")
	buff.WriteString("\n")

	var all []interface{}
	for _, v := range r.Snippets {
		all = append(all, v)
	}
	for _, v := range r.Pouches {
		if v.Name != "" {
			all = append(all, v)
		}
	}

	//buff.Write(listSnippets())
	buff.Write(listHorizontal(all))

	buff.WriteString(style.Fmt(style.Subdued, fmt.Sprintf("\n\n   %d/50 Pouches", len(r.Pouches)-1)))
	if models.ClientIsNew(r.LastUpdate) {
		buff.WriteString(style.Fmt(style.Subdued, fmt.Sprintf("          kwk auto-updated to %s %s", models.Client.Version, humanTime(r.LastUpdate))))
	} else {
		buff.WriteString("\n")
	}
	buff.WriteString("\n\n")
	for _, v := range r.Personal {
		if v.Name == "inbox" {
			if v.UnOpened > 0 {
				buff.WriteString(fmt.Sprintf(" 📬 Inbox %d", v.UnOpened))

			} else {
				buff.WriteString(" 📪  inbox")
			}
		} else if v.Name == "settings" {
			buff.WriteString("   ⚙  settings")
		}
	}
	buff.WriteString("\n\n")
	return buff.String()
}

func formatDescription(in string) string {
	w := strings.Split(style.WrapString(in, 25), "\n")
	for i, v :=  range w {
		w[i] = style.Fmt(style.Subdued, v)
	}
	join := strings.Join(w, "\n")
	if len(join) > 25 * 4 {
		join = join[0:(25*4)]
		i := strings.LastIndex(join, " ")
		if i > 0 {
			join = join[0:]
		}

	}
	return join
}

func listSnippets(list *models.SnippetList) string {
	buf := new(bytes.Buffer)
	buf.WriteString("\n")
	fmt.Fprint(buf, style.Fmt(style.Cyan, "kwk.co/"+list.Username+"/")+list.Pouch+"/\n\n")
	tbl := tablewriter.NewWriter(buf)
	tbl.SetHeader([]string{"Name", "Status", "Preview"})
	tbl.SetAutoWrapText(false)
	tbl.SetBorder(false)
	tbl.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	tbl.SetCenterSeparator("")
	tbl.SetColumnSeparator(" ")
	tbl.SetRowLine(true)
	tbl.SetAutoFormatHeaders(false)
	tbl.SetHeaderLine(true)
	tbl.SetColWidth(5)

	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for _, v := range list.Items {
		var executed string
		if v.RunStatusTime > 0 {
			executed = fmt.Sprintf("%s  %s", statusString(v), style.Fmt(style.Subdued, humanize.Time(time.Unix(v.RunStatusTime, 0))))
		} else {
			executed = statusString(v)
		}
		preview := chunkFormatSnippet(v, models.Prefs().ExpandLines)
		tbl.Append([]string{
			style.Fmt(style.Cyan, v.SnipName.String()) + "\n\n" + formatDescription(v.Description),
			executed + "\n" + style.Fmt(style.Subdued, fmt.Sprintf("↻ %2d", v.RunCount)),
			//strings.Join(v.Tags, ", "),
			preview,
		})
	}
	tbl.Render()
	if len(list.Items) == 0 {
		fmt.Println(style.Fmt(style.Yellow, "Create some snippets to fill this view!\n"))
	}
	fmt.Fprintf(buf, "\n%d of %d records\n\n", len(list.Items), list.Total)
	fmt.Fprint(buf, "\n\n")

	return buf.String()
}

func alphaSearchResult(result models.SearchResult) string {
	if result.Highlights == nil {
		result.Highlights = map[string]string{}
	}
	if result.Highlights["snip"] == "" {
		result.Highlights["snip"] = result.Snip
	}
	lines := highlightsToLines(result.Highlights)
	f := ""
	for _, line := range lines {
		f = f + line.Key[:4] + "\u2847  " + line.Line + "\n"
	}
	f = style.Fmt(style.Subdued, f)
	f = style.ColourSpan(style.Black, f, "<em>", "</em>", style.Subdued)
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

type ColorFunc func(int interface{}) string

func blue(in interface{}) string {
	return style.Fmt(style.Cyan, fmt.Sprintf("%v", in))
}

func yellow(in interface{}) string {
	return style.Fmt(style.Yellow, fmt.Sprintf("%v", in))
}

func red(in interface{}) string {
	return style.Fmt(style.Red, fmt.Sprintf("%v", in))
}

func subdued(in interface{}) string {
	return style.Fmt(style.Subdued, fmt.Sprintf("%v", in))
}

func uri(text string) string {
	text = strings.Replace(text, "https://", "", 1)
	text = strings.Replace(text, "http://", "", 1)
	text = strings.Replace(text, "www.", "", 1)
	if len(text) >= 40 {
		text = text[0:10] + "..." + text[len(text)-30:]
	}
	return text
}

func PrettyPrint(obj interface{}) {
	fmt.Println("")
	p, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Print(string(p))
	fmt.Print("\n\n")
}
