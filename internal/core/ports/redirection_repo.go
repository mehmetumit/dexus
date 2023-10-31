package ports

import (
	"context"
	"errors"
)

var (
	ErrRedirectionNotFound = errors.New("redirection not found")
)

type RedirectionRepo interface {
	Get(ctx context.Context, from string) (string, error)
}
