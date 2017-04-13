package config

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"bitbucket.com/sharingmachine/kwkcli/file"
)

func Test_JsonFile(t *testing.T) {
	Convey("JsonFile test", t, func() {
		Convey(`Should create, update, delete and get a setting`, func() {
			sys := file.New()
			s := NewJsonSettings(sys, "test_settings")
			key := "user"
			expected := models.User{
				Email: "richard@kwk.co",
				Token: "asdfsdfsdfuiu",
			}
			s.Upsert(key, expected)
			actual := &models.User{}
			s.Get(key, actual)
			So(*actual, ShouldResemble, expected)

			expectedUpdate := models.User{
				Email: "richard@kwk.io",
				Token: "asdfsdfsdfuiu",
			}
			s.Upsert(key, expectedUpdate)
			s.Get(key, actual)
			So(*actual, ShouldResemble, expectedUpdate)

			err := s.Delete(key)
			So(err, ShouldBeNil)

			actual = &models.User{}
			err = s.Get(key, actual)
			So(err.Error(), ShouldEqual, "Not found.")
			p, _ := sys.GetDirPath("test_settings")
			os.RemoveAll(p)
		})
	})
}
