package system

import (
	"fmt"
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func Test_System(t *testing.T) {
	Convey("System test", t, func() {
		s := New()

		Convey(`Should create a directory if not exists`, func() {
			dir := "test_dir"
			path, err := s.GetDirPath(dir)
			So(err, ShouldBeNil)
			So(path, should.Equal, fmt.Sprintf("/Users/%s/Library/Caches/kwk/%s", os.Getenv("USER"), dir))
			fi, err := os.Stat(path)
			So(fi.IsDir(), should.BeTrue)
			err = os.RemoveAll(path)
			So(err, ShouldBeNil)
		})

		Convey(`Should check file exists and not exists`, func() {
			path, err := s.WriteToFile(".", "testfile.js", "some text")
			So(err, should.BeNil)
			ok, err := s.Exists(path)
			So(ok, should.BeTrue)
			err = os.RemoveAll(path)
			So(err, should.BeNil)
			ok, err = s.Exists(path)
			So(ok, should.BeFalse)
			So(err, should.BeNil)
		})

		Convey(`Should write and read from file`, func() {
			dir := "test_dir"
			uri := "git status"
			path, err := s.GetDirPath(dir)
			fullKey := "test.bash"
			p, err := s.WriteToFile(dir, fullKey, uri)
			So(err, ShouldBeNil)
			So(p, should.Equal, fmt.Sprintf("/Users/%s/Library/Caches/kwk/%s/%s", os.Getenv("USER"), dir, fullKey))
			txt, err := s.ReadFromFile(dir, fullKey)
			So(txt, should.Equal, uri)
			err = os.RemoveAll(path)
			So(err, ShouldBeNil)
		})
	})
}
