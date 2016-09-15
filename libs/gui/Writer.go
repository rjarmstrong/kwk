package gui

import "fmt"

type Writer struct {
}

func (w *Writer) Print(input interface{}){
	fmt.Println(input)
}
