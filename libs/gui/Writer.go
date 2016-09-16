package gui


func NewInteraction(templates map[string]Template) *Interaction {
	return &Interaction{templates:templates}
}

type Interaction struct {
	templates map[string]Template
}

type Template func(input interface{})

var templates = map[string]Template{}

func (w *Interaction) Respond(templateName string, input interface{}) {
	templates[templateName](input)
}
