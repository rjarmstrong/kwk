package dlg

// Dialogue creates CLI ui elements to perform common interactions.
type Dialog interface {

	// Modal creates a yes/no prompt with a given template and data.
	Modal(templateName string, data interface{}) *DialogResponse

	// FormField renders a prompt with a templated label to take user input.
	FormField(templateName string, data interface{}, mask bool) *DialogResponse

	// MultiChoice is a special modal with multiple possible choices.
	MultiChoice(templateName string, header interface{}, options interface{}) *DialogResponse
}