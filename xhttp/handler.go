package xhttp

import (
	"context"
	"errors"
	"github.com/VaynerAkaWalo/go-toolkit/xctx"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

const (
	StatusCode xctx.ContextKey = "status_code"
	Duration   xctx.ContextKey = "duration"
	Error      xctx.ContextKey = "error"
	Method     xctx.ContextKey = "method"
	Path       xctx.ContextKey = "path"
)

type errorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type statusCatcher struct {
	http.ResponseWriter
	statusCode int
}

type HttpHandler func(http.ResponseWriter, *http.Request) error

func (handler HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := context.WithValue(r.Context(), xctx.Transaction, uuid.New().String())
	ctx = context.WithValue(ctx, Method, r.Method)
	ctx = context.WithValue(ctx, Path, r.URL.Path)

	var code int
	catcher := &statusCatcher{ResponseWriter: w}

	err := handler(catcher, r.WithContext(ctx))
	if err != nil {
		var httpError *HttpError
		if errors.As(err, &httpError) {
			code = httpError.Code
		}

		_ = WriteResponse(w, code, errorResponse{Message: err.Error(), Code: code})
		ctx = context.WithValue(ctx, Error, err.Error())
	}

	if catcher.statusCode != 0 {
		code = catcher.statusCode
	}

	ctx = context.WithValue(ctx, Duration, time.Since(start).Milliseconds())
	ctx = context.WithValue(ctx, StatusCode, code)

	if code < 400 {
		slog.InfoContext(ctx, "request completed")
	} else {
		slog.ErrorContext(ctx, "request failed")
	}
}

func (s *statusCatcher) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}
