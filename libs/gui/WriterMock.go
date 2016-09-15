package gui

type WriterMock struct {
	PrintCalledWith interface{}
}

func (w *WriterMock) Print(input interface{}){
	w.PrintCalledWith = input
}
