package app

import (
	"context"
	"net/url"

	"github.com/mehmetumit/dexus/internal/core/ports"
)

type AppConfig struct {
	Logger          ports.Logger
	RedirectionRepo ports.RedirectionRepo
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
	redirectionTo, err := a.RedirectionRepo.Get(ctx, p)
	if err != nil {
		if err == ports.ErrRedirectionNotFound {
			a.Logger.Error("redirection not found on repo:", err)
		} else {
			a.Logger.Error("internal redirection repo error:", err)
		}
		return nil, err
	}

	to, err := url.Parse(redirectionTo)
	if err != nil {
		a.Logger.Error("internal parsing error:", err)
		return nil, err
	}

	return to, nil
}
