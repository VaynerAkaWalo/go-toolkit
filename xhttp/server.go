package xhttp

import (
	"log/slog"
	"net/http"
)

type (
	Server struct {
		Addr         string
		Handlers     []RouteHandler
		AuthNHandler AuthenticationHandler
	}

	Router struct {
		*http.ServeMux
		authenticationHandler AuthenticationHandler
	}

	RouteHandler interface {
		RegisterRoutes(router *Router)
	}
)

func (server *Server) ListenAndServe() error {
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    server.Addr,
		Handler: mux,
	}

	apiRouter := Router{mux, server.AuthNHandler}
	for _, routeHandler := range server.Handlers {
		routeHandler.RegisterRoutes(&apiRouter)
	}

	slog.Info("Starting http server at port " + server.Addr)
	return httpServer.ListenAndServe()
}

func (router *Router) RegisterHandler(route string, routeHandler func(http.ResponseWriter, *http.Request) error) {
	router.Handle(route, httpHandler{routeHandler, NoOpAuthenticationHandler{}})
}

func (router *Router) RegisterHandlerWithAuth(route string, routeHandler func(http.ResponseWriter, *http.Request) error) {
	router.Handle(route, httpHandler{routeHandler, router.authenticationHandler})
}
