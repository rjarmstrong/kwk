package app

import (
	"bufio"
	"github.com/howeyc/gopass"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/vwrite"
	"github.com/siddontang/go/num"
	"io"
)

// Dialogue creates CLI ui elements to perform common interactions.
type Dialog interface {
	// Modal creates a yes/no prompt with a given template and data.
	Modal(handler vwrite.Handler, autoYes bool) *DialogResponse

	FormField(field vwrite.Handler, mask bool) (*DialogResponse, error)

	// MultiChoice is a special modal with multiple possible choices.
	MultiChoice(vwrite.Handler, []*types.Snippet) (*types.Snippet, error)
}

// Dialogue Response carries the users input back to the calling code.
type DialogResponse struct {
	Ok    bool
	Value interface{}
}

func NewDialog(w vwrite.Writer, r io.Reader) *StdDialog {
	return &StdDialog{Writer: w, reader: bufio.NewReader(r)}
}

// StdDialogue is the default dialogue type.
type StdDialog struct {
	vwrite.Writer
	reader *bufio.Reader
}

func (d *StdDialog) Modal(handler vwrite.Handler, autoYes bool) *DialogResponse {
	r := &DialogResponse{}
	if autoYes {
		r.Ok = true
		return r
	}
	d.Writer.Write(handler)
	yesNo, _, _ := d.reader.ReadRune()
	r.Ok = string(yesNo) == "y" || string(yesNo) == "Y"
	return r
}

func (d *StdDialog) MultiChoice(question vwrite.Handler, list []*types.Snippet) (*types.Snippet, error) {
	d.Write(question)
	input, _, err := d.reader.ReadRune()
	if err != nil {
		return nil, err
	}
	i, err := num.ParseInt(string(input))
	if i > len(list) || err != nil {
		return d.MultiChoice(out.FreeText("Please choose one snippet by number: "), list)
	}
	return list[i-1], nil
}

func (d *StdDialog) FormField(field vwrite.Handler, mask bool) (*DialogResponse, error) {
	d.Write(field)
	var value []byte
	var err error
	if mask {
		value, err = gopass.GetPasswdMasked()
	} else {
		value, _, err = d.reader.ReadLine()
	}
	if err != nil {
		return nil, err
	}
	return &DialogResponse{
		Ok:    true,
		Value: string(value),
	}, nil
}
