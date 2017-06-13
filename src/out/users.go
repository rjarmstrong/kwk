package out

import (
	"fmt"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"github.com/rjarmstrong/kwk/src/style"
	"io"
)

var UserChooseUsername = vwrite.HandlerFunc(func(w io.Writer) {
	Prompt(w, "Choose a memorable username: ")
})

var UserChoosePassword = vwrite.HandlerFunc(func(w io.Writer) {
	Prompt(w, "And enter a password (1 num, 1 cap, 8 chars): ")
})

var UserUsernameField = vwrite.HandlerFunc(func(w io.Writer) {
	Prompt(w, "Your kwk Username: ")
})

var UserPasswordField = vwrite.HandlerFunc(func(w io.Writer) {
	Prompt(w, "Your Password: ")
})

var UserEmailField = vwrite.HandlerFunc(func(w io.Writer) {
	Prompt(w, "Whats your email?: ")
})

var UserInviteTokenField = vwrite.HandlerFunc(func(w io.Writer) {
	fmt.Fprint(w, "Your kwk invite code:  ")
})

func UserPasswordChanged() vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintln(w, "Your password has been changed for next login.")
	})
}

func UserProfile(u *types.User) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "You are: %s\n", u.Username)
	})
}
func UserSignedIn(username string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "\n%sWelcome back, %s!\n", style.Margin, username)
	})
}

func UserSignedOut() vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Info(w, "You are now signed out.")
		SignedOut().Write(w)
	})
}

func UserSignedUp(username string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Welcome to kwk %s!\n You're signed in already.\n", username)
	})
}

func UserPasswordResetSent(email string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Password reset instructions have been sent to: %s\n"+
			"Once you have received the token run: kwk change-password <token>\n"+
			"\n", email)
	})
}
