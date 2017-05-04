package gokwk

import (
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/metadata"
	"testing"
)

func Test_RPC(t *testing.T) {
	Convey("RPC", t, func() {

		Convey(`Upgrade`, func() {
			Convey(`Given the current settings has a profile (signed in user) should add token to context`, func() {
				t := &config.PersisterMock{}
				token := "sometoken234234234"
				t.GetHydrates = &models.User{Token: token}
				h := NewHeaders(t)
				md, _ := metadata.FromContext(h.Context())
				So(md[models.TokenHeaderName][0], ShouldEqual, token)
			})

			Convey(`Given the current settings does not have a profile (signed in user) should not add token to context`, func() {
				t := &config.PersisterMock{}
				h := NewHeaders(t)
				md, _ := metadata.FromContext(h.Context())
				So(md[models.TokenHeaderName][0], ShouldBeEmpty)
			})
		})

	})
}
