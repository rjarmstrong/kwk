package app

import (
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/types/vwrite"
)

type system struct {
	vwrite.Writer
	updater Updater
}

func NewSystem(w vwrite.Writer, u Updater) *system {
	return &system{Writer: w, updater: u}
}

func (c *system) Update() error {
	return c.updater.Run()
}

func (c *system) GetVersion() error {
	return c.EWrite(out.Version(cliInfo))
}
