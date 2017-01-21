package app

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
)

type AccountCli struct {
	service  account.Manager
	settings config.Persister
	tmpl.Writer
	dlg.Dialog
}

func NewAccountCli(u account.Manager, s config.Persister, w tmpl.Writer, d dlg.Dialog) *AccountCli {
	return &AccountCli{service: u, settings: s, Writer: w, Dialog: d}
}

func (c *AccountCli) Get() {
	u := &models.User{}
	if err := c.settings.Get(models.ProfileFullKey, u, 0); err != nil {
		c.Render("api:not-authenticated", nil)
	} else {
		c.Render("account:profile", u)
	}
}

func (c *AccountCli) SignUp() {

	email := c.FormField("account:signup:email", nil, false).Value.(string)
	username := c.FormField("account:signup:username", nil, false).Value.(string)
	password := c.FormField("account:signup:password", nil, true).Value.(string)

	if u, err := c.service.SignUp(email, username, password); err != nil {
		c.HandleErr(err)
	} else {
		if len(u.Token) > 50 {
			c.settings.Upsert(models.ProfileFullKey, u)
			c.Render("account:signedup", u)
		}
	}
}

func (c *AccountCli) SignIn(username string, password string) {
	if username == "" {
		username = c.FormField("account:usernamefield", nil, false).Value.(string)
	}
	if password == "" {
		password = c.FormField("account:passwordfield", nil, true).Value.(string)
	}
	if u, err := c.service.SignIn(username, password); err != nil {
		c.HandleErr(err)
	} else {
		if len(u.Token) > 50 {
			c.settings.Upsert(models.ProfileFullKey, u)
			c.Render("account:signedin", u)
		}
	}
}

func (c *AccountCli) SignOut() {
	if err := c.service.Signout(); err != nil {
		c.HandleErr(err)
		return
	}
	if err := c.settings.Delete(models.ProfileFullKey); err != nil {
		c.Render("error", err)
		return
	}
	c.Render("account:signedout", nil)
}
