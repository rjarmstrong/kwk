package out

import (
	"fmt"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/vwrite"
	"io"
)

var UserChooseUsername = vwrite.HandlerFunc(func(w io.Writer) {
	fmt.Fprint(w, "Choose a memorable username: ")
})

var UserChooseAPassword = vwrite.HandlerFunc(func(w io.Writer) {
	fmt.Fprint(w, "And enter a password (1 num, 1 cap, 8 chars): ")
})

var UserUsernameField = vwrite.HandlerFunc(func(w io.Writer) {
	fmt.Fprint(w, "Your kwk Username: ")
})

var UserPasswordField = vwrite.HandlerFunc(func(w io.Writer) {
	fmt.Fprint(w, "Your Password: ")
})

var UserEmailField = vwrite.HandlerFunc(func(w io.Writer) {
	fmt.Fprint(w, "Whats your email?: ")
})

var UserInviteTokenField = vwrite.HandlerFunc(func(w io.Writer) {
	fmt.Fprint(w, "Your kwk invite code:  ")
})

var UserPasswordChanged = vwrite.HandlerFunc(func(w io.Writer) {
	fmt.Fprintln(w, "Your password has been changed for next login.")
})

func UserProfile(u *types.User) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "You are: %s\n", u.Username)
	})
}
func UserSignedIn(username string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprintf(w, "Welcome back, %s!\n", username)
	})
}

var UserSignedOut = vwrite.HandlerFunc(func(w io.Writer) {
	fmt.Fprintln(w, "You are now signed out.")
})

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
