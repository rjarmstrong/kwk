package integration

import (
	. "github.com/smartystreets/goconvey/convey"
	"bytes"
	"testing"
)

func Test_App(t *testing.T) {
	Convey("ACCOUNT COMMANDS", t, func() {
		w := &bytes.Buffer{}
		reader := &bytes.Buffer{}
		kwk := getApp(reader, w)

		Convey(`Profile`, func() {
			wrongCreds := "Wrong username or password.\n"

			Convey(`SIGNUP`, func() {
				Convey(`Should signup`, func() {
					signup(reader, kwk)
					So(w.String(), ShouldContainSubstring, "Welcome to kwk testuser! You're signed in already.\n")
					w.Reset()
				})

				Convey(`Should get profile`, func() {
					kwk.Run("me")
					So(w.String(), ShouldEqual, "You are: testuser!\n")
				})

				Convey(`Should signout`, func() {
					kwk.Run("signout")
					w.Reset()
					kwk.Run("me")
					So(w.String(), ShouldEqual, notLoggedIn)
				})

				Convey(`Should notify if email used`, func() {
					reader.WriteString("test@kwk.co\n")
					reader.WriteString("testuser2\n")
					reader.WriteString("TestPassword1\n")
					kwk.Run("signup")
					So(w.String(), ShouldContainSubstring, "Email has been taken.\n")
					w.Reset()
				})

				Convey(`Should notify if username has been taken`, func() {
					reader.WriteString("test2@kwk.co\n")
					reader.WriteString("testuser\n")
					reader.WriteString("TestPassword1\n")
					kwk.Run("signup")
					So(w.String(), ShouldContainSubstring, "Username has been taken.\n")
					w.Reset()
				})
			})

			Convey(`SIGNIN`, func() {

				Convey(`Should successfully signin`, func() {
					kwk.Run("signin", username, password)
					So(w.String(), ShouldEqual, "Welcome back testuser!\n")
					w.Reset()
				})
				Convey(`Should print profile`, func() {
					signup(reader, kwk)
					w.Reset()
					kwk.Run("profile")
					So(w.String(), ShouldEqual, "You are: testuser!\n")
					kwk.Run("signout")
				})
				Convey(`Should not signin with wrong password`, func() {
					kwk.Run("signin", "richard", "D1llbuckWrong")
					So(w.String(), ShouldEqual, wrongCreds)
					w.Reset()
				})
				Convey(`Should not signin with wrong username`, func() {
					kwk.Run("signin", "richardWrong", "D1llbuck")
					So(w.String(), ShouldEqual, wrongCreds)
					kwk.Run("signout")
					w.Reset()
				})
				Convey(`When calling method with requires account should reply signin prompt`, func() {
					kwk.Run("ls")
					So(w.String(), ShouldEqual, notLoggedIn)
					w.Reset()
				})
			})
		})
	})
}
