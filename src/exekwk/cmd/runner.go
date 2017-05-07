package cmd

import (
	"github.com/kwk-super-snippets/types"
)

type Runner interface {
	Run(s *types.Snippet, args []string) error
	Edit(s *types.Snippet) error
}
