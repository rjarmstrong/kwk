package app

import (
	"github.com/kwk-links/kwk-cli/libs/services/system"
	"github.com/kwk-links/kwk-cli/libs/services/users"
	"github.com/kwk-links/kwk-cli/libs/services/gui"
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