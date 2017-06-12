package _wip

type WriterMock struct {
	RenderCalledWith    []interface{}
	RenderCallCount     int
	RenderErrCalledWith error
}

func (t *WriterMock) HandleErr(error error) {
	t.RenderErrCalledWith = error
}

func (t *WriterMock) Render(templateName string, data interface{}) {
	t.RenderCallCount += 1
	t.RenderCalledWith = []interface{}{templateName, data}
}

func (t *WriterMock) Out(templateName string, data interface{}) {

}
