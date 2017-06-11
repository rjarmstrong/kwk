package handlers

import (
	"github.com/rjarmstrong/kwk/src/cli"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk/src/store"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/errs"
	"github.com/rjarmstrong/kwk-types/vwrite"
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

// SignUp
func (c *Users) SignUp() error {
	// TASK: Add dynamic 'exists' checks
	res, _ := c.FormField(out.UserEmailField, false)
	email := res.Value.(string)
	res, _ = c.FormField(out.UserUsernameField, false)
	username := res.Value.(string)
	// TASK: Add client side validation check
	res, _ = c.FormField(out.UserPasswordField, true)
	password := res.Value.(string)
	//res, _ = c.FormField(out.UserInviteTokenField, false)
	//inviteCode := res.Value.(string)

	req := &types.SignUpRequest{Email: email, Username: username, Password: password}
	u, err := c.client.SignUp(c.cxf(), req)
	if err != nil {
		return err
	}
	if len(u.AccessToken) > 50 {
		err := c.doc.Upsert(userDocName, cli.UserWithToken{AccessToken: u.AccessToken, User: *u.User})
		if err != nil {
			return err
		}
		c.EWrite(out.UserSignedUp(u.User.Username))
	}
	return nil
}

// LogIn
func (c *Users) LogIn(username string, password string) error {
	if username == "" {
		res, _ := c.FormField(out.UserUsernameField, false)
		username = res.Value.(string)

	}
	if password == "" {
		res, _ := c.FormField(out.UserPasswordField, true)
		password = res.Value.(string)
	}
	ures, err := c.client.SignIn(c.cxf(),
		&types.SignInRequest{Username: username, Password: password, PrivateView: c.prefs.PrivateView})
	if err != nil {
		return err
	}
	c.pr.User = *ures.User
	c.pr.AccessToken = ures.AccessToken

	if len(ures.AccessToken) > 50 {
		err := c.doc.Upsert(userDocName, cli.UserWithToken{AccessToken: ures.AccessToken, User: *ures.User})
		if err != nil {
			return err
		}
		c.Write(out.UserSignedIn(ures.User.Username))
		// TASK: RootPrinter should be in UserSignedIn
		return c.rootPrinter(ures.Root)
	}
	return nil
}

// LogOut
func (c *Users) LogOut() error {
	err := c.doc.DeleteAll()
	if err != nil {
		return err
	}
	return c.EWrite(out.UserSignedOut)
}

// ChangePassword
func (c *Users) ChangePassword() error {
	p := &types.ChangeRequest{}
	res, _ := c.Dialog.FormField(out.FreeText("Please provide an email or username: "), false)
	p.Username = res.Value.(string)

	res, _ = c.FormField(out.FreeText("Current password: "), true)
	p.ExistingPassword = res.Value.(string)
	// TASK: Client side validation
	res, _ = c.FormField(out.FreeText("New password: "), true)
	p.NewPassword = res.Value.(string)
	_, err := c.client.ChangePassword(c.cxf(), p)
	if err != nil {
		return err
	}
	return c.EWrite(out.UserPasswordChanged)
}

// ResetPassword
func (c *Users) ResetPassword(email string) error {
	if email == "" {
		res, _ := c.FormField(out.FreeText("Enter your email to reset your password:  "), false)
		email = res.Value.(string)
	}
	req := &types.ResetRequest{Email: email}
	_, err := c.client.ResetPassword(c.cxf(), req)
	if err != nil {
		return err
	}
	return c.EWrite(out.UserPasswordResetSent(email))
}

// Profile
func (c *Users) Profile() error {
	if c.pr == nil {
		return errs.NotAuthenticated
	}
	return c.EWrite(out.UserProfile(&c.pr.User))
}

// LoadPrincipal load the currently signed in user from cache.
func (c *Users) LoadPrincipal(pr *cli.UserWithToken) {
	c.doc.Get(userDocName, pr, 0)
}
