package gui

import (
	"bufio"
	"os"
	"github.com/siddontang/go/num"
	"fmt"
	"reflect"
)

func NewDialogues(w ITemplateWriter) *Dialogues {
	return &Dialogues{writer:w}
}

type Dialogues struct {
	writer ITemplateWriter
}

func (d *Dialogues) Modal(templateName string, data interface{}) *DialogueResponse {
	reader := bufio.NewReader(os.Stdin)
	d.writer.Render(templateName, data)
	yesNo, _, _ := reader.ReadRune()
	return &DialogueResponse{
		Ok:string(yesNo) == "y",
	}
}

func (d *Dialogues) MultiChoice(templateName string, header interface{}, options interface{}) *DialogueResponse {
	reader := bufio.NewReader(os.Stdin)
	// TODO: Render header and options
	items := InterfaceSlice(options)
	fmt.Println(header)
	d.writer.Render(templateName, items)
	fmt.Println()
	value, _, _ := reader.ReadRune()
	// upper and lower contraints
	if i, err := num.ParseInt(string(value)); err != nil {
		panic(err)
	} else {
		return &DialogueResponse{
			Value: items[i],
		}
	}
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}

func (d *Dialogues) Field(templateName string, data interface{}) *DialogueResponse {
	reader := bufio.NewReader(os.Stdin)
	d.writer.Render(templateName, data)
	value, _, _ := reader.ReadLine()
	reader.Reset(nil)
	return &DialogueResponse{
		Value:string(value),
	}
}

type DialogueResponse struct {
	Ok    bool
	Value interface{}
}