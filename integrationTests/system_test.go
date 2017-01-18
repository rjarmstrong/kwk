package integration

import (
	"bitbucket.com/sharingmachine/kwkcli/system"
	"bytes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"os"
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
					So(w.String(), should.ContainSubstring, "kwk ")
					w.Reset()
				})
			})
		})
	})
}
