package gui

type IInteraction interface {
	Respond(templateName string, input interface{})
}