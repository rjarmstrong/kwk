package cmd

import (
	"bitbucket.com/sharingmachine/types"
)

type Runner interface {
	Run(s *types.Snippet, args []string) error
	Edit(s *types.Snippet) error
}
