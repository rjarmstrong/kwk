package gui

type InteractionMock struct {
	RespondCalledWith []interface{}
}

func (w *InteractionMock) Respond(templateName string, input interface{}){
	w.RespondCalledWith = []interface{}{templateName, input}
}