package integration

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/smartystreets/assertions/should"
	"bytes"
	_ "github.com/go-sql-driver/mysql"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"os"
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
				Convey(`Should show version not set error`, func() {
					v := os.Getenv(system.APP_VERSION)
					os.Setenv(system.APP_VERSION, "")
					kwk.Run("version")
					So(w.String(), should.Equal, system.APP_VERSION + " has not been set.\n")
					w.Reset()
					os.Setenv(system.APP_VERSION, v)
				})
			})
		})
	})
}
