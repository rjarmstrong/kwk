package app

import (
	"bitbucket.com/sharingmachine/kwkcli/system"
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
)

type SystemCli struct {
	service       system.ISystem
	accountManage account.Manager
	tmpl.Writer
	rpc rpc.Sys
}

func NewSystemCli(s system.ISystem, r rpc.Sys, u account.Manager, w tmpl.Writer) *SystemCli {
	return &SystemCli{service: s, accountManage: u, Writer: w, rpc: r}
}

func (c *SystemCli) Upgrade() {
	if err := c.service.Upgrade(); err != nil {
		c.HandleErr(err)
	} else {
		c.Render("system:upgraded", nil)
	}
}

func (c *SystemCli) GetVersion() {
	if v, err := c.service.GetVersion(); err != nil {
		c.HandleErr(err)
	} else {
		c.Render("system:version", map[string]string{
			"version": v})
	}
}

func (c *SystemCli) TestAppErr(multi bool) {
	c.HandleErr(c.rpc.TestAppError(multi))
}

func (c *SystemCli) TestTransErr() {
	c.HandleErr(c.rpc.TestTransportError())
}
