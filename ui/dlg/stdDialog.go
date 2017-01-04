package dlg

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"github.com/siddontang/go/num"
	"github.com/howeyc/gopass"
	"reflect"
	"bufio"
)

func New(w tmpl.Writer, reader *bufio.Reader) *StdDialog {
	return &StdDialog{writer: w, reader: reader}
}

// StdDialogue is the default dialogue type.
type StdDialog struct {
	writer tmpl.Writer
	reader *bufio.Reader
}

func (d *StdDialog) Modal(templateName string, data interface{}) *DialogResponse {
	d.writer.Render(templateName, data)
	yesNo, _, _ := d.reader.ReadRune()
	return &DialogResponse{
		Ok: string(yesNo) == "y",
	}
}

func (d *StdDialog) MultiChoice(templateName string, header interface{}, options interface{}) *DialogResponse {
	d.writer.Render("dialog:header", header)
	o := InterfaceSlice(options)
	d.writer.Render(templateName, options)
	value, _, _ := d.reader.ReadLine()
	if i, err := num.ParseInt(string(value)); err != nil {
		panic(err)
	} else {
		if i > len(o) {
			d.MultiChoice(templateName, header, options)
		}
		return &DialogResponse{
			Value: o[i-1],
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

func (d *StdDialog) FormField(templateName string, data interface{}, mask bool) *DialogResponse {
	d.writer.Render(templateName, data)
	var value []byte
	var err error
	if mask {
		value, err = gopass.GetPasswdMasked()
	} else {
		value, _, err = d.reader.ReadLine()
	}
	if err != nil {
		panic(err.Error())
	}
	//d.reader.Reset(nil)
	return &DialogResponse{
		Value: string(value),
	}
}
