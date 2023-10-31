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
		//Internal err
		return nil, err
	}
	cachedTo, err := a.Cacher.Get(ctx, key)
	if err != nil {
		if err != ports.ErrKeyNotFound {
			//Internal err
			return nil, err
		}
		//Cache miss
		redirectionTo, err := a.RedirectionRepo.Get(ctx, u.String())
		if err != nil {
			//internal err
			return nil, err
		}
		dataTo = redirectionTo

	} else {
		//Cache hit
		dataTo = string(cachedTo)
	}
	to, err := url.Parse(dataTo)

	return to, nil
}
