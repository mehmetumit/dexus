package app

import (
	"context"
	"net/url"

	"github.com/mehmetumit/dexus/internal/core/ports"
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

func (a *App) FindRedirect(ctx context.Context, u *url.URL) (*url.URL, error) {
	var dataTo string
	key, err := a.Cacher.GenKey(ctx, u.String())
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
		a.Logger.Debug("Cache miss:", string(cachedTo))
		redirectionTo, err := a.RedirectionRepo.Get(ctx, u.String())
		if err != nil {
			a.Logger.Error("internal redirection repo error:", err)
			return nil, err
		}
		dataTo = redirectionTo

	} else {
		dataTo = string(cachedTo)
		a.Logger.Debug("Cache hit", dataTo)
	}
	to, err := url.Parse(dataTo)
	if err != nil {
		a.Logger.Error("internal parsing error:", err)
		return nil, err
	}

	return to, nil
}
