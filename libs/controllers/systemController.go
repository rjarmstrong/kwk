package controllers

import (
	"github.com/kwk-links/kwk-cli/libs/services/system"
	"github.com/kwk-links/kwk-cli/libs/services/users"
)

type SystemController struct {
	service system.ISystem
	users users.IUsers
}

func NewSystemController(s system.ISystem, u users.IUsers) *SystemController {
	return &SystemController{service:s, users:u}
}

func (c *SystemController) Upgrade(){
   c.service.Upgrade()
}

func (c *SystemController) GetVersion(){
   c.service.GetVersion()
}

func (c *SystemController) ChangeDirectory(username string){

}