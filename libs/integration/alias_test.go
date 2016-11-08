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
			Convey(`Should create new blank alias if no content or alias is given`, func() {
				signup(reader, kwk)
				w.Reset()
				reader.WriteString("7\n")
				kwk.Run("new")
				So(lastLine(w.String()), should.ContainSubstring, ".url created.")
				w.Reset()
			})
		})

		Convey(`INSPECT`, func() {
			Convey(`Should inspect alias`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.go")
				So(lastLine(w.String()), should.Equal, "hello.go created.")
				w.Reset()
				kwk.Run("inspect", "hello.go")
				So(w.String(), should.Resemble, "Alias: testuser/hello.go\nRuntime: golang\nURI: echo \"hello\"\nVersion: 1\nTags: ")
				w.Reset()
			})
		})


		Convey(`CAT`, func() {
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
		})

		Convey(`RENAME`, func() {
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

		Convey(`DELETE`, func() {
			Convey(`Should prompt to delete alias and delete when 'y'`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()

				reader.WriteString("y\n")

				kwk.Run("rm", "hello.js")
				So(w.String(), should.Resemble, "Are you sure you want to delete hello.js? y/nhello.js deleted.")
				w.Reset()

				kwk.Run("get", "hello.js")
				So(w.String(), should.Resemble, "alias: hello.js not found\n")
				w.Reset()
			})

			Convey(`Should delete when partial key given which is not ambiguous`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()

				reader.WriteString("y\n")

				kwk.Run("rm", "hello")
				w.Reset()

				kwk.Run("get", "hello.js")
				So(w.String(), should.Resemble, "alias: hello.js not found\n")
				w.Reset()
			})


			Convey(`Should prompt to delete alias and not delete when not 'y'`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()

				reader.WriteString("b\n")

				kwk.Run("rm", "hello.js")
				So(w.String(), should.Resemble, "Are you sure you want to delete hello.js? y/nhello.js was pardoned.")
				w.Reset()
			})
		})

		Convey(`PATCH`, func() {
			Convey(`Should patch alias (find and replace)`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()

				kwk.Run("patch", "hello.js",  "echo", "printf")
				So(w.String(), should.Resemble, "hello.js patched.")
				w.Reset()
				kwk.Run("get", "hello.js")
				So(w.String(), should.Resemble, "printf \"hello\"")
				w.Reset()
			})

		})

		Convey(`TAG`, func() {
			Convey(`Should tag an alias`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				kwk.Run("tag", "hello.js", "tag1")
				So(w.String(), should.Resemble, "hello.js tagged.")
				w.Reset()
			})
			Convey(`Should show error if no tag given`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				kwk.Run("tag", "hello.js")
				So(w.String(), should.Resemble, "Please provide at least one tag.\n")
				w.Reset()
			})
			Convey(`Should untag an alias`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				kwk.Run("tag", "hello.js", "tag1")
				So(w.String(), should.Resemble, "hello.js tagged.")
				w.Reset()
				kwk.Run("untag", "hello.js", "tag1")
				So(w.String(), should.Resemble, "hello.js untagged.")
				w.Reset()
			})
			Convey(`Should show error when untagging if tag does not exist`, func() {
				signup(reader, kwk)
				w.Reset()
				kwk.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				kwk.Run("tag", "hello.js", "tag1")
				So(w.String(), should.Resemble, "hello.js tagged.")
				w.Reset()
				kwk.Run("untag", "hello.js", "donkey")
				So(w.String(), should.Resemble, "None of these tag(s) apply to 'hello.js': donkey\n")
				w.Reset()
			})
		})
	})
}
