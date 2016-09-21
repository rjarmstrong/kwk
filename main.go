package main

import (
	//"os"
	//"github.com/kwk-links/kwk-cli/libs/app"
	//"github.com/kwk-links/kwk-cli/libs/api"
	//"github.com/kwk-links/kwk-cli/libs/system"
	//"github.com/kwk-links/kwk-cli/libs/gui"
	//"github.com/kwk-links/kwk-cli/libs/openers"
)
import (
	//"os"
	//"github.com/kwk-links/kwk-cli/libs/gui"
	//"github.com/kwk-links/kwk-cli/libs/openers"
	//"github.com/kwk-links/kwk-cli/libs/app"
	//"github.com/kwk-links/kwk-cli/libs/system"
	//"github.com/kwk-links/kwk-cli/libs/settings"
	//"github.com/kwk-links/kwk-cli/libs/api"
	"github.com/kwk-links/kwk-cli/libs/rpc"
	"fmt"
	"github.com/kwk-links/kwk-cli/libs/services/users"
)

func main() {
	//TODO:
	 //Change settings directory
	//Setup interaction templates
	//Create a System type
	// Move project to bitbucket path

	conn := rpc.Conn("127.0.0.1:7777");
	usrs := users.New(conn)
	u, err := usrs.SignIn("dilld", "D1llbuck")
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

	//os.Setenv("version", "v0.0.1")
	////s.SetVersion("v0.0.1")
	//templates := map[string]gui.Template{}
	//i := gui.NewInteraction(templates)
	//s := system.NewSystem()
	//sett := settings.NewSettings(s, "settings")
	//a := api.New(sett)
	//o := openers.NewOpener(s, a)
	//kwkApp := app.NewKwkApp(a, s, sett, i, o)
	//kwkApp.App.Run(os.Args)
}
