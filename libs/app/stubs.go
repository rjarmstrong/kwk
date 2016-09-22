package app

import (
	"github.com/kwk-links/kwk-cli/libs/services/settings"
	"github.com/kwk-links/kwk-cli/libs/services/system"
	"github.com/kwk-links/kwk-cli/libs/services/aliases"
	"github.com/kwk-links/kwk-cli/libs/services/openers"
	"github.com/kwk-links/kwk-cli/libs/services/users"
	"github.com/kwk-links/kwk-cli/libs/services/gui"
)

func CreateAppStub() *KwkApp {
	s := &system.SystemMock{}
	t := &settings.SettingsMock{}
	a := &aliases.AliasesMock{}
	o := &openers.OpenerMock{}
	u := &users.UsersMock{}
	d := &gui.DialogueMock{}
	w := &gui.TemplateWriterMock{}
	app := NewKwkApp(a, s, t, o, u, d, w)
	return app
}
