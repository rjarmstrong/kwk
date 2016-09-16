package gui

type IWriter interface {
	PrintWithTemplate(templateName string, input interface{})
}