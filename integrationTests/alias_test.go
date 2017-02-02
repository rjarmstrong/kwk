package integration

import (
	"bytes"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Alias(t *testing.T) {
	cleanup()
	w := &bytes.Buffer{}
	reader := &bytes.Buffer{}
	app := getApp(reader, w)
	signup(reader, app)
	w.Reset()

	Convey("ALIAS COMMANDS", t, func() {
		Convey(`NEW`, func() {
			Convey(`Given no extension should use txt by default`, func() {
				//signin(reader, app)
				app.Run("new", "echo \"hello\"", "hello")
				So(lastLine(w.String()), ShouldResemble, "hello.txt created.")
				w.Reset()
			})
			Convey(`Should create new url with an extension`, func() {
				app.Run("new", "echo \"hello\"", "hello.go")
				So(lastLine(w.String()), ShouldEqual, "hello.go created.")
				w.Reset()
			})
			Convey(`Should create new blank alias if no content or alias is given`, func() {
				app.Run("new")
				So(lastLine(w.String()), ShouldContainSubstring, ".txt created.")
				w.Reset()
			})
		})

		Convey(`UPDATE`, func() {
			Convey(`Should add a description`, func() {
				app.Run("new", "echo \"hello\"", "hello.go")
				w.Reset()
				description := "This is for saying hello."
				app.Run("describe", "hello.go", description)
				So(w.String(), ShouldResemble, "Description updated:\n\x1b[36mThis is for saying hello.\x1b[0m")
				w.Reset()
			})
		})

		Convey(`INSPECT`, func() {
			Convey(`Should inspect alias`, func() {
				app.Run("new", "echo \"hello\"", "hello.go")
				So(lastLine(w.String()), ShouldEqual, "hello.go created.")
				w.Reset()
				app.Run("describe", "hello.go", "Hi there!")
				w.Reset()
				app.Run("get", "hello.go")
				w.Reset()
				app.Run("inspect", "hello.go")
				So(w.String(), ShouldResemble, "\nsnippet: testuser/hello.go\nRuntime: golang\nURI: echo \"hello\"\nVersion: 1\nTags: \nWeb: \x1b[4mhttp://aus.kwk.co/testuser/hello.go\x1b[0m\nDescription: Hi there!\nRun count: 2\n\n")
				w.Reset()
			})
		})

		Convey(`CAT`, func() {
			Convey(`Should cat unambiguous alias`, func() {
				uri := "echo \"hello\""
				app.Run("new", uri, "hello.go")
				So(lastLine(w.String()), ShouldEqual, "hello.go created.")
				w.Reset()
				app.Run("cat", "hello")
				So(w.String(), ShouldEqual, uri)
				w.Reset()
			})
			Convey(`Should prompt cat ambiguous alias`, func() {
				uri := "echo \"hello\""
				uri2 := "console.log('hello');"
				app.Run("new", uri, "hello.go")
				app.Run("new", uri2, "hello.js")
				w.Reset()
				app.Run("cat", "hello")
				So(w.String(), ShouldResemble, "That snippet is ambiguous please run it again with the extension:\nhello.go\nhello.js\n")
				w.Reset()
			})
		})

		Convey(`RENAME`, func() {
			Convey(`Should rename an alias`, func() {
				app.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				app.Run("mv", "hello.js", "dong.js")
				So(w.String(), ShouldResemble, "hello.js renamed to dong.js")
				w.Reset()
			})
			Convey(`Should rename an alias and auto-add the extension if not given in the new key`, func() {
				app.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				app.Run("mv", "hello.js", "dong")
				So(w.String(), ShouldResemble, "hello.js renamed to dong.js")
				w.Reset()
			})
			Convey(`Should rename an alias and even if no extension is given for the original key`, func() {
				app.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				app.Run("mv", "hello", "dong")
				So(w.String(), ShouldResemble, "hello.js renamed to dong.js")
				w.Reset()
			})
		})

		Convey(`DELETE`, func() {
			Convey(`Should prompt to delete alias and delete when 'y'`, func() {
				app.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()

				reader.WriteString("y\n")

				app.Run("rm", "hello.js")
				So(w.String(), ShouldResemble, "Are you sure you want to delete hello.js? y/nhello.js deleted.")
				w.Reset()

				app.Run("get", "hello.js")
				So(w.String(), ShouldResemble, "snippet: hello.js not found\n")
				w.Reset()
			})

			Convey(`Should delete when partial key given which is not ambiguous`, func() {
				app.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()

				reader.WriteString("y\n")

				app.Run("rm", "hello")
				w.Reset()

				app.Run("get", "hello.js")
				So(w.String(), ShouldResemble, "snippet: hello.js not found\n")
				w.Reset()
			})

			Convey(`Should prompt to delete alias and not delete when not 'y'`, func() {
				app.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()

				reader.WriteString("b\n")

				app.Run("rm", "hello.js")
				So(w.String(), ShouldResemble, "Are you sure you want to delete hello.js? y/nhello.js was pardoned.")
				w.Reset()
			})
		})

		Convey(`PATCH`, func() {
			Convey(`Should patch alias (find and replace)`, func() {
				app.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()

				app.Run("patch", "hello.js", "echo", "printf")
				So(w.String(), ShouldResemble, "hello.js patched.")
				w.Reset()
				app.Run("get", "hello.js")
				So(w.String(), ShouldResemble, "printf \"hello\"")
				w.Reset()
			})

		})

		Convey(`TAG`, func() {
			Convey(`Should tag an alias`, func() {
				app.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				app.Run("tag", "hello.js", "tag1")
				So(w.String(), ShouldResemble, "hello.js tagged.")
				w.Reset()
			})
			Convey(`Should show error if no tag given`, func() {
				app.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				app.Run("tag", "hello.js")
				So(w.String(), ShouldResemble, "Please provide at least one tag.\n")
				w.Reset()
			})
			Convey(`Should untag an alias`, func() {
				app.Run("new", "echo \"hello\"", "hello.js")
				w.Reset()
				app.Run("tag", "hello.js", "tag1")
				So(w.String(), ShouldResemble, "hello.js tagged.")
				w.Reset()
				app.Run("untag", "hello.js", "tag1")
				So(w.String(), ShouldResemble, "hello.js untagged.")
				w.Reset()
			})
			//Convey(`Should show error when untagging if tag does not exist`, func() {
			//	signup(reader, kwk)
			//	w.Reset()
			//	app.Run("new", "echo \"hello\"", "hello.js")
			//	w.Reset()
			//	app.Run("tag", "hello.js", "tag1")
			//	So(w.String(), ShouldResemble, "hello.js tagged.")
			//	w.Reset()
			//	app.Run("untag", "hello.js", "donkey")
			//	So(w.String(), ShouldResemble, "None of these tag(s) apply to 'hello.js': donkey\n")
			//	w.Reset()
			//})
		})
	})
}
