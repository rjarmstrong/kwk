package cmd

import
(
	"bitbucket.com/sharingmachine/types"
)

type RunnerMock struct {
	RunCalledWith  []interface{}
	EditCalledWith *types.Snippet
}

func (o *RunnerMock) Run(alias *types.Snippet, args []string) error {
	o.RunCalledWith = []interface{}{alias, args}
	return nil
}

func (o *RunnerMock) Edit(alias *types.Snippet) error {
	o.EditCalledWith = alias
	return nil
}
