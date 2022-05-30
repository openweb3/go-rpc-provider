package rpc

import "net/http"

type HandlerFuncMiddleware func(next http.HandlerFunc) http.HandlerFunc

func (s *Server) HookServeHttp(middleware HandlerFuncMiddleware) {
	s.serveHTTPMiddlewares = append(s.serveHTTPMiddlewares, middleware)
}

func (s *Server) HookOnCreateWebsocketHandler(middleware HandlerFuncMiddleware) {
	s.serveWebsocktMiddlewares = append(s.serveHTTPMiddlewares, middleware)
}

func (s *Server) genHttpHandlerNestedWare(core http.HandlerFunc) http.HandlerFunc {
	nested := core
	for i := len(s.serveHTTPMiddlewares) - 1; i >= 0; i-- {
		nested = s.serveHTTPMiddlewares[i](nested)
	}
	return nested
}

func (s *Server) genWebsocketHandlerNestedWare(core http.HandlerFunc) http.HandlerFunc {
	nested := core
	for i := len(s.serveWebsocktMiddlewares) - 1; i >= 0; i-- {
		nested = s.serveWebsocktMiddlewares[i](nested)
	}
	return nested
}
