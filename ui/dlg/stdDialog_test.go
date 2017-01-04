package dlg

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"testing"
	"bytes"
	"bufio"
)

func Test_StdDialog(t *testing.T) {
	Convey("Dialog test", t, func() {
		w := &tmpl.WriterMock{}
		b := &bytes.Buffer{}
		reader := bufio.NewReader(b)
		d := New(w, reader)

		Convey(`Given a numeric input should return dialog option`, func() {
			options := []string{"apple", "banana", "pear"}
			b.WriteString("2\n")
			r := d.MultiChoice("dialog:choose", "some header", options)
			So(w.RenderCalledWith[0], should.Resemble, "dialog:choose")
			So(w.RenderCallCount, should.Equal, 2)
			So(r.Value, should.Resemble, "banana")
		})

		Convey(`Given a non-numeric input should re-ask for options`, func() {
			options := []string{"apple", "banana", "pear"}
			b.WriteString("sdfs\n")
			b.WriteString("1\n")
			r := d.MultiChoice("dialog:choose", "some header", options)
			So(w.RenderCalledWith[0], should.Resemble, "dialog:choose")
			So(w.RenderCallCount, should.Equal, 4)
			So(r.Value, should.Resemble, "apple")
		})

		Convey(`Given an out of range input should re-ask for options`, func() {
			options := []string{"apple", "banana", "pear"}
			b.WriteString("4\n")
			b.WriteString("2\n")
			r := d.MultiChoice("dialog:choose", "some header", options)
			So(w.RenderCalledWith[0], should.Resemble, "dialog:choose")
			So(w.RenderCallCount, should.Equal, 4)
			So(r.Value, should.Resemble, "banana")
		})

	})
}