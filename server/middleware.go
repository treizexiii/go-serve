package server

import (
	"fmt"
	"net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler

type Middleware struct {
	Name       string
	Middleware MiddlewareFunc
	Path       []string
	Method     []string
}

func (m *Middleware) HasPath() bool {
	return len(m.Path) > 0
}

func (m *Middleware) HasMethod() bool {
	return len(m.Method) > 0
}

func (m *Middleware) Apply(path string, method string) bool {
	if m.HasPath() {
		pathMatch := false
		for _, p := range m.Path {
			if p == path || m.pathMatch(p, path) {
				pathMatch = true
				break
			}
		}

		if !pathMatch {
			return false
		}
	}

	if m.HasMethod() {
		methodMatch := false
		for _, m := range m.Method {
			if m == method {
				methodMatch = true
				break
			}
		}
		if !methodMatch {
			return false
		}
	}

	return true
}

func (m *Middleware) pathMatch(pattern, path string) bool {
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(path) >= len(prefix) && path[:len(prefix)] == prefix
	}

	return pattern == path
}

func (m *Middleware) GetScope() string {
	scope := "global"
	if m.HasPath() {
		scope = fmt.Sprintf("path: %v, methods: %v", m.Path, m.Method)
	}

	return scope
}
