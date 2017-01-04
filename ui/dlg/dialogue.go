package dlg

// Dialogue creates CLI ui elements to perform common interactions.
type Dialogue interface {

	// Modal creates a yes/no prompt with a given template and data.
	Modal(templateName string, data interface{}) *DialogueResponse

	// FormField renders a prompt with a templated label to take user input.
	FormField(templateName string, data interface{}, mask bool) *DialogueResponse

	// Multichoice is a special modal with multiple possible choices.
	MultiChoice(templateName string, header interface{}, options interface{}) *DialogueResponse
}