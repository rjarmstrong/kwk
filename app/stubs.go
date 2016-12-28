package app

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/openers"
	"bitbucket.com/sharingmachine/kwkcli/search"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/system"
	"bitbucket.com/sharingmachine/kwkcli/users"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
)

func CreateAppStub() *KwkApp {
	s := &system.MockSystem{}
	t := &config.SettingsMock{}
	a := &snippets.Mock{}
	o := &openers.OpenerMock{}
	u := &users.UsersMock{}
	d := &dlg.MockDialogue{}
	w := &tmpl.MockWriter{}
	h := &search.TermMock{}
	app := New(a, s, t, o, u, d, w, h)
	return app
}
