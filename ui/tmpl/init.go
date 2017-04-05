package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"text/template"
	"fmt"
	"time"
	"io"
	"text/tabwriter"
)

var Templates = map[string]*template.Template{}
var Printers = map[string]func(w io.Writer, in interface{}){}

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

const(
	KWK_HOME = "kwk.co"
	MARGIN   = "  "
	TWOLINES = "\n\n"
)

func init() {
	// Aliases
	add("dashboard", style.Fmt16(style.Cyan, logo)+"{{. | listRoot }}", template.FuncMap{"listRoot": listRoot })

	add("snippet:updated", MARGIN+ "👍  Description updated:\n{{ .Description | blue }}\n\n", template.FuncMap{"blue": blue})
	add("api:not-found", "{{. | yellow }} Not found\n", template.FuncMap{"yellow": yellow})
	add("snippet:cloned", MARGIN +"👍  Cloned as {{.String | blue}}\n\n", template.FuncMap{"blue": blue})
	add("snippet:new", MARGIN+ "👍  {{. | blue }} created "+style.Icon_Pouch+"\n\n", template.FuncMap{"blue": blue})
	add("snippet:newprivate", MARGIN+ "👍  {{.String | blue }} created "+style.Icon_PrivatePouch+"\n\n", template.FuncMap{"blue": blue})
	add("snippet:cat", "{{.Snip}}\n", nil)
	add("snippet:edited", MARGIN + "👍  Successfully updated {{ .String | blue }}\n\n", template.FuncMap{"blue": blue})
	add("snippet:editing", "{{ \"Editing... \" | blue }}\nPlease edit file and save.\n - NB: After saving, no changes will be saved until running kwk edit <name> again.\n - Ctrl+C to abort.\n", template.FuncMap{"blue": blue})
	add("snippet:edit-prompt", "{{ .String | blue }} doesn't exist - would you like create it? [y/n] \n", template.FuncMap{"blue": blue})

	add("snippet:ambiguouscat", "That snippet is ambiguous please run it again with the extension:\n{{range .Items}}{{.String | blue }}\n{{ end }}", template.FuncMap{"blue": blue})
	add("snippet:list", "{{. | listPouch }}", template.FuncMap{"listPouch": listPouch })
	add("pouch:list-root", "{{. | listRoot }}", template.FuncMap{"listRoot": listRoot })

	add("snippet:tag", "{{.String | blue }} tagged.\n", template.FuncMap{"blue": blue})
	add("snippet:untag", "{{.String | blue }} untagged.\n", template.FuncMap{"blue": blue})
	add("snippet:renamed", "{{.originalName | blue }} renamed to {{.newName | blue }}\n", template.FuncMap{"blue": blue})
	add("snippet:patched", "{{.String | blue }} patched.\n", template.FuncMap{"blue": blue})

	add("snippet:check-delete", "Are you sure you want to delete snippet {{. | yellow }}? [y/n] ", template.FuncMap{"yellow": yellow})
	add("snippet:deleted", "Snippets {{. | blue }} deleted.\n", template.FuncMap{"blue": blue})
	add("snippet:not-deleted", "Snippets {{. | blue }} NOT deleted.\n", template.FuncMap{"blue": blue})

	add("snippet:moved-root", "{{ .Quant | blue }} snippet(s) moved to root.\n", template.FuncMap{"blue": blue})
	add("snippet:moved-pouch", "{{ .Quant | blue }} snippet(s) moved to pouch {{ .Pouch | blue }}\n", template.FuncMap{"blue": blue})
	add("snippet:create-pouch", "{{ \"Would you like to create the snippet in a new pouch? [y/n] \" | yellow }} ?  ", template.FuncMap{"yellow": yellow})
	add("snippet:inspect", "{{ . | inspect }}", template.FuncMap{"inspect":  inspect })

	add("pouch:not-deleted", "{{. | blue }} was NOT deleted.\n", template.FuncMap{"blue": blue})
	add("pouch:deleted", "{{. | blue }} was deleted.\n", template.FuncMap{"blue": blue})

	add("pouch:check-delete", "Are you sure you want to delete pouch {{. | yellow }}? [y/n] ", template.FuncMap{"yellow": yellow})
	add("pouch:created", "Pouch: {{. | blue }} created.\n", template.FuncMap{"blue": blue})
	add("pouch:renamed", "Pouch: {{. | blue }} renamed.\n", template.FuncMap{"blue": blue})
	add("pouch:locked", "🔒  pouch {{. | blue }} locked.\n", template.FuncMap{"blue": blue})
	add("pouch:unlocked", "🔓  pouch {{. | blue }} unlocked and public\n.", template.FuncMap{"blue": blue})
	add("pouch:not-locked", "Pouch: {{. | blue }} NOT locked.\n", template.FuncMap{"blue": blue})
	add("pouch:check-unlock", "Are you sure you want pouch 👝  {{. | blue }} public ? [y/n] ", template.FuncMap{"blue": blue})

	// System
	add("system:upgraded", "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n   Successfully upgraded!  \n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n", nil)
	add("system:version", "kwk version:\n CLI: {{ .String | blue }} released {{ .Time | time }}\n API: {{ .Api.String | blue}}\n", template.FuncMap{"blue": blue, "time": humanTime })

	// Account
	add("account:signedup", "Welcome to kwk {{.Username | blue }}!\n You're signed in already.\n", template.FuncMap{"blue": blue})
	addColor("account:usernamefield", "Your Kwk Username: ", blue)
	addColor("account:passwordfield", "Your Password: ", blue)
	addColor("account:password-changed", "You password has been changed for next login.\n", blue)
	add("account:signedin", "Welcome back {{.Username | blue }}!\n", template.FuncMap{"blue": blue})
	addColor("account:signedout", "And you're signed out.\n", blue)
	add("account:profile", "You are: {{.Username | blue}}!\n", template.FuncMap{"blue": blue})
	add("account:reset-sent", "Password reset instructions have been sent to: {{ .  | blue }}\n" +
		"*****\n" +
		"Use `kwk change-password` once you have received instructions.\n" +
		"*****\n", template.FuncMap{"blue": blue})

	add("env:changed", style.InfoDeskPerson+"  {{ \"env.yml\" | blue }} set to: {{.| blue }}\n", template.FuncMap{"blue": blue})

	addColor("account:signup:email", "What's your email? ", blue)
	addColor("account:signup:username", "Choose a great username: ", blue)
	addColor("account:signup:password", "And enter a password (1 num, 1 cap, 8 chars): ", blue)
	addColor("account:signup:invite-code", "Your kwk invite code: ", blue)

	// errors
	add("validation:title", "{{. | yellow }}\n", template.FuncMap{"yellow": yellow})
	add("validation:multi-line", " - {{ .Desc | yellow }}\n", template.FuncMap{"yellow": yellow})
	add("validation:one-line", style.Warning + "  {{ .Desc | yellow }}\n", template.FuncMap{"yellow": yellow})

	add("api:not-authenticated", "{{ \"Please login to continue: kwk login\" | yellow }}\n", template.FuncMap{"yellow": yellow})
	add("api:not-implemented", "{{ \"The kwk cli is a greater version than supported by kwk API.\" | yellow }}\n", template.FuncMap{"yellow": yellow})
	add("api:denied", "{{ \"Permission denied\" | yellow }}\n", template.FuncMap{"yellow": yellow})
	addColor("api:error", "\n"+style.Fire + "  We have a code RED error. \n- To report type: kwk upload-errors \n- You can also try to upgrade: npm update kwkcli -g\n", red)
	addColor("api:not-available", style.Ambulance + "  Kwk is DOWN! Please try again in a bit.\n", yellow)
	add("api:exists", "{{ \"That item already exists.\" | yellow }}\n", template.FuncMap{"yellow": yellow})
	add("free-text", "{{.}}", nil)
	addPrinters()
}


