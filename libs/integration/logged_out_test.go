package integration

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"bytes"
)

func Test_Logged_out(t *testing.T) {
	cleanup()
	w := &bytes.Buffer{}
	reader := &bytes.Buffer{}
	app := getApp(reader, w)

	Convey(`Not logged in when creating new`, t, func() {
		w.Reset()
		app.Run("new", "http://somelink.com")
		So(w.String(), should.Equal, notLoggedIn)
		w.Reset()
	})

}