package app

import (
	"bitbucket.com/sharingmachine/kwkcli/system"
	"bitbucket.com/sharingmachine/kwkcli/users"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
)

type SystemCli struct {
	service system.ISystem
	users   users.IUsers
	tmpl.Writer
}

func NewSystemCli(s system.ISystem, u users.IUsers, w tmpl.Writer) *SystemCli {
	return &SystemCli{service: s, users: u, Writer: w}
}

func (c *SystemCli) Upgrade() {
	if err := c.service.Upgrade(); err != nil {
		c.Render("error", err)
	} else {
		c.Render("system:upgraded", nil)
	}
}

func (c *SystemCli) GetVersion() {
	if v, err := c.service.GetVersion(); err != nil {
		c.Render("error", err)
	} else {
		c.Render("system:version", map[string]string{
			"version": v})
	}
}
