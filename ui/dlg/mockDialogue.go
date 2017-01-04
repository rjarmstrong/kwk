package dlg

type DialogueMock struct {
	LastModalCalledWith   []interface{}
	CallHistory           []interface{}
	ReturnItem            *DialogueResponse
	FieldCallHistory      []interface{}
	FieldResponse         *DialogueResponse
	FieldResponseMap map[string]interface{}
	MultiChoiceCalledWith []interface{}
	MultiChoiceResponse   *DialogueResponse
}

func (d *DialogueMock) Modal(templateName string, data interface{}) *DialogueResponse {
	d.LastModalCalledWith = []interface{}{templateName, data}
	d.CallHistory = append(d.CallHistory, d.LastModalCalledWith)
	return d.ReturnItem
}

func (d *DialogueMock) FormField(templateName string, data interface{}, mask bool) *DialogueResponse {
	d.FieldCallHistory = append(d.FieldCallHistory, []interface{}{templateName, data})
	if d.FieldResponseMap[templateName] != nil {
		return &DialogueResponse{Value:d.FieldResponseMap[templateName], Ok:true}
	}
	return d.FieldResponse
}

func (d *DialogueMock) MultiChoice(templateName string, header interface{}, options interface{}) *DialogueResponse {
	d.MultiChoiceCalledWith = []interface{}{templateName, header, options}
	return d.MultiChoiceResponse
}
