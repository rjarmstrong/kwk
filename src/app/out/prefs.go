package out

var prefs Prefs

type Prefs struct {
	ListAll          bool //List all items including private. TODO: implement on api in search SEARCH.
	Naked            bool
	ListHorizontal   bool
	SlimRows         int
	ExpandedRows     int
	AlwaysExpandRows bool
	RowSpaces        bool
	RowLines         bool
	DisablePreview   bool
}
