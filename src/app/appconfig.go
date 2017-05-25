package app

type CLIConfig struct {
	Profile  bool   `default:"false" json:"KWK_PROFILE"`
	Debug    bool   `default:"false" json:"KWK_DEBUG"`
	APIHost  string `default:"api.kwk.co:443" json:"KWK_APIHOST"`
	TestMode bool   `default:"false" json:"KWK_TESTMODE"`
}
