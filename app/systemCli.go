package app

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/update"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/persist"
)

type SystemCli struct {
	tmpl.Writer
	rpc     rpc.Service
	updater *update.Runner
}

func NewSystemCli(s persist.IO, r rpc.Service, w tmpl.Writer, p persist.Persister) *SystemCli {
	return &SystemCli{rpc: r, Writer: w, updater:update.NewRunner(p)}
}

func (c *SystemCli) Update() {
	err := c.updater.Run()
	if err != nil {
		c.HandleErr(err)
	}
}

func (c *SystemCli) GetVersion() {
	c.Render("system:version", &models.Client)
}

func (c *SystemCli) TestAppErr(multi bool) {
	c.HandleErr(c.rpc.TestAppError(multi))
}

func (c *SystemCli) TestTransErr() {
	c.HandleErr(c.rpc.TestTransportError())
}
