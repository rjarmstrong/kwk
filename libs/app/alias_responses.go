package app

import (
	"fmt"
	"time"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"bufio"
	"os"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/api"
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
}
