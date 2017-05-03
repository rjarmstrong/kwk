package app

import (
	"bitbucket.com/sharingmachine/kwkcli/src/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/src/update"
)

type SystemCli struct {
	tmpl.Writer
	updater *update.Runner
}

func NewSystemCli(w tmpl.Writer, u *update.Runner) *SystemCli {
	return &SystemCli{ Writer: w, updater:u}
}

func (c *SystemCli) Update() {
	err := c.updater.Run()
	if err != nil {
		c.HandleErr(err)
	}
}

func (c *SystemCli) GetVersion() {
	c.Render("system:version", CLIInfo)
}