package gui

import (
	"bufio"
	"os"
	"github.com/siddontang/go/num"
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

func (d *Dialogues) MultiChoice(templateName string, header interface{}, options ...interface{}) *DialogueResponse {
	reader := bufio.NewReader(os.Stdin)
	// TODO: Render header and options
	//d.writer.Render(templateName, data)
	value, _, _ := reader.ReadRune()
	// upper and lower contraints
	if i, err := num.ParseInt(string(value)); err != nil {
		panic(err)
	} else {
		return &DialogueResponse{
			Value: options[i],
		}
	}
}

func (d *Dialogues) Field(templateName string, data interface{}) *DialogueResponse {
	reader := bufio.NewReader(os.Stdin)
	d.writer.Render(templateName, data)
	value, _, _ := reader.ReadLine()
	return &DialogueResponse{
		Value:string(value),
	}
}

type DialogueResponse struct {
     Ok bool
     Value interface{}
}