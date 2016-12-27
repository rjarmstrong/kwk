package dlg

import (
	"bufio"
	"fmt"
	"github.com/siddontang/go/num"
	"reflect"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
)

func New(w tmpl.Writer, reader *bufio.Reader) *StdDialogue {
	return &StdDialogue{writer: w, reader: reader}
}

// StdDialogue is the default dialogue type.
type StdDialogue struct {
	writer tmpl.Writer
	reader *bufio.Reader
}

func (d *StdDialogue) Modal(templateName string, data interface{}) *DialogueResponse {
	d.writer.Render(templateName, data)
	yesNo, _, _ := d.reader.ReadRune()
	return &DialogueResponse{
		Ok: string(yesNo) == "y",
	}
}

func (d *StdDialogue) MultiChoice(templateName string, header interface{}, options interface{}) *DialogueResponse {
	items := InterfaceSlice(options)
	fmt.Println(header)
	d.writer.Render(templateName, items)
	fmt.Println()
	value, _, _ := d.reader.ReadLine()
	// upper and lower contraints
	if i, err := num.ParseInt(string(value)); err != nil {
		panic(err)
	} else {
		if i > len(items) {
			d.MultiChoice(templateName, header, options)
		}
		return &DialogueResponse{
			Value: items[i-1],
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

func (d *StdDialogue) Field(templateName string, data interface{}) *DialogueResponse {
	d.writer.Render(templateName, data)
	value, _, err := d.reader.ReadLine()
	if err != nil {
		panic(err.Error())
	}
	//d.reader.Reset(nil)
	return &DialogueResponse{
		Value: string(value),
	}
}
