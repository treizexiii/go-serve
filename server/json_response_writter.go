package server

import "net/http"

type jsonResponseWriter struct {
	http.ResponseWriter
	// config     JSONSerializationConfig
	request    *http.Request
	written    bool
	statusCode int
}

func (jw *jsonResponseWriter) handleContextResponse() {
	ctx := jw.request.Context()

	if err, ok := ctx.Value(ResponseKeyMessage).(string); ok {
		jw.writeError(err)
		return
	}

	if data := ctx.Value(ResponseKeyData); data != nil {
		status := http.StatusOK

		response := jw.buildContextResponse(data, status)
		jw.statusCode = status
		jw.writeHeader(status)
		jw.writeJson(response)

	}
}

func (jw *jsonResponseWriter) buildContextResponse(data any, status int) any {
	panic("unimplemented")
}

func (jw *jsonResponseWriter) writeError(err string) {
	panic("unimplemented")
}

func (jw *jsonResponseWriter) writeJson(response any) {
	panic("unimplemented")
}

func (jw *jsonResponseWriter) writeHeader(code int) {
	jw.statusCode = code
	jw.ResponseWriter.WriteHeader(code)
}

func JSONSerializationMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jsonWriter := &jsonResponseWriter{
				ResponseWriter: w,
				request:        r,
			}

			next.ServeHTTP(jsonWriter, r)

			if !jsonWriter.written {
				jsonWriter.handleContextResponse()
			}
		})
	}
}
