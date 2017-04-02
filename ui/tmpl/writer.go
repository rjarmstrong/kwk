package tmpl

type Writer interface {
	Render(templateName string, data interface{})
 	Out(templateName string, data interface{})
/*
 HandleErr may require a specific error implementation as a parameter. Please check implementation for details.
 */
	HandleErr(error error)
}