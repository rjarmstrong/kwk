package dlg

import (
	. "github.com/smartystreets/goconvey/convey"
	_ "bitbucket.com/sharingmachine/kwkcli/src/ui/tmpl"
	"testing"
	_ "bytes"
	_ "bufio"
)

func Test_StdDialog(t *testing.T) {
	Convey("Dialog test", t, func() {
		//w := &tmpl.WriterMock{}
		//b := &bytes.Buffer{}
		//reader := bufio.NewReader(b)
		//d := New(w, reader)

		//Convey(`Given a numeric input should return dialog option`, func() {
		//	options := []string{"apple", "banana", "pear"}
		//	b.WriteString("2\n")
		//	r := d.MultiChoice("dialog:choose", "some header", options)
		//	So(w.RenderCalledWith[0], ShouldResemble, "dialog:choose")
		//	So(w.RenderCallCount, ShouldEqual, 2)
		//	So(r.Value, ShouldResemble, "banana")
		//})
		//
		//Convey(`Given a non-numeric input should re-ask for options`, func() {
		//	options := []string{"apple", "banana", "pear"}
		//	b.WriteString("sdfs\n")
		//	b.WriteString("1\n")
		//	r := d.MultiChoice("dialog:choose", "some header", options)
		//	So(w.RenderCalledWith[0], ShouldResemble, "dialog:choose")
		//	So(w.RenderCallCount, ShouldEqual, 4)
		//	So(r.Value, ShouldResemble, "apple")
		//})
		//
		//Convey(`Given an out of range input should re-ask for options`, func() {
		//	options := []string{"apple", "banana", "pear"}
		//	b.WriteString("4\n")
		//	b.WriteString("2\n")
		//	r := d.MultiChoice("dialog:choose", "some header", options)
		//	So(w.RenderCalledWith[0], ShouldResemble, "dialog:choose")
		//	So(w.RenderCallCount, ShouldEqual, 4)
		//	So(r.Value, ShouldResemble, "banana")
		//})

	})
}
