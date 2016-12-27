package dlg

type MockDialogue struct {
	LastModalCalledWith   []interface{}
	CallHistory           []interface{}
	ReturnItem            *DialogueResponse
	FieldCallHistory      []interface{}
	FieldResponse         *DialogueResponse
	MultiChoiceCalledWith []interface{}
	MultiChoiceResponse   *DialogueResponse
}

func (d *MockDialogue) Modal(templateName string, data interface{}) *DialogueResponse {
	d.LastModalCalledWith = []interface{}{templateName, data}
	d.CallHistory = append(d.CallHistory, d.LastModalCalledWith)
	return d.ReturnItem
}

func (d *MockDialogue) Field(templateName string, data interface{}) *DialogueResponse {
	d.FieldCallHistory = append(d.FieldCallHistory, []interface{}{templateName, data})
	return d.FieldResponse
}

func (d *MockDialogue) MultiChoice(templateName string, header interface{}, options interface{}) *DialogueResponse {
	d.MultiChoiceCalledWith = []interface{}{templateName, header, options}
	return d.MultiChoiceResponse
}
