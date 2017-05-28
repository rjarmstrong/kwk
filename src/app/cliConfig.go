package app

type CLIConfig struct {
	CpuProfile  bool   `default:"false" json:"KWK_PROFILE"`
	Debug       bool   `default:"false" json:"KWK_DEBUG"`
	APIHost     string `default:"api.kwk.co:443" json:"KWK_APIHOST"`
	TestMode    bool   `default:"false" json:"KWK_TESTMODE"`
	SnippetPath string `default:"snippets" json:"KWK_CACHEPATH"` // Snippets cache
	DocPath     string `default:"docs" json:"KWK_DOCSPATH"` // Any other non-snippet files
}
