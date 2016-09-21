package controllers

import (
	"github.com/kwk-links/kwk-cli/libs/services/system"
)

type SystemController struct {
	service system.ISystem
}

func NewSystemController(s system.ISystem) *SystemController {
	return &SystemController{service:s}
}

func (c *SystemController) Upgrade(){

}

func (c *SystemController) GetVersion(){

}

func (c *SystemController) ChangeDirectory(username string){

}