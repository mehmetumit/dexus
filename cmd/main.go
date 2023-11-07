package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mehmetumit/dexus/helper"
	"github.com/mehmetumit/dexus/internal/adapters/fileredirect"
	"github.com/mehmetumit/dexus/internal/adapters/memcache"
	"github.com/mehmetumit/dexus/internal/adapters/rest"
	"github.com/mehmetumit/dexus/internal/adapters/stdlog"
	"github.com/mehmetumit/dexus/internal/core/app"
)

var (
	Version string
	Commit  string
)

// Composition root
func main() {
	// Load with these defaults if envs are not set
	envCfg := helper.LoadConfigWithDefaults(helper.Config{
		DebugLevel:       true,
		Host:             "",
		Port:             "8080",
		ReadTimeoutMS:    1000 * time.Millisecond,
		WriteTimeoutMS:   1000 * time.Millisecond,
		CacheTTLSec:      60 * time.Second,
		YamlRedirectPath: "configs/redirection.yaml",
	})

	logger := stdlog.NewStdLog(envCfg.DebugLevel)
	logger.Info(fmt.Sprintf("Version: %s | Commit: %s", Version, Commit))

	cacher := memcache.NewMemCache(logger)

	redirectRepo, err := fileredirect.NewYamlRedirect(logger, envCfg.YamlRedirectPath)
	if err != nil {
		os.Exit(1)
	}
	appCore := app.NewApp(app.AppConfig{
		Logger:          logger,
		RedirectionRepo: redirectRepo,
	})
	rest.NewServer(logger, &rest.ServerConfig{
		Addr: envCfg.Host + ":" + envCfg.Port,
		//IdleTimeout:  90 * time.Second,
		ReadTimeout:  envCfg.ReadTimeoutMS,
		WriteTimeout: envCfg.WriteTimeoutMS,
	}, appCore, rest.WithCacher(cacher, logger, envCfg.CacheTTLSec))
}
