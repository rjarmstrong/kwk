package dlg

type DialogMock struct {
	LastModalCalledWith   []interface{}
	CallHistory           []interface{}
	ReturnItem            *DialogResponse
	FieldCallHistory      []interface{}
	FieldResponse         *DialogResponse
	FieldResponseMap map[string]interface{}
	MultiChoiceCalledWith []interface{}
	MultiChoiceResponse   *DialogResponse
}

func (d *DialogMock) Modal(templateName string, data interface{}, autoYes bool) *DialogResponse {
	d.LastModalCalledWith = []interface{}{templateName, data, autoYes}
	d.CallHistory = append(d.CallHistory, d.LastModalCalledWith)
	return d.ReturnItem
}

func (d *DialogMock) FormField(templateName string, data interface{}, mask bool) *DialogResponse {
	d.FieldCallHistory = append(d.FieldCallHistory, []interface{}{templateName, data})
	if d.FieldResponseMap[templateName] != nil {
		return &DialogResponse{Value:d.FieldResponseMap[templateName], Ok:true}
	}
	return d.FieldResponse
}

func (d *DialogMock) MultiChoice(templateName string, header interface{}, options interface{}) *DialogResponse {
	d.MultiChoiceCalledWith = []interface{}{templateName, header, options}
	return d.MultiChoiceResponse
}
