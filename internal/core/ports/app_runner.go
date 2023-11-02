package ports

import (
	"context"
	"net/url"
)

type AppRunner interface {
	FindRedirect(ctx context.Context, p string) (*url.URL, error)
}
