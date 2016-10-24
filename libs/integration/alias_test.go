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
			Convey(`Should create new url`, func() {
				signup(reader, kwk)
				w.Reset()

				reader.WriteString("1\n")
				kwk.Run("new", "echo \"hello\"", "hello")
				So(lastLine(w.String()), should.Equal, "hello.go created.")
				w.Reset()
			})
		})
	})
}
