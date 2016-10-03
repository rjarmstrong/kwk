package system

import (
	"fmt"
	"io"
	"github.com/fatih/color"
)

func Help(w io.Writer, template string, data interface{}) {
	c := color.New(color.FgCyan).Add(color.Bold)
	c.Printf("\n ===================================================================== ")
	c.Printf("\n ~~~~~~~~~~~~~~~~~~~~~~~~   KWK Power Links.  ~~~~~~~~~~~~~~~~~~~~~~~~ \n\n")
	c.Printf(" The ultimate URI manager. Create short and memorable codes called\n")
	c.Printf(" `aliass` to store URLs, computer paths, AppLinks etc.\n\n")
	c.Printf(" Usage: kwk [alias|cmd] [subcmd] [args]\n")
	fmt.Print("\n e.g.: `kwk open got-spoilers` to open all G.O.T. spoiler websites.\n")

	c.Printf("\n Commands:\n")
	fmt.Print("    <alias,..>                      - Open a alias or tag or execute the alias in the default app or script                                                               runtime.\n")
	fmt.Print("    new        <uri> [name]           - Create a new alias, optionally provide a memorable name\n")

	fmt.Print("    list       [tag,..] [and|or|not]  - List aliass, filter by tags\n")
	fmt.Print("    search     [term] [tag]           - * Search aliass and their metadata by keyword, filter by tags\n")
	fmt.Print("    suggest    <uri>                  - * List suggested aliass or tags for the given uri\n")
	fmt.Print("    tag        <alias> [tag,..]     - Add tags to a alias\n")
	fmt.Print("    untag      <alias> [tag,..]     - Remove tags from a alias\n")
	fmt.Print("    inspect    <alias>              - Look at the details of a kwk link\n")
	fmt.Print("    title      <alias> [text]       - Used to make a title for the resource\n")
	fmt.Println()
	fmt.Print("    rename     <alias> <alias>    - Update alias name <old> <new>\n")
	fmt.Print("    patch      <alias> <uri>|<oldstr> <newstr> - Patch a uri completely or partially. Auto increments the version\n")
	fmt.Print("    edit       <alias>              - Edit a alias in default editor\n")
	fmt.Print("    delete     <alias>              - Deletes alias with warning prompt. Will give 404.\n")
	fmt.Print("    covert     <alias>              - Open in covert (incognito mode)\n")
	fmt.Print("    cat        <alias>              - Gets URIs without navigating. (Copies first to clipboard)\n")

	c.Printf("\n Community:\n")
	fmt.Print("    cd    	  [username|tag]         - *Navigate to another users or tag. Used when searching or listing \n")
	fmt.Print("    fork      <alias> <uri> [oldstr] [newstr] - makes a copy of a alias and optionally replaces a string. \n")
	fmt.Print("    comment    [tag][alias][text]   - *Comment on tag or alias\n")
	fmt.Print("    profile    [username]             - *View a profile ??:\n")
	fmt.Print("    share      [alias|tag] [handle] - *Share with someone with a given handle:\n")
	fmt.Print("                                        twitter, email, kwk username\n")

	c.Printf("\n Analytics:\n")
	fmt.Print("    stats      [alias][tag]         - *Get statistics and filter by alias or tag\n")

	c.Printf("\n Account:\n")
	fmt.Print("    login      <username><password>   - Login with secret key.\n")
	fmt.Print("    me                                - View profile create ascii profile.\n")
	fmt.Print("    logout                            - Clears locally cached secret key.\n")
	fmt.Print("    signup     <email> <password> <username>  - Sign-up with a username.\n")

	fmt.Print("\n\n  * Filter only Tags: today yesterday thisweek lastweek thismonth lastmonth thisyear lastyear")
	fmt.Print("\n ** aliass are case sensitive")

	fmt.Print("\n\n More Commands: `kwk [admin|device] help`")

	//Day II: fmt.Print("	lock       <alias> <pin>          - Lock a alias with a pin\n")
	//Day II: fmt.Print("	subscribe  <domain>	            - Subscribe with custom domain. Free for 30 days.\n")
	//Day II: fmt.Print("	rate <alias> 9	            - Subscribe with custom domain. Free for 30 days.\n")
	//Day II: fmt.Print("	note <alias> "I like this one	    - Subscribe with custom domain. Free for 30 days.\n")

	//fmt.Printf("\n Admin:\n")
	//fmt.Printf("	cache       ls                  - List locally cached aliass.\n")
	//fmt.Printf("	cache       clear               - Clears any locally cached data.\n")
	//fmt.Printf("	upgrade                    	- Downloads and upgrades kwk cli client.\n")
	//fmt.Printf("	config      warn  [on|off]      - Warns if attempting to open dodgy alias.\n")
	//fmt.Printf("	config      quiet [on|off]      - Prevents links from being printed to console.\n")
	//fmt.Printf("	version                    	\n")
	c.Printf("\n ===================================================================== \n\n")
}
