package app

import (
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"github.com/kwk-super-snippets/cli/src/cli"
	"github.com/kwk-super-snippets/cli/src/out"
	"log"
)

func GetConfig() *cli.AppConfig {
	err := envconfig.Process("KWK", cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	out.DebugEnabled = cfg.Debug
	if cfg.TestMode && cfg.APIHost != cli.DefaultApiHost {
		cfg.APIHost = cli.DefaultTestApiHost
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	out.Debug("CONFIG: %s", string(b))
	return cfg
}
