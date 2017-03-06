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
	if u, err := c.service.Get(); err != nil {
		c.Render("api:not-authenticated", nil)
	} else {
		c.Render("account:profile", u)
	}
}

func (c *AccountCli) SignUp() {

	email := c.TemplateFormField("account:signup:email", nil, false).Value.(string)
	username := c.TemplateFormField("account:signup:username", nil, false).Value.(string)
	password := c.TemplateFormField("account:signup:password", nil, true).Value.(string)
	inviteCode := c.TemplateFormField("account:signup:invite-code", nil, false).Value.(string)

	if u, err := c.service.SignUp(email, username, password, inviteCode); err != nil {
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
		username = c.TemplateFormField("account:usernamefield", nil, false).Value.(string)
	}
	if password == "" {
		password = c.TemplateFormField("account:passwordfield", nil, true).Value.(string)
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

func (c *AccountCli) ChangePassword() {
	p := models.ChangePasswordParams{}
	res := c.Dialog.FormField("Please provide an email or username:")
	p.Email = res.Value.(string)
	res = c.Dialog.TemplateFormField("account:passwordfield", nil, true)
	p.ExistingPassword = res.Value.(string)
	res = c.Dialog.TemplateFormField("account:passwordfield", nil, true)
	p.NewPassword = res.Value.(string)
	_, err := c.service.ChangePassword(p)
	if err != nil {
		c.HandleErr(err)
	}
	c.Render("account:password-changed", nil)
}

func (c *AccountCli) ResetPassword(email string) {
	if email == "" {
		res := c.Dialog.TemplateFormField("account:signup:email", nil , false)
		email = res.Value.(string)
	}
	_, err := c.service.ResetPassword(email)
	if err != nil {
		c.HandleErr(err)
	}
	c.Render("account:reset-sent", email)
}
