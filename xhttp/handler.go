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

type (
	httpHandler struct {
		withErrorHandler      func(http.ResponseWriter, *http.Request) error
		authenticationHandler AuthenticationHandler
	}

	errorResponse struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	statusCatcher struct {
		http.ResponseWriter
		statusCode int
	}
)

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ctx := context.WithValue(r.Context(), xctx.Transaction, uuid.New().String())
	ctx = context.WithValue(ctx, Method, r.Method)
	ctx = context.WithValue(ctx, Path, r.URL.Path)

	ctx, httpError := h.authenticationHandler.Authorize(ctx, r)
	if httpError != nil {
		_ = WriteResponse(w, httpError.Code, errorResponse{Message: httpError.Error(), Code: httpError.Code})
		ctx = context.WithValue(ctx, Error, httpError.Error())
		h.logRequestCompletion(ctx, httpError.Code, start)
		return
	}

	var code int
	catcher := &statusCatcher{ResponseWriter: w}

	err := h.withErrorHandler(catcher, r.WithContext(ctx))
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

	h.logRequestCompletion(ctx, code, start)
}

func (h httpHandler) logRequestCompletion(ctx context.Context, code int, start time.Time) {
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
