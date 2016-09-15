package gui

import "fmt"

type IWriter interface {
	Print(input interface{})
}

type Writer struct {
}

func (w *Writer) Print(input interface{}){
	fmt.Println(input)
}

type WriterMock struct {
	PrintCalledWith interface{}
}

func (w *WriterMock) Print(input interface{}){
	w.PrintCalledWith = input
}