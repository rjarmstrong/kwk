package dlg

import "bitbucket.com/sharingmachine/kwkcli/models"

type DialogMock struct {
	LastModalCalledWith   []interface{}
	CallHistory           []interface{}
	ReturnItem            *DialogResponse
	FieldCallHistory      []interface{}
	FieldResponse         *DialogResponse
	FieldResponseMap map[string]interface{}
	MultiChoiceCalledWith []interface{}
	MultiChoiceResponse   *models.Snippet
}

func (d *DialogMock) Modal(templateName string, data interface{}, autoYes bool) *DialogResponse {
	d.LastModalCalledWith = []interface{}{templateName, data, autoYes}
	d.CallHistory = append(d.CallHistory, d.LastModalCalledWith)
	return d.ReturnItem
}

func (d *DialogMock) TemplateFormField(templateName string, data interface{}, mask bool) *DialogResponse {
	d.FieldCallHistory = append(d.FieldCallHistory, []interface{}{templateName, data})
	if d.FieldResponseMap[templateName] != nil {
		return &DialogResponse{Value:d.FieldResponseMap[templateName], Ok:true}
	}
	return d.FieldResponse
}

func (d *DialogMock) FormField(label string) *DialogResponse {
  panic("not impl")
}

func (d *DialogMock) MultiChoice(templateName string, header interface{}, list []*models.Snippet) *models.Snippet {
	d.MultiChoiceCalledWith = []interface{}{templateName, header, list}
	return d.MultiChoiceResponse
}
