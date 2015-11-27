package context

import (
	"net/http"
	"golang.org/x/net/context"
)


type ContextHandler interface {
	ServerHTTP(context.Context, http.ResponseWriter, *http.Request)
}
type handler struct {
	ctx     context.Context
	handler ContextHandler
}


type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (f ContextHandlerFunc) ContextHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	f(ctx, w, req)
}


func NewHandler(h ContextHandler) http.Handler {
	return NewHandlerWithContext(context.Background(), h)
}

func NewHandlerWithContext(ctx context.Context, h ContextHandler) http.Handler {
	return &handler{
		ctx:     ctx,
		handler: h,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.handler.ServeHTTP(h.ctx, w, req)
}

