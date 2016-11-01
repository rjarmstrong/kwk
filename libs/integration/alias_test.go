package integration

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/smartystreets/assertions/should"
	"bytes"
	_ "github.com/go-sql-driver/mysql"
)

func Test_Alias(t *testing.T) {
	Convey("ALIAS COMMANDS", t, func() {
		w := &bytes.Buffer{}
		reader := &bytes.Buffer{}
		kwk := getApp(reader, w)
		cleanup()

		Convey(`NEW`, func() {
			Convey(`Not logged in when creating new`, func() {
				kwk.Run("logout")
				So(w.String(), should.Equal, "And you're signed out.\n")
				w.Reset()
				kwk.Run("new", "http://somelink.com")
				So(w.String(), should.Equal, notLoggedIn)
				w.Reset()
			})
			Convey(`Given no extension should prompt and assign choosen url`, func() {
				signup(reader, kwk)
				w.Reset()

				reader.WriteString("7\n")
				kwk.Run("new", "echo \"hello\"", "hello")
				So(lastLine(w.String()), should.Equal, "hello.url created.")
				w.Reset()
			})
			Convey(`Should create new url with an extension`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.go")
				So(lastLine(w.String()), should.Equal, "hello.go created.")
				w.Reset()
			})
			Convey(`Should inspect alias`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.go")
				So(lastLine(w.String()), should.Equal, "hello.go created.")
				w.Reset()
				kwk.Run("inspect", "hello.go")
				So(w.String(), should.Resemble, "Alias: testuser/hello.go\nRuntime: golang\nURI: echo \"hello\"\nVersion: 1")
				w.Reset()
			})
			Convey(`Should cat unambiguous alias`, func() {
				signup(reader, kwk)
				w.Reset()
				uri := "echo \"hello\""
				kwk.Run("new", uri, "hello.go")
				So(lastLine(w.String()), should.Equal, "hello.go created.")
				w.Reset()
				kwk.Run("cat", "hello")
				So(w.String(), should.Equal, uri)
				w.Reset()
			})

			Convey(`Should prompt cat ambiguous alias`, func() {
				signup(reader, kwk)
				w.Reset()
				uri := "echo \"hello\""
				uri2 := "console.log('hello');"
				kwk.Run("new", uri, "hello.go")
				kwk.Run("new", uri2, "hello.js")
				w.Reset()
				kwk.Run("cat", "hello")
				So(w.String(), should.Resemble, "That alias is ambiguous please run it again with the extension:\nhello.go\nhello.js\n")
				w.Reset()
			})
			Convey(`Should rename an alias`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				kwk.Run("mv", "hello.js", "dong.js")
				So(w.String(), should.Resemble, "hello.js renamed to dong.js")
				w.Reset()
			})
			Convey(`Should rename an alias and auto-add the extension if not given in the new key`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				kwk.Run("mv", "hello.js", "dong")
				So(w.String(), should.Resemble, "hello.js renamed to dong.js")
				w.Reset()
			})
			Convey(`Should rename an alias and even if no extension is given for the original key`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				kwk.Run("mv", "hello", "dong")
				So(w.String(), should.Resemble, "hello.js renamed to dong.js")
				w.Reset()
			})
		})
	})
}
