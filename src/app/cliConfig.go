package app

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"encoding/json"
	"log"
)

type CLIConfig struct {
	CpuProfile  bool   `default:"false" json:"KWK_PROFILE"`
	Debug       bool   `default:"false" json:"KWK_DEBUG"`
	APIHost     string `default:"api.kwk.co:443" json:"KWK_APIHOST"`
	TestMode    bool   `default:"false" json:"KWK_TESTMODE"`
	SnippetPath string `default:"snippets" json:"KWK_CACHEPATH"` // Snippets cache.
	DocPath     string `default:"docs" json:"KWK_DOCSPATH"`      // Any other non-snippet files.
	UserDocName string `default:"user" json:"KWK_USERDOCNAME"`   // File where user credentials are stored.
}

func GetConfig() *CLIConfig {
	err := envconfig.Process("KWK", cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	out.DebugEnabled = cfg.Debug
	if cfg.TestMode {
		cfg.APIHost = "localhost:8000"
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	out.Debug("CONFIG: %s", string(b))
	return cfg
}