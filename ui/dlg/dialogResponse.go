package dlg

// Dialogue Response carries the users input back to the calling code.
type DialogResponse struct {
	Ok    bool
	Value interface{}
}