func addPrinters() {
	Printers["search:alpha"] = alphaSearchResult
	Printers["dialog:choose"] = multiChoice
}

func humanTime(t int64) string {
	return style.Time(time.Unix(t, 0))
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

func multiChoice(w io.Writer, in interface{}) {
	list := in.([]*models.Snippet)
	fmt.Fprint(w, "\n")
	if len(list) == 1 {
		fmt.Fprintf(w, "%sDid you mean: %s? y/n\n\n", MARGIN, style.Fmt256(style.Color_PouchCyan, list[0].String()))
		return
	}
	t := tabwriter.NewWriter(w, 5, 1, 3, ' ', tabwriter.TabIndent)
	for i, v := range list {
		if i%3 == 0 {
			t.Write([]byte(MARGIN))
		}
		fmt256 := style.Fmt16(style.Cyan, i+1)
		t.Write([]byte(fmt.Sprintf("%s %s", fmt256, v.SnipName.String())))
		x := i + 1
		if x%3 == 0 {
			t.Write([]byte("\n"))
		} else {
			t.Write([]byte("\t"))
		}
	}
	t.Write([]byte("\n\n"))
	t.Flush()
	fmt.Fprint(w, MARGIN + style.Fmt256(style.Color_PouchCyan, "Please select a snippet: "))
}

type ColorFunc func(int interface{}) string

func blue(in interface{}) string {
	return style.Fmt16(style.Cyan, fmt.Sprintf("%v", in))
}

func yellow(in interface{}) string {
	return style.Fmt16(style.Yellow, fmt.Sprintf("%v", in))
}

func red(in interface{}) string {
	return style.Fmt16(style.Red, fmt.Sprintf("%v", in))
}

func subdued(in interface{}) string {
	return style.Fmt16(style.Subdued, fmt.Sprintf("%v", in))
}
