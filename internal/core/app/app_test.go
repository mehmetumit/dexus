package app

import (
	"context"
	"testing"

	"github.com/mehmetumit/dexus/internal/core/ports"
	"github.com/mehmetumit/dexus/internal/mocks"
)

func newTestApp(t testing.TB) *App {
	t.Helper()
	mockLogger := mocks.NewMockLogger()
	mockRedirectionRepo := mocks.NewMockRedirectionRepo()
	mockCacher := mocks.NewMockCacher()
	mockLogger.SetDebugLevel(true)
	return NewApp(
		AppConfig{
			Logger:          mockLogger,
			RedirectionRepo: mockRedirectionRepo,
			Cacher:          mockCacher,
		},
	)
}

func TestApp_FindRedirect(t *testing.T) {
	redirectionMap := mocks.MockRedirectionMap{
		"test1": "https://test1.com",
		"test2": "https://test2.com",
		"test3": "https://test3.com",
	}
	notFoundKeys := []string{"not-found", "not/found", ""}
	t.Run("Find Redirect From Repo", func(t *testing.T) {
		app := newTestApp(t)
		ctx := context.Background()
		app.RedirectionRepo.(*mocks.MockRedirectionRepo).SetMockRedirectionMap(redirectionMap)
		for k, v := range redirectionMap {
			u, err := app.FindRedirect(ctx, k)
			if err != nil {
				t.Errorf("Expected err nil, got %v", err)
			}

			if u.String() != v {
				t.Errorf("Expected redirection %v, got %v", v, u.String())

			}
		}
	})
	t.Run("Find Not Found Redirect From Repo", func(t *testing.T) {
		app := newTestApp(t)
		app.RedirectionRepo.(*mocks.MockRedirectionRepo).SetMockRedirectionMap(redirectionMap)
		for _, v := range notFoundKeys {
			_, err := app.FindRedirect(context.Background(), v)
			if err != ports.ErrRedirectionNotFound {
				t.Errorf("Expected err %v, got %v", ports.ErrRedirectionNotFound, err)
			}
		}
	})
	t.Run("Find Redirect From Cache", func(t *testing.T) {
		app := newTestApp(t)
		app.RedirectionRepo.(*mocks.MockRedirectionRepo).SetMockRedirectionMap(redirectionMap)
		ctx := context.Background()
		for k := range redirectionMap {
			app.FindRedirect(ctx, k)
		}
		for k, v := range redirectionMap {
			keyHash, err := app.Cacher.GenKey(ctx, k)
			if err != nil {
				t.Errorf("Expected err nil, got %v", err)
			}
			gotVal, err := app.Cacher.Get(ctx, keyHash)
			if err != nil {
				t.Errorf("Expected err nil, got %v", err)
			}
			if gotVal != v {
				t.Errorf("Expected cache val %v, got %v", v, gotVal)
			}
			_, err = app.FindRedirect(ctx, k)
			if err != nil {
				t.Errorf("Expected err nil, got %v", err)
			}
		}

	})

}
