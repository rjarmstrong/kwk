package dlg

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"github.com/siddontang/go/num"
	"github.com/howeyc/gopass"
	"reflect"
	"bufio"
	"bitbucket.com/sharingmachine/kwkcli/models"
)

func New(w tmpl.Writer, reader *bufio.Reader) *StdDialog {
	return &StdDialog{writer: w, reader: reader}
}

// StdDialogue is the default dialogue type.
type StdDialog struct {
	writer tmpl.Writer
	reader *bufio.Reader
}

func (d *StdDialog) Modal(templateName string, data interface{}, autoYes bool) *DialogResponse {
	r := &DialogResponse{}
	if autoYes {
		r.Ok = true
		return r
	}
	d.writer.Render(templateName, data)
	yesNo, _, _ := d.reader.ReadRune()
	r.Ok = string(yesNo) == "y" || string(yesNo) == "Y"
	return r
}

func (d *StdDialog) MultiChoice(templateName string, header interface{}, list []*models.Snippet) *DialogResponse {
	d.writer.Render("dialog:header", header)
	d.writer.Render(templateName, list)
	input, _, _ := d.reader.ReadLine()
	i, err := num.ParseInt(string(input))
	if i > len(list) || err != nil {
		return d.MultiChoice(templateName, "Please choose a number.", list)
	}
	return &DialogResponse{
		Value: list[i-1],
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

func (d *StdDialog) FormField(label string) *DialogResponse {
	d.writer.Render("free-text", label)
	var value []byte
	var err error
	value, _, err = d.reader.ReadLine()
	if err != nil {
		panic(err.Error())
	}
	//d.reader.Reset(nil)
	return &DialogResponse{
		Ok:true,
		Value: string(value),
	}
}

func (d *StdDialog) TemplateFormField(templateName string, data interface{}, mask bool) *DialogResponse {
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
