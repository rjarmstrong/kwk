package app

import (
	"bitbucket.com/sharingmachine/kwkcli/src/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/src/update"
)

type system struct {
	tmpl.Writer
	updater *update.Runner
}

func NewSystemCli(w tmpl.Writer, u *update.Runner) *system {
	return &system{ Writer: w, updater: u}
}

func (c *system) Update() {
	err := c.updater.Run()
	if err != nil {
		c.HandleErr(err)
	}
}

func (c *system) GetVersion() {
	c.Render("system:version", CLIInfo)
}