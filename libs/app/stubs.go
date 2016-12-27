package app

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/services/aliases"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/openers"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/search"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/users"
	"bitbucket.com/sharingmachine/kwkcli/libs/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/libs/ui/tmpl"
)

func CreateAppStub() *KwkApp {
	s := &system.SystemMock{}
	t := &settings.SettingsMock{}
	a := &aliases.AliasesMock{}
	o := &openers.OpenerMock{}
	u := &users.UsersMock{}
	d := &dlg.MockDialogue{}
	w := &tmpl.MockWriter{}
	h := &search.SearchMock{}
	app := NewKwkApp(a, s, t, o, u, d, w, h)
	return app
}
