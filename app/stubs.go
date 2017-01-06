package app

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/cmd"
	"bitbucket.com/sharingmachine/kwkcli/search"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/system"
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
)

func CreateAppStub() *KwkApp {
	s := &system.SystemMock{}
	t := &config.SettingsMock{}
	a := &snippets.ServiceMock{}
	o := &cmd.RunnerMock{}
	u := &account.ManagerMock{}
	w := &tmpl.WriterMock{}
	h := &search.TermMock{}
	d := &dlg.DialogMock{}
	api := &rpc.SysMock{}
	app := New(a, s, t, o, u, d, w, h, api)
	return app
}
