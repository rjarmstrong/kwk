package app

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/users"
	"bitbucket.com/sharingmachine/kwkcli/libs/models"
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
	u := &models.User{}
	if err := c.settings.Get(models.ProfileFullKey, u); err != nil {
		c.Render("account:notloggedin", nil)
	} else {
		c.Render("account:profile", u)
	}
}

func (c *AccountController) SignUp(){

	email := c.Field("account:signup:email", nil).Value.(string)
	username := c.Field("account:signup:username", nil).Value.(string)
	password := c.Field("account:signup:password", nil).Value.(string)

	if u, err := c.service.SignUp(email, username, password); err != nil {
		c.Render("error", err)
	} else {
		if len(u.Token) > 50 {
			c.settings.Upsert(models.ProfileFullKey, u)
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
			c.settings.Upsert(models.ProfileFullKey, u)
			c.Render("account:signedin", u)
		}
	}
}

func (c *AccountController) SignOut(){
	if err := c.service.Signout(); err != nil {
		c.Render("error", err)
		return
	}
	if err := c.settings.Delete(models.ProfileFullKey); err != nil {
		c.Render("error", err)
		return
	}
	c.Render("account:signedout", nil)
}

func (c *AccountController) ChangeDirectory(username string) {
    if err := c.settings.ChangeDirectory(username); err != nil {
	    c.Render("error", err)
    } else {
	    c.Render("account:cd", map[string]string{"username" : username})
    }
}