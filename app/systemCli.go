package app

import (
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"fmt"
)

type SystemCli struct {
	service       sys.Manager
	accountManage account.Manager
	tmpl.Writer
	rpc rpc.Service
}

func NewSystemCli(s sys.Manager, r rpc.Service, u account.Manager, w tmpl.Writer) *SystemCli {
	return &SystemCli{service: s, accountManage: u, Writer: w, rpc: r}
}

func (c *SystemCli) Upgrade() {
	panic("not implemented")
}

func (c *SystemCli) GetVersion() {
	cliV, err := c.service.GetVersion()
	if err != nil {
		c.HandleErr(err)
	}
	apiV, err := c.rpc.GetApiInfo(); if err != nil {
		c.HandleErr(err)
	}
	c.Render("system:version", map[string]string{
		"cliVersion": cliV,
		"apiVersion": fmt.Sprintf("%s+%s", apiV.Version, apiV.Build),
	})
}

func (c *SystemCli) TestAppErr(multi bool) {
	c.HandleErr(c.rpc.TestAppError(multi))
}

func (c *SystemCli) TestTransErr() {
	c.HandleErr(c.rpc.TestTransportError())
}
