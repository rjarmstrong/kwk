package integration

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/smartystreets/assertions/should"
	"bitbucket.com/sharingmachine/kwkcli/libs/rpc"
	"bytes"
	"bufio"
)

func Test_App(t *testing.T) {
	Convey("ACCOUNT COMMANDS", t, func() {
		conn := rpc.Conn("127.0.0.1:7777");
		defer conn.Close()
		w := &bytes.Buffer{}
		reader := &bytes.Buffer{}
		r := bufio.NewReader(reader)
		kwk := createApp(conn, w, r)

		Convey(`Profile`, func() {
			notLoggedIn := "You are not logged in please log in: kwk login <username> <password>\n"

			Convey(`SIGNIN`, func() {
				wrongCreds := "Wrong username or password."

				Convey(`Should successfully signin`, func() {
					kwk.Run("signin", "richard", "D1llbuck")
					So(w.String(), should.Equal, "Welcome back richard!")
					w.Reset()
				})
				Convey(`Should not signin with wrong password`, func() {
					kwk.Run("signin", "richard", "D1llbuckWrong")
					So(w.String(), should.Equal, wrongCreds)
					w.Reset()
				})
				Convey(`Should not signin with wrong username`, func() {
					kwk.Run("signin", "richardWrong", "D1llbuck")
					So(w.String(), should.Equal, wrongCreds)
					kwk.Run("signout")
					w.Reset()
				})
				Convey(`When calling method with requires account should reply signin prompt`, func() {
					kwk.Run("ls")
					So(w.String(), should.Equal, notLoggedIn)
					w.Reset()
				})
			})

			Convey(`SIGNUP`, func() {
				Convey(`Should successfully signup`, func() {
					reader.WriteString("email@dill.com\n")
					reader.WriteString("username\n")
					reader.WriteString("Password1\n")
					kwk.Run("signup")
					So(w.String(), should.ContainSubstring, "Welcome to kwk")
					w.Reset()
				})
				Convey(`Should notify if email used`, func() {
					reader.WriteString("email@dill.com\n")
					kwk.Run("signup")
					So(w.String(), should.ContainSubstring, "Email is in use on Kwk.")
					w.Reset()
				})
				Convey(`Should notify if username has been taken`, func() {
					reader.WriteString("email@dill.com\n")
					reader.WriteString("rjarmstrong\n")
					kwk.Run("signup")
					So(w.String(), should.ContainSubstring, "Email is in use on Kwk.")
					w.Reset()
				})
			})

				//Convey(`Should print not logged in`, func() {
				//	kwk.Run("profile")
				//	So(w.String()(), should.Equal, notLoggedIn)
				//	w.Reset()
				//})
				//Convey(`Should print profile`, func() {
				//	kwk.Run("signin", "richard", "D1llbuck")
				//	w.Reset()
				//	kwk.Run("profile")
				//	So(w.String()(), should.Equal, "You are: richard!")
				//	kwk.Run("signout")
				//})
		})
	})
}
