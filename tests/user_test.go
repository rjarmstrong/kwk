package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_UserRoutes(t *testing.T) {
	app := NewApp()
	//d := app.Dialogue.(*dlg.DialogMock)

	// Should signin
	p := app.App.Command("signin")
	app.App.Run([]string{"[app]", "signin"})

	// Should Get Profile
	p := app.App.Command("me")
	assert.Equal(t, "me", p.Name)
}

//func Test_App(t *testing.T) {
//	Convey("ACCOUNT COMMANDS", t, func() {
//		Convey(`Signin`, func() {
//			Convey(`Should run by name`, func() {
//				p := app.App.Command("signin")
//				So(p, ShouldNotBeNil)
//				p2 := app.App.Command("login")
//				So(p2.Name, ShouldEqual, p.Name)
//			})
//			Convey(`Should call api signin`, func() {
//				app.App.Run([]string{"[app]", "signin", "richard", "password"})
//				So(u.LoginCalledWith[0], ShouldEqual, "richard")
//				So(u.LoginCalledWith[1], ShouldEqual, "password")
//			})
//			Convey(`Should call api signin and enter form details`, func() {
//				d.FieldResponse = &dlg.DialogResponse{Value: "richard"}
//				app.App.Run([]string{"[app]", "signin"})
//			})
//			Convey(`Should save details to file when signed in`, func() {
//				d.FieldResponse = &dlg.DialogResponse{Value: "richard"}
//				u.SignInResponse = &models.User{
//					Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dnZWRJbkFzIjoiYWRtaW4iLCJpYXQiOjE0MjI3Nzk2Mzh9.gzSraSYS8EXBxLN_oWnFSRgCzcmJmMjLiuyu5CSpyHI",
//				}
//				app.App.Run([]string{"[app]", "signin"})
//				So(t.UpsertCalledWith, ShouldResemble, []interface{}{profileFileName, u.SignInResponse})
//			})
//		})
//
//		Convey(`Signup`, func() {
//			Convey(`Should run by name`, func() {
//				p := app.App.Command("signup")
//				So(p, ShouldNotBeNil)
//				p2 := app.App.Command("register")
//				So(p2.Name, ShouldEqual, p.Name)
//			})
//			 Convey(`Should call api signup`, func() {
//				d.FieldResponseMap = map[string]interface{}{}
//				d.FieldResponseMap["account:signup:email"] = "richard@kwk.co"
//				d.FieldResponseMap["account:signup:username"] = "richard"
//				d.FieldResponseMap["account:signup:password"] = "password"
//
//				app.App.Run([]string{"[app]", "signup"})
//				So(u.SignupCalledWith[0], ShouldEqual, "richard@kwk.co")
//				So(u.SignupCalledWith[1], ShouldEqual, "richard")
//				So(u.SignupCalledWith[2], ShouldEqual, "password")
//				 d.FieldResponseMap = nil
//			})
//		})
//
//		Convey(`Signout`, func() {
//			Convey(`Should run by name`, func() {
//				p := app.App.Command("signout")
//				So(p, ShouldNotBeNil)
//				p2 := app.App.Command("logout")
//				So(p2.Name, ShouldEqual, p.Name)
//			})
//			Convey(`Should call api signout`, func() {
//				app.App.Run([]string{"[app]", "signout"})
//				So(u.SignoutCalled, ShouldBeTrue)
//				So(t.DeleteCalledWith, ShouldResemble, profileFileName)
//			})
//		})
//	})
//}
