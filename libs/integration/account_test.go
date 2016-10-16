package integration

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/smartystreets/assertions/should"
	"bitbucket.com/sharingmachine/kwkcli/libs/rpc"
	"bytes"
	"bufio"
	_ "github.com/go-sql-driver/mysql"
)

func Test_App(t *testing.T) {
	Convey("ACCOUNT COMMANDS", t, func() {
		conn := rpc.Conn("127.0.0.1:7777");
		//defer conn.Close() // Not sure if we should be closing this??
		w := &bytes.Buffer{}
		reader := &bytes.Buffer{}
		r := bufio.NewReader(reader)
		kwk := createApp(conn, w, r)

		const (
			email="test@kwk.co"
			username="testuser"
			password="TestPassword1"
		)
		signup := func() {
			reader.WriteString(email + "\n")
			reader.WriteString(username + "\n")
			reader.WriteString(password + "\n")
			kwk.Run("signup")
		}

		Convey(`Profile`, func() {
			notLoggedIn := "You are not logged in please log in: kwk login <username> <password>\n"

			Convey(`SIGNUP`, func() {
				Convey(`Should signup`, func() {
					signup()
					So(w.String(), should.ContainSubstring, "Welcome to kwk")
					w.Reset()
				})

				Convey(`Should get profile`, func() {
					kwk.Run("me")
					So(w.String(), should.Equal, "You are: testuser!")
				})

				Convey(`Should signout`, func() {
					kwk.Run("signout")
					w.Reset()
					kwk.Run("me")
					So(w.String(), should.Equal, notLoggedIn)
				})

				Convey(`Should notify if email used`, func() {
					reader.WriteString("test@kwk.co\n")
					reader.WriteString("testuser2\n")
					reader.WriteString("TestPassword1\n")
					kwk.Run("signup")
					So(w.String(), should.ContainSubstring, "Email has been taken.")
					w.Reset()
				})

				Convey(`Should notify if username has been taken`, func() {
					reader.WriteString("test2@kwk.co\n")
					reader.WriteString("testuser\n")
					reader.WriteString("TestPassword1\n")
					kwk.Run("signup")
					So(w.String(), should.ContainSubstring, "Username has been taken.")
					w.Reset()
				})
			})

			Convey(`SIGNIN`, func() {
				wrongCreds := "Wrong username or password."

				Convey(`Should successfully signin`, func() {
					kwk.Run("signin", username, password)
					So(w.String(), should.Equal, "Welcome back testuser!")
					w.Reset()
				})
				Convey(`Should print profile`, func() {
					signup()
					w.Reset()
					kwk.Run("profile")
					So(w.String(), should.Equal, "You are: testuser!")
					kwk.Run("signout")
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
		})
	})
}
