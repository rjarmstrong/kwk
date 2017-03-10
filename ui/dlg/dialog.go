package dlg

import "bitbucket.com/sharingmachine/kwkcli/models"

// Dialogue creates CLI ui elements to perform common interactions.
type Dialog interface {

	// Modal creates a yes/no prompt with a given template and data.
	Modal(templateName string, data interface{}, autoYes bool) *DialogResponse

	// FormField renders a prompt with a templated label to take user input.
	TemplateFormField(templateName string, data interface{}, mask bool) *DialogResponse

	FormField(label string) *DialogResponse

	// MultiChoice is a special modal with multiple possible choices.
	MultiChoice(templateName string, header interface{}, list []*models.Snippet) *DialogResponse
}