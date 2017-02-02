package integration

import (
	. "github.com/smartystreets/goconvey/convey"
	"bytes"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func Test_System(t *testing.T) {
	Convey("SYSTEM COMMANDS", t, func() {
		w := &bytes.Buffer{}
		reader := &bytes.Buffer{}
		kwk := getApp(reader, w)

		Convey(`SYSTEM`, func() {
			Convey(`VERSION`, func() {
				Convey(`Should get version`, func() {
					kwk.Run("version")
					So(w.String(), ShouldContainSubstring, "kwk ")
					w.Reset()
				})
			})
		})
	})
}
