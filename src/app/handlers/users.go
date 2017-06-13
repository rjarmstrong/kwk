package handlers

import (
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/errs"
	"github.com/rjarmstrong/kwk-types/valid"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"github.com/rjarmstrong/kwk/src/cli"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk/src/store"
	"github.com/rjarmstrong/kwk/src/style"
	"os"
)

const userDocName = "user"

type Users struct {
	client types.UsersClient
	doc    store.Doc
	vwrite.Writer
	out.Dialog
	cxf         cli.ContextFunc
	pr          *cli.UserWithToken
	prefs       *out.Prefs
	rootPrinter cli.RootPrinter
}

func NewUsers(pr *cli.UserWithToken, uc types.UsersClient, doc store.Doc, w vwrite.Writer,
	d out.Dialog, c cli.ContextFunc, prefs *out.Prefs, rp cli.RootPrinter) *Users {

	return &Users{
		prefs:       prefs,
		pr:          pr,
		client:      uc,
		doc:         doc,
		Writer:      w,
		Dialog:      d,
		cxf:         c,
		rootPrinter: rp,
	}
}

// TASK: Add client side validation check
// SignUp
func (c *Users) SignUp() error {
	res, _ := c.FormField(out.UserEmailField, false)
	email := res.Value.(string)
	if !valid.Test(email, valid.RgxEmail) {
		// check email
		out.Warn(os.Stdout, "Invalid email")
		res, _ = c.FormField(out.UserEmailField, false)
		email = res.Value.(string)
	}
	res, _ = c.FormField(out.UserChooseUsername, false)
	username := res.Value.(string)
	res, _ = c.FormField(out.UserChoosePassword, true)
	c.Write(out.FreeText(style.Margin))
	password := res.Value.(string)

	req := &types.SignUpRequest{Email: email, Username: username, Password: password}
	u, err := c.client.SignUp(c.cxf(), req)
	if err != nil {
		return err
	}
	err = c.doc.Upsert(userDocName, cli.UserWithToken{AccessToken: u.AccessToken, User: *u.User})
	if err != nil {
		return err
	}
	return c.EWrite(out.UserSignedUp(u.User.Username))
}

// LogIn
func (c *Users) SignIn(username string) error {
	if username == "" {
		res, _ := c.FormField(out.UserUsernameField, false)
		username = res.Value.(string)
	}
	res, _ := c.FormField(out.UserPasswordField, true)
	c.Write(out.FreeText(style.Margin))
	password := res.Value.(string)

	ures, err := c.client.SignIn(c.cxf(),
		&types.SignInRequest{Username: username, Password: password, PrivateView: c.prefs.PrivateView})
	if err != nil {
		return err
	}
	c.pr.User = *ures.User
	c.pr.AccessToken = ures.AccessToken

	err = c.doc.Upsert(userDocName, cli.UserWithToken{AccessToken: ures.AccessToken, User: *ures.User})
	if err != nil {
		return err
	}
	c.Write(out.UserSignedIn(ures.User.Username))
	return c.rootPrinter(ures.Root)
}

// LogOut
func (c *Users) SignOut() error {
	err := c.doc.DeleteAll()
	if err != nil {
		return err
	}
	return c.EWrite(out.UserSignedOut())
}

// TASK: Client side validation
// ChangePassword
func (c *Users) ChangePassword() error {
	p := &types.ChangeRequest{}
	res, _ := c.Dialog.FormField(out.FreeText("Please provide an email or username: "), false)
	p.Username = res.Value.(string)
	res, _ = c.FormField(out.FreeText("Current password: "), true)
	p.ExistingPassword = res.Value.(string)
	res, _ = c.FormField(out.FreeText("New password: "), true)
	p.NewPassword = res.Value.(string)
	_, err := c.client.ChangePassword(c.cxf(), p)
	if err != nil {
		return err
	}
	return c.EWrite(out.UserPasswordChanged())
}

// ForgotPassword
func (c *Users) ForgotPassword() error {
	res, _ := c.FormField(out.FreeText("Enter your email to reset your password:  "), false)
	email := res.Value.(string)
	req := &types.ResetRequest{Email: email}
	_, err := c.client.ResetPassword(c.cxf(), req)
	if err != nil {
		return err
	}
	return c.EWrite(out.UserPasswordResetSent(email))
}

// Profile
func (c *Users) Profile() error {
	if !c.pr.HasAccessToken() {
		return errs.NotAuthenticated
	}
	return c.EWrite(out.UserProfile(&c.pr.User))
}

// LoadPrincipal load the currently signed in user from cache.
func (c *Users) LoadPrincipal(pr *cli.UserWithToken) {
	c.doc.Get(userDocName, pr, 0)
}
