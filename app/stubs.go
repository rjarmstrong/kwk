package app

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/cmd"
	"bitbucket.com/sharingmachine/kwkcli/user"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/persist"
)

func CreateAppStub() *KwkApp {
	f := &persist.IoMock{}
	t := &persist.PersisterMock{}
	a := &snippets.ServiceMock{}
	o := &cmd.RunnerMock{}
	u := &user.AccountMock{}
	w := &tmpl.WriterMock{}
	d := &dlg.DialogMock{}
	api := &rpc.SysMock{}
	return NewApp(a, f, t, o, u, d, w, api)
}
