package gui


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
