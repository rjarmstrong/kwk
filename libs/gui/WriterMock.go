package gui

type WriterMock struct {
	PrintCalledWith []interface{}
}

func (w *WriterMock) PrintWithTemplate(templateName string, input interface{}){
	w.PrintCalledWith = []interface{}{templateName, input}
}
