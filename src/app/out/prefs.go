package out

var prefs Prefs

// Prefs to control display formatting.
type Prefs struct {
	// Always opens browser in covert mode, when set to true.
	Covert bool

	// If true then `kwk <snipname>` will NOT execute the snippet without the `run|r` parameter.
	// In this case `view|v` command will be required to view the details of a snippet`
	RequireRunKeyword bool

	// When running snippets this is the duration a snippet should run before timing out.
	// Set to 0 for infinite.
	CommandTimeout int64

	AutoYes bool // TASK: Security

	ListAll           bool // List all snippets including private.
	GlobalSearch      bool // Search all accounts (not just yours).
	Naked             bool // Disable ansi styling, useful for scripting kwk.
	SnippetThumbRows  int  // Number of lines of code to show in lists.
	ExpandedThumbRows int  // Number of lines of code to show when list is expanded.
	ExpandedRows      bool // When true use value from ExpandedThumbRows
	RowSpaces         bool // Add blank lines between snippets in lists.
	RowLines          bool // Hard lines (hard-rules) between snippets in lists.
	DisablePreview    bool // Do not save a preview of the output of running snippets.
	Quiet             bool // Only display snippet URIs when listing, useful for scripting.
	ListHorizontal    bool // When true list snippets horizontally
}
