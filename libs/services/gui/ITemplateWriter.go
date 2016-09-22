package gui

type ITemplateWriter interface {
	Render(templateName string, data interface{})
}