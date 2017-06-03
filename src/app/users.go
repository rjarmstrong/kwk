package app

import (
	"github.com/kwk-super-snippets/cli/src/out"
	"github.com/kwk-super-snippets/cli/src/store"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/kwk-super-snippets/types/vwrite"
	"os"
)

type users struct {
	client types.UsersClient
	doc    store.Doc
	vwrite.Writer
	Dialog
	dash *Dashboard
}

type UserWithToken struct {
	AccessToken string `json:"access_token"`
	User        types.User
}

func (m *UserWithToken) HasAccessToken() bool {
	return m.AccessToken != ""
}

func NewUsers(u types.UsersClient, s store.Doc, w vwrite.Writer, d Dialog, dash *Dashboard) *users {
	return &users{client: u, doc: s, Writer: w, Dialog: d, dash: dash}
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

	req := &types.SignUpRequest{Email: email, Username: username, Password: password, InviteCode: inviteCode}
	u, err := c.client.SignUp(Ctx(), req)
	if err != nil {
		return err
	}
	if len(u.AccessToken) > 50 {
		err := c.doc.Upsert(cfg.UserDocName, UserWithToken{AccessToken: u.AccessToken, User: *u.User})
		if err != nil {
			return err
		}
		c.EWrite(out.UserSignedUp(u.User.Username))
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
	u, err := c.client.SignIn(Ctx(), &types.SignInRequest{Username: username, Password: password})
	if err != nil {
		return err
	}
	if len(u.AccessToken) > 50 {
		err := c.doc.Upsert(cfg.UserDocName, UserWithToken{AccessToken: u.AccessToken, User: *u.User})
		if err != nil {
			return err
		}
		c.Write(out.UserSignedIn(u.User.Username))
		c.dash.GetWriter()(os.Stdout, "", nil)
	}
	return nil
}

func (c *users) SignOut() error {
	err := c.doc.DeleteAll()
	if err != nil {
		return err
	}
	return c.EWrite(out.UserSignedOut)
}

func (c *users) ChangePassword() error {
	p := &types.ChangeRequest{}
	res, _ := c.Dialog.FormField(out.FreeText("Please provide an email or username: "), false)
	p.Username = res.Value.(string)

	res, _ = c.FormField(out.FreeText("Current password: "), true)
	p.ExistingPassword = res.Value.(string)
	res, _ = c.FormField(out.FreeText("New password: "), true)
	p.NewPassword = res.Value.(string)
	_, err := c.client.ChangePassword(Ctx(), p)
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
	req := &types.ResetRequest{Email: email}
	_, err := c.client.ResetPassword(Ctx(), req)
	if err != nil {
		return err
	}
	return c.EWrite(out.UserPasswordResetSent(email))
}

func (c *users) Profile() error {
	if principal == nil {
		return errs.NotAuthenticated
	}
	return c.EWrite(out.UserProfile(&principal.User))
}
