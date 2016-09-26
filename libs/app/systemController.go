package app

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/users"
)

type SystemController struct {
	service system.ISystem
	users   users.IUsers
	gui.ITemplateWriter
}

func NewSystemController(s system.ISystem, u users.IUsers, w gui.ITemplateWriter) *SystemController {
	return &SystemController{service:s, users:u, ITemplateWriter: w}
}

func (c *SystemController) Upgrade() {
	if err := c.service.Upgrade(); err != nil {
		c.Render("error", err)
	} else {
		c.Render("system:upgraded", nil)
	}
}

func (c *SystemController) GetVersion() {
	if v, err := c.service.GetVersion(); err != nil {
		c.Render("error", err)
	} else {
		c.Render("system:version", map[string]string{
			"version": v,})
	}
}