package system

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
	"fmt"
	"os"
)

func Test_UserService(t *testing.T) {
	Convey("Manage settings", t, func() {
		Convey(`Should create, update, delete and get a setting`, func() {
			s := NewSettings("test_leveldb")
			defer s.Close()
			key := "user"
			expected := User{
				Email:"richard@kwk.co",
				Token:"asdfsdfsdfuiu",
			}
			s.Upsert(key, expected)
			actual := &User{}
			s.Get(key, actual)
			So(*actual, should.Resemble, expected)

			expectedUpdate := User{
				Email:"richard@kwk.io",
				Token:"asdfsdfsdfuiu",
			}
			s.Upsert(key, expectedUpdate)
			s.Get(key, actual)
			So(*actual, should.Resemble, expectedUpdate)

			err := s.Delete(key)
			So(err, ShouldBeNil)

			actual = &User{}
			err = s.Get(key, actual)
			So(err.Error(), should.Equal, "Not found.")
		})

		Convey(`Should create a directory if not exists`, func() {
			dir := "test_dir"
			path, err := GetDirPath(dir)
			So(err, ShouldBeNil)
			So(path, should.Equal, fmt.Sprintf("/Users/%s/Library/Caches/kwk/%s", os.Getenv("USER"), dir))
			fi, err := os.Stat(path)
			So(fi.IsDir(), should.BeTrue)
			err = os.RemoveAll(path)
			So(err, ShouldBeNil)
		})


		Convey(`Should write and read from file`, func() {
			dir := "test_dir"
			uri := "git status"
			path, err := GetDirPath(dir)
			fullKey := "test.bash"
			p, err := WriteToFile(dir, fullKey, uri)
			So(err, ShouldBeNil)
			So(p, should.Equal, fmt.Sprintf("/Users/%s/Library/Caches/kwk/%s/%s", os.Getenv("USER"), dir, fullKey))

			txt, err := ReadFromFile(dir, fullKey)
			So(txt, should.Equal, uri)

			err = os.RemoveAll(path)
			So(err, ShouldBeNil)
		})
	})
}


