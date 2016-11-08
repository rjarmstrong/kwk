package app

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/aliases"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/openers"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/users"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/search"
)

func CreateAppStub() *KwkApp {
	s := &system.SystemMock{}
	t := &settings.SettingsMock{}
	a := &aliases.AliasesMock{}
	o := &openers.OpenerMock{}
	u := &users.UsersMock{}
	d := &gui.DialogueMock{}
	w := &gui.TemplateWriterMock{}
	h := &search.SearchMock{}
	app := NewKwkApp(a, s, t, o, u, d, w, h)
	return app
}