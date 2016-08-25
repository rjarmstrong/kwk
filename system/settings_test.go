package system

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/assertions/should"
	"testing"
)

func Test_UserService(t *testing.T) {
	Convey("Manage settings", t, func() {
		Convey(`Should create, update, delete and get a setting`, func() {
			s := NewSettings("test2.kwk.bolt.db")
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

			s.Delete(key)
			actual = &User{}
			err := s.Get(key, actual)
			So(err.Error(), should.Equal, "Key 'user' does not exist.")
		})
	})
}


