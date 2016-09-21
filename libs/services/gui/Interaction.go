package gui

import "fmt"

func NewInteraction(templates map[string]Template) *Interaction {
	return &Interaction{templates:templates}
}

type Interaction struct {
	templates map[string]Template
}

type Template func(input interface{}) interface{}

var templates = map[string]Template{}

func (w *Interaction) Respond(templateName string, input interface{}) interface{} {
	return templates[templateName](input)
}


// Respond2 can handle both interface and error responses
// e.g. one output or one form
func (w *Interaction) Respond2(templateName string) Response {
	return Response(func(input interface{}, error error) {
		if error != nil {
			fmt.Println(error)
			return nil
		}
		return templates[templateName](input)
	})
}