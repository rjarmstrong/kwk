package gui

type InteractionMock struct {
	LastRespondCalledWith []interface{}
	CallHistory [][]interface{}
	ReturnItem            interface{}
}

func (i *InteractionMock) Respond(templateName string, input interface{}) interface{} {
	i.LastRespondCalledWith = []interface{}{templateName, input}
	i.CallHistory = append(i.CallHistory, i.LastRespondCalledWith)
	return i.ReturnItem
}