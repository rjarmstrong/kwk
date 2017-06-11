package out

import (
	"bufio"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"github.com/siddontang/go/num"
	"io"
)

// Dialogue creates CLI ui elements to perform common interactions.
type Dialog interface {
	// Modal creates a yes/no prompt with a given template and data.
	Modal(handler vwrite.Handler, autoYes bool) *DialogResponse

	FormField(field vwrite.Handler, mask bool) (*DialogResponse, error)

	ChooseSnippet(s []*types.Snippet) *types.Snippet
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

// ChooseSnippet is similar to MultiChoice but returns immediately if there is only one snippet
func (d *StdDialog) ChooseSnippet(s []*types.Snippet) *types.Snippet {
	if len(s) == 1 {
		return s[0]
	} else if len(s) > 1 {
		s, _ := d.multiChoice(FreeText("Multiple matches. Choose a snippet:  "), s)
		return s
	}
	return nil
}

func (d *StdDialog) multiChoice(question vwrite.Handler, list []*types.Snippet) (*types.Snippet, error) {
	d.Write(question)
	for i, v := range list {
		if i%3 == 0 {
			d.Writer.Write(FreeText("\n"))
		}
		d.Writer.Write(FreeText(fmt.Sprintf("%d)  %s    ", i+1, v.Alias.URI())))
	}
	d.Writer.Write(FreeText("\n"))
	input, _, err := d.reader.ReadLine()
	if err != nil {
		return nil, err
	}
	i, err := num.ParseInt(string(input))
	if i > len(list) || err != nil {
		return d.multiChoice(FreeText("Please choose one snippet by number: "), list)
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
