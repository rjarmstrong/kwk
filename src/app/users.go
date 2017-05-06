package app

import (
	"bitbucket.com/sharingmachine/kwkcli/src/app/out"
	"bitbucket.com/sharingmachine/kwkcli/src/gokwk"
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/kwkcli/src/persist"
	"bitbucket.com/sharingmachine/types/vwrite"
	"os"
)

type users struct {
	acc  gokwk.Users
	conf persist.Persister
	vwrite.Writer
	Dialog
	dash *Dashboard
}

func NewAccount(u gokwk.Users, s persist.Persister, w vwrite.Writer, d Dialog, dash *Dashboard) *users {
	return &users{acc: u, conf: s, Writer: w, Dialog: d, dash: dash}
}

func (c *users) Get() error {
	u, err := c.acc.Get()
	if err != nil {
		return err
	}
	return c.EWrite(out.UserProfile(u))
}

func (c *users) SignUp() error {
	res, _ := c.FormField(out.UserEmailField, false)
	email := res.Value.(string)
	res, _ = c.FormField(out.UserUsernameField, false)
	username := res.Value.(string)
	res, _ = c.FormField(out.UserPasswordField, true)
	password := res.Value.(string)
	res, _ = c.FormField(out.UserInviteTokenField, false)
	inviteCode := res.Value.(string)

	u, err := c.acc.SignUp(email, username, password, inviteCode)
	if err != nil {
		return err
	}
	if len(u.Token) > 50 {
		err := c.conf.Upsert(models.ProfileFullKey, u)
		if err != nil {
			return err
		}
		c.EWrite(out.UserSignedUp(u.Username))
	}
	return nil
}

func (c *users) SignIn(username string, password string) error {
	if username == "" {
		res, _ := c.FormField(out.UserUsernameField, false)
		username = res.Value.(string)

	}
	if password == "" {
		res, _ := c.FormField(out.UserPasswordField, true)
		password = res.Value.(string)
	}
	u, err := c.acc.SignIn(username, password)
	if err != nil {
		return err
	}
	if len(u.Token) > 50 {
		err := c.conf.Upsert(models.ProfileFullKey, u)
		if err != nil {
			return err
		}
		c.Write(out.UserSignedIn(u.Username))
		c.dash.GetWriter()(os.Stdout, "", nil)
	}
	return nil
}

func (c *users) SignOut() error {
	err := c.acc.Signout()
	if err != nil {
		return err
	}
	return c.EWrite(out.UserSignedOut)
}

func (c *users) ChangePassword() error {
	p := models.ChangePasswordParams{}
	res, _ := c.Dialog.FormField(out.FreeText("Please provide an email or username: "), false)
	p.Username = res.Value.(string)

	res, _ = c.FormField(out.FreeText("Current password: "), true)
	p.ExistingPassword = res.Value.(string)
	res, _ = c.FormField(out.FreeText("New password: "), true)
	p.NewPassword = res.Value.(string)
	_, err := c.acc.ChangePassword(p)
	if err != nil {
		return err
	}
	return c.EWrite(out.UserPasswordChanged)
}

func (c *users) ResetPassword(email string) error {
	if email == "" {
		res, _ := c.FormField(out.FreeText("Enter your email to reset your password:  "), false)
		email = res.Value.(string)
	}
	_, err := c.acc.ResetPassword(email)
	if err != nil {
		return err
	}
	return c.EWrite(out.UserPasswordResetSent(email))
}
