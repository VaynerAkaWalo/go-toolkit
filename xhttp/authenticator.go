package xhttp

import (
	"context"
	"github.com/VaynerAkaWalo/go-toolkit/xctx"
	"net/http"
)

const (
	UserId        xctx.ContextKey = "user_id"
	Token         string          = "X-AUTH-TOKEN"
	AuthSchema    string          = "X-AUTH-SCHEMA"
	SessionCookie string          = "session_id"
	SessionV1     string          = "SessionV1"
)

type User struct {
	UserId string
}

type (
	Authenticator struct {
		provider     AuthenticationProvider
		strategies   []authenticationStrategy
		excludePaths map[string]bool
	}

	AuthenticationProvider interface {
		FetchUser(ctx context.Context, token string, schema string) (User, error)
	}

	authenticationStrategy interface {
		resolveTokenAndSchema(r *http.Request) (bool, string, string)
	}

	sessionInCookieAuthenticationStrategy struct{}
	tokenAuthenticationStrategy           struct{}
)

func NewAuthenticator(provider AuthenticationProvider, excludePaths ...string) Authenticator {
	excludePathsMap := make(map[string]bool)
	for _, path := range excludePaths {
		excludePathsMap[path] = true
	}

	return Authenticator{
		provider:     provider,
		strategies:   []authenticationStrategy{sessionInCookieAuthenticationStrategy{}, tokenAuthenticationStrategy{}},
		excludePaths: excludePathsMap,
	}
}

func (authN Authenticator) authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	currentPath := r.Method + r.URL.Path
	isExcluded, _ := authN.excludePaths[currentPath]

	if isExcluded {
		return ctx, nil
	}

	var found bool
	var token, schema string
	for _, strategy := range authN.strategies {
		found, token, schema = strategy.resolveTokenAndSchema(r)

		if found {
			break
		}
	}

	if !found {
		return ctx, NewError("unable to find authentication credentials", http.StatusUnauthorized)
	}

	user, err := authN.provider.FetchUser(ctx, token, schema)
	if err != nil {
		return ctx, NewError("unable to authenticate with given credentials", http.StatusUnauthorized)
	}

	ctx = context.WithValue(ctx, UserId, user.UserId)
	return ctx, nil
}

func (sessionStrategy sessionInCookieAuthenticationStrategy) resolveTokenAndSchema(r *http.Request) (bool, string, string) {
	sessionCookie, err := r.Cookie(SessionCookie)
	if err != nil {
		return false, "", ""
	}

	return true, sessionCookie.Value, SessionV1
}

func (t tokenAuthenticationStrategy) resolveTokenAndSchema(r *http.Request) (bool, string, string) {
	token := r.Header.Get(Token)
	if token == "" {
		return false, token, ""
	}

	schema := r.Header.Get(AuthSchema)
	if schema == "" {
		return false, token, schema
	}

	return true, token, schema
}
