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
	mockLogger.SetDebugLevel(true)
	return NewApp(
		AppConfig{
			Logger:          mockLogger,
			RedirectionRepo: mockRedirectionRepo,
		},
	)
}

func TestApp_FindRedirect(t *testing.T) {
	redirectionMap := mocks.MockRedirectionMap{
		"test1": "https://test1.com",
		"test2": "https://test2.com",
		"test3": "https://test3.com",
	}
	notFoundKeys := []string{
		"not-found",
		"not/found",
		"",
	}
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

}
