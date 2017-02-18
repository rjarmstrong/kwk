package app

import (
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/update"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"fmt"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"bitbucket.com/sharingmachine/kwkcli/config"
)

type SystemCli struct {
	service       sys.Manager
	accountManage account.Manager
	tmpl.Writer
	rpc rpc.Service
	updater *update.Runner
}

func NewSystemCli(s sys.Manager, r rpc.Service, u account.Manager, w tmpl.Writer, p config.Persister) *SystemCli {
	return &SystemCli{service: s, accountManage: u, Writer: w, rpc: r, updater:update.NewRunner(p)}
}

func (c *SystemCli) Update() {
	err := c.updater.Run()
	if err != nil {
		c.HandleErr(err)
	}
}

func (c *SystemCli) GetVersion(cliVersion string) {
	apiV, err := c.rpc.GetApiInfo(); if err != nil {
		c.HandleErr(err)
	}
	c.Render("system:version", map[string]string{
		"cliVersion": cliVersion,
		"apiVersion": fmt.Sprintf("%s+%s", apiV.Version, apiV.Build),
	})
}

func (c *SystemCli) TestAppErr(multi bool) {
	c.HandleErr(c.rpc.TestAppError(multi))
}

func (c *SystemCli) TestTransErr() {
	c.HandleErr(c.rpc.TestTransportError())
}
