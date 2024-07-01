package view

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/anthdm/superkit/kit"
	"github.com/anthdm/superkit/kit/middleware"
)

// Asset is a view helper that returns the full asset path as a
// string based on the given asset name.
//
//	view.Asset("styles.css") // => /public/assets/styles.css.
func Asset(name string) string {
	return fmt.Sprintf("/public/assets/%s", name)
}

// getContextValue is a helper function to retrieve a value from the context.
// It returns the value if present, otherwise returns the provided default value.
func getContextValue[T any](ctx context.Context, key interface{}, defaultValue T) T {
	value, ok := ctx.Value(key).(T)
	if !ok {
		return defaultValue
	}
	return value
}

// Auth is a view helper function that returns the current Auth.
// If Auth is not set, a default Auth will be returned.
//
//	view.Auth(ctx)
func Auth(ctx context.Context) kit.Auth {
	return getContextValue(ctx, kit.AuthKey{}, kit.DefaultAuth{})
}

// URL is a view helper that returns the current URL.
// The request path can be accessed with:
//
//	view.URL(ctx).Path // => ex. /login
func URL(ctx context.Context) *url.URL {
	return getContextValue(ctx, middleware.RequestKey{}, &http.Request{}).URL
}

// Request is a view helper that returns the current http request.
// The request can be accessed with:
//
//	view.Request(ctx)
func Request(ctx context.Context) *http.Request {
	return getContextValue(ctx, middleware.RequestKey{}, &http.Request{})
}
