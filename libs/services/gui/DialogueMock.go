package gui

type DialogueMock struct {
	LastModalCalledWith   []interface{}
	CallHistory           []interface{}
	ReturnItem            *DialogueResponse
	FieldCallHistory      []interface{}
	FieldResponse         *DialogueResponse
	MultiChoiceCalledWith []interface{}
	MultiChoiceResponse   *DialogueResponse
}

func (d *DialogueMock) Modal(templateName string, data interface{}) *DialogueResponse {
	d.LastModalCalledWith = []interface{}{templateName, data}
	d.CallHistory = append(d.CallHistory, d.LastModalCalledWith)
	return d.ReturnItem
}

func (d *DialogueMock) Field(templateName string, data interface{}) *DialogueResponse {
	d.FieldCallHistory = append(d.FieldCallHistory, []interface{}{templateName, data})
	return d.FieldResponse
}

func (d *DialogueMock) MultiChoice(templateName string, header interface{}, options ...interface{}) *DialogueResponse {
	d.MultiChoiceCalledWith = []interface{}{templateName, header, options}
	return d.MultiChoiceResponse
}