package app

import (
	"context"
	"net/url"
	"time"

	"github.com/mehmetumit/dexus/internal/core/ports"
)

var (
	defaultCacheTTL = 10 * time.Second
)

type AppConfig struct {
	Logger          ports.Logger
	RedirectionRepo ports.RedirectionRepo
	Cacher          ports.Cacher
}
type App struct {
	AppConfig
}

func NewApp(cfg AppConfig) *App {
	return &App{
		cfg,
	}
}

func (a *App) FindRedirect(ctx context.Context, p string) (*url.URL, error) {
	var dataTo string
	key, err := a.Cacher.GenKey(ctx, p)
	if err != nil {
		a.Logger.Error("internal key generation error:", err)
		return nil, err
	}
	cachedTo, err := a.Cacher.Get(ctx, key)
	if err != nil {
		if err != ports.ErrKeyNotFound {
			a.Logger.Error("internal cache error:", err)
			return nil, err
		}
		a.Logger.Debug("Cache miss:", p)
		redirectionTo, err := a.RedirectionRepo.Get(ctx, p)
		if err != nil {
			if err == ports.ErrRedirectionNotFound {
				a.Logger.Error("redirection not found on repo:", err)
			} else {
				a.Logger.Error("internal cache error:", err)
			}
			return nil, err
		}
		//Set cache after cache miss
		if err := a.Cacher.Set(ctx, key, redirectionTo, defaultCacheTTL); err != nil {
			a.Logger.Error("cache set error:", err)
		}
		dataTo = redirectionTo

	} else {
		dataTo = cachedTo
		a.Logger.Debug("Cache hit", dataTo)
	}
	to, err := url.Parse(dataTo)
	if err != nil {
		a.Logger.Error("internal parsing error:", err)
		return nil, err
	}

	return to, nil
}
