package app

import (
	"bitbucket.com/sharingmachine/kwkcli/src/snippets"
	"bitbucket.com/sharingmachine/kwkcli/src/cmd"
	"bitbucket.com/sharingmachine/kwkcli/src/user"
	"bitbucket.com/sharingmachine/kwkcli/src/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/src/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/src/rpc"
	"bitbucket.com/sharingmachine/kwkcli/src/persist"
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
