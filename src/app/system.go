package app

import (
	"bitbucket.com/sharingmachine/kwkcli/src/app/out"
	"bitbucket.com/sharingmachine/kwkcli/src/update"
	"bitbucket.com/sharingmachine/types/vwrite"
)

type system struct {
	vwrite.Writer
	updater *update.Runner
}

func NewSystemCli(w vwrite.Writer, u *update.Runner) *system {
	return &system{Writer: w, updater: u}
}

func (c *system) Update() error {
	return c.updater.Run()
}

func (c *system) GetVersion() error {
	return c.EWrite(out.Version(CLIInfo))
}
