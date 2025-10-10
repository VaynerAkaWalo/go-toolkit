package xhttp

import (
	"context"
	"github.com/VaynerAkaWalo/go-toolkit/xctx"
	"net/http"
)

const (
	UserId xctx.ContextKey = "user_id"
)

type (
	AuthenticationHandler interface {
		Authorize(context.Context, *http.Request) (context.Context, *HttpError)
	}

	NoOpAuthenticationHandler struct{}
)

func (handler NoOpAuthenticationHandler) Authorize(ctx context.Context, r *http.Request) (context.Context, *HttpError) {
	return ctx, nil
}
