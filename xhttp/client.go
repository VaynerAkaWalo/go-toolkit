package xhttp

import (
	"context"
	"github.com/VaynerAkaWalo/go-toolkit/xctx"
	"io"
	"net/http"
)

const (
	TxHeader string = "X-TX-ID"
)

func NewRequest(ctx context.Context, method string, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequestWithContext(ctx, method, url, body)

	txId, ok := ctx.Value(xctx.Transaction).(string)
	if ok && txId != "" {
		req.Header.Set(TxHeader, txId)
	}

	return req
}
