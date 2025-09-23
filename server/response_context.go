package server

import (
	"context"
	"net/http"
)

type contextKey string

const (
	ContextKeyRequestId contextKey = "requestId"
	ResponseKeyData     contextKey = "data"
	ResponseKeyMessage  contextKey = "message"
	ResponseKeyStatus   contextKey = "status"
)

type ResponseContext struct {
	ctx context.Context
}

func NewResponseContext(ctx context.Context) *ResponseContext {
	return &ResponseContext{ctx: ctx}
}

func (rc *ResponseContext) WithData(data interface{}) *ResponseContext {
	rc.ctx = context.WithValue(rc.ctx, ResponseKeyData, data)
	return rc
}

func (rc *ResponseContext) WithMessage(message string) *ResponseContext {
	rc.ctx = context.WithValue(rc.ctx, ResponseKeyMessage, message)
	return rc
}

func (rc *ResponseContext) WithStatus(status int) *ResponseContext {
	rc.ctx = context.WithValue(rc.ctx, ResponseKeyStatus, status)
	return rc
}

func (rc *ResponseContext) Context() context.Context {
	return rc.ctx
}

type JSONHandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func ConvertJSONHandler(handler JSONHandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := handler(w, r)

		if err != nil {
			ctx := context.WithValue(r.Context(), ResponseKeyMessage, err.Error())
			_ = r.WithContext(ctx)
			return
		}

		if data != nil {
			ctx := context.WithValue(r.Context(), ResponseKeyData, data)
			_ = r.WithContext(ctx)
		}
	}

}
