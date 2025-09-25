package middlewares

import "net/http"

type MiddlewareFunc func(http.Handler) http.Handler

type MiddlewareInfo struct {
	Name       string
	Middleware MiddlewareFunc
}
