package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	//TODO env
	debugLevel := true
	logger := stdlog.NewStdLog(debugLevel)
	logger.Info(fmt.Sprintf("Version: %s | Commit: %s", Version, Commit))
	cacher := memcache.NewMemCache(logger)
	fPath := filepath.FromSlash("test_data/redirection.yaml")
	redirectRepo, err := fileredirect.NewYamlRedirect(logger, fPath)
	if err != nil {
		os.Exit(1)
	}
	appCore := app.NewApp(app.AppConfig{
		Logger:          logger,
		RedirectionRepo: redirectRepo,
	})
	rest.NewServer(logger, &rest.ServerConfig{
		Addr: "localhost:8080",
		//IdleTimeout:  90 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}, appCore, rest.WithCacher(cacher, logger, 5*time.Second))
}
