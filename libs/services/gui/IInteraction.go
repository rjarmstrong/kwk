package gui

type IInteraction interface {
	Respond(templateName string) Response
}

type Response func(input interface{}, error error) interface{}