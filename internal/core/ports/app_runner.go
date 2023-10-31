package ports

import (
	"context"
	"net/url"
)

type AppRunner interface {
	FindRedirect(ctx context.Context, u *url.URL) (*url.URL, error)
}
