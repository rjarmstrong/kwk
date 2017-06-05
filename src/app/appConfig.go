package app

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/kwk-super-snippets/cli/src/out"
	"encoding/json"
	"log"
	"github.com/kwk-super-snippets/cli/src/cli"
)

func GetConfig() *cli.AppConfig {
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