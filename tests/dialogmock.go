package tests

type DialogMock struct {
	LastModalCalledWith   []interface{}
	CallHistory           []interface{}
	ReturnItem            *DialogResponse
	FieldCallHistory      []interface{}
	FieldResponse         *DialogResponse
	FieldResponseMap      map[string]interface{}
	MultiChoiceCalledWith []interface{}
	MultiChoiceResponse   *types.Snippet
}

func (d *DialogMock) Modal(templateName string, data interface{}, autoYes bool) *DialogResponse {
	d.LastModalCalledWith = []interface{}{templateName, data, autoYes}
	d.CallHistory = append(d.CallHistory, d.LastModalCalledWith)
	return d.ReturnItem
}

func (d *DialogMock) TemplateFormField(templateName string, data interface{}, mask bool) (*DialogResponse, error) {
	d.FieldCallHistory = append(d.FieldCallHistory, []interface{}{templateName, data})
	if d.FieldResponseMap[templateName] != nil {
		return &DialogResponse{Value: d.FieldResponseMap[templateName], Ok: true}, nil
	}
	return d.FieldResponse, nil
}

func (d *DialogMock) FormField(label string) (*DialogResponse, error) {
	panic("not impl")
}

func (d *DialogMock) MultiChoice(templateName string, header interface{}, list []*types.Snippet) *types.Snippet {
	d.MultiChoiceCalledWith = []interface{}{templateName, header, list}
	return d.MultiChoiceResponse
}
