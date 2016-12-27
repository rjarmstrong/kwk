package tmpl

type Writer interface {
	Render(templateName string, data interface{})
}
