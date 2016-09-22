package app

import (
	"github.com/kwk-links/kwk-cli/libs/services/users"
	"github.com/kwk-links/kwk-cli/libs/services/gui"
	"github.com/kwk-links/kwk-cli/libs/services/settings"
)

type AccountController struct {
	service users.IUsers
	settings settings.ISettings
	gui.ITemplateWriter
	gui.IDialogues
}

func NewAccountController(u users.IUsers, s settings.ISettings, w gui.ITemplateWriter, d gui.IDialogues) *AccountController {
	return &AccountController{service:u, settings:s, ITemplateWriter: w, IDialogues: d}
}

func (c *AccountController) Get(){
	if u, err := c.service.Get(); err != nil {
		c.Render("error", err)
		c.Render("account:notloggedin", nil)
	} else {
		c.Render("account:profile", u)
	}
}

func (c *AccountController) SignUp(email string, username string, password string){
	if u, err := c.service.SignUp(email, username, password); err != nil {
		c.Render("error", err)
	} else {
		if len(u.Token) > 50 {
			c.settings.Upsert("me", u)
			c.Render("account:signedup", u.Username)
		}
	}
}

func (c *AccountController) SignIn(username string, password string){
	if username == "" {
		username = c.Field("account:usernamefield", nil).Value.(string)
	}
	if password == "" {
		password = c.Field("account:passwordfield", nil).Value.(string)
	}
	if u, err := c.service.SignIn(username, password); err != nil {
		c.Render("error", err)
	} else {
		if len(u.Token) > 50 {
			c.settings.Upsert("me", u)
			c.Field("account:signedin", u.Username)
		}
	}
}

func (c *AccountController) SignOut(){
	c.service.Signout()
}

func (c *AccountController) ChangeDirectory(username string) {
    if err := c.settings.ChangeDirectory(username); err != nil {
	    c.Render("error", err)
    } else {
	    c.Render("account:cd", map[string]string{"username" : username})
    }
}