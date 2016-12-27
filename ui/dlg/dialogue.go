package dlg

// Dialogue creates CLI ui elements to perform common interactions.
type Dialogue interface {
	Modal(templateName string, data interface{}) *DialogueResponse
	Field(templateName string, data interface{}) *DialogueResponse
	MultiChoice(templateName string, header interface{}, options interface{}) *DialogueResponse
}