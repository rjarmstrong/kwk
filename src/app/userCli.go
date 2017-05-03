package app

import (
	"bitbucket.com/sharingmachine/kwkcli/src/gokwk"
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/kwkcli/src/persist"
	"bitbucket.com/sharingmachine/kwkcli/src/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/src/ui/tmpl"
	"os"
)

type UserCli struct {
	acc  gokwk.Users
	conf persist.Persister
	tmpl.Writer
	dlg.Dialog
	dash *Dashboard
}

func NewAccountCli(u gokwk.Users, s persist.Persister, w tmpl.Writer, d dlg.Dialog, dash *Dashboard) *UserCli {
	return &UserCli{acc: u, conf: s, Writer: w, Dialog: d, dash: dash}
}

func (c *UserCli) Get() {
	if u, err := c.acc.Get(); err != nil {
		c.Render("api:not-authenticated", nil)
	} else {
		c.Render("account:profile", u)
	}
}

func (c *UserCli) SignUp() {

	res, _ := c.TemplateFormField("account:signup:email", nil, false)
	email := res.Value.(string)
	res, _ = c.TemplateFormField("account:signup:username", nil, false)
	username := res.Value.(string)
	res, _ = c.TemplateFormField("account:signup:password", nil, true)
	password := res.Value.(string)
	res, _ = c.TemplateFormField("account:signup:invite-code", nil, false)
	inviteCode := res.Value.(string)

	if u, err := c.acc.SignUp(email, username, password, inviteCode); err != nil {
		c.HandleErr(err)
	} else {
		if len(u.Token) > 50 {
			c.conf.Upsert(models.ProfileFullKey, u)
			c.Render("account:signedup", u)
		}
	}
}

func (c *UserCli) SignIn(username string, password string) {
	if username == "" {
		res, _ := c.TemplateFormField("account:usernamefield", nil, false)
		username = res.Value.(string)

	}
	if password == "" {
		res, _ := c.TemplateFormField("account:passwordfield", nil, true)
		password = res.Value.(string)
	}
	if u, err := c.acc.SignIn(username, password); err != nil {
		c.HandleErr(err)
	} else {
		if len(u.Token) > 50 {
			c.conf.Upsert(models.ProfileFullKey, u)
			c.Render("account:signedin", u)
			c.dash.GetWriter()(os.Stdout, "", nil)
		}
	}
}

func (c *UserCli) SignOut() {
	err := c.acc.Signout()
	if err != nil {
		c.HandleErr(err)
		return
	}
	c.Render("account:signedout", nil)
}

func (c *UserCli) ChangePassword() {
	p := models.ChangePasswordParams{}
	res, _ := c.Dialog.FormField("Please provide an email or username:")
	p.Email = res.Value.(string)
	res, _ = c.Dialog.TemplateFormField("account:passwordfield", nil, true)
	p.ExistingPassword = res.Value.(string)
	res, _ = c.Dialog.TemplateFormField("account:passwordfield", nil, true)
	p.NewPassword = res.Value.(string)
	_, err := c.acc.ChangePassword(p)
	if err != nil {
		c.HandleErr(err)
	}
	c.Render("account:password-changed", nil)
}

func (c *UserCli) ResetPassword(email string) {
	if email == "" {
		res, _ := c.Dialog.TemplateFormField("account:signup:email", nil, false)
		email = res.Value.(string)
	}
	_, err := c.acc.ResetPassword(email)
	if err != nil {
		c.HandleErr(err)
	}
	c.Render("account:reset-sent", email)
}
