package server

import (
	"context"
	"encoding/json"
	"fmt"
	"goserve/middlewares"
	"goserve/responses"
	"net/http"
	"strings"
)

const (
	CONTENT_TYPE string = "application/json; charset=utf-8"
)

type jsonResponseWriter struct {
	http.ResponseWriter
	config     JSONSerializationConfig
	request    *http.Request
	written    bool
	statusCode int
}

func (jw *jsonResponseWriter) handleContextResponse() {
	ctx := jw.request.Context()

	if err, ok := ctx.Value(ResponseKeyMessage).(error); ok {
		jw.writeError(err)
		return
	}

	if data := ctx.Value(ResponseKeyData); data != nil {
		status := http.StatusOK

		response := jw.buildContextResponse(ctx, data)
		jw.statusCode = status
		jw.writeHeader(status)
		jw.writeJSON(response)

	}
}

func (jw *jsonResponseWriter) buildContextResponse(ctx context.Context, data interface{}) interface{} {
	// response := SuccessResponse{
	// 	Success:   true,
	// 	Timestamp: time.Now().Unix(),
	// 	Data:      data,
	// }

	// // Ajouter le message s'il existe
	// if msg, ok := ctx.Value(ResponseMessageKey).(string); ok {
	// 	response.Message = msg
	// }

	// // Ajouter les métadonnées s'elles existent
	// if meta, ok := ctx.Value(ResponseMetaKey).(*Meta); ok {
	// 	response.Meta = meta
	// }

	return data
}

func (jw *jsonResponseWriter) writeError(err error) {
	jw.statusCode = http.StatusInternalServerError
	jw.WriteHeader(jw.statusCode)

	errorData := jw.config.ErrorWrapper(err, jw.request)
	jw.writeJSON(errorData)
}

func (jw *jsonResponseWriter) writeJSON(data interface{}) (int, error) {
	jw.setJSONHeaders()

	var jsonData []byte
	var err error

	if jw.config.PrettyPrint {
		jsonData, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonData, err = json.Marshal(data)
	}

	if err != nil {
		errorResponse := responses.InternalError("Failed to serialize response")
		jsonData, _ = json.Marshal(errorResponse)
		jw.writeHeader(http.StatusInternalServerError)
	}

	return jw.ResponseWriter.Write(jsonData)
}

func (jw *jsonResponseWriter) writeHeader(code int) {
	jw.statusCode = code
	jw.ResponseWriter.WriteHeader(code)
}

func (jw *jsonResponseWriter) Write(data []byte) (int, error) {
	if !jw.written {

		jw.written = true

		if jw.isJSONContent() && jw.isJSONValid(data) {
			jw.setJSONHeaders()
			jw.ResponseWriter.Write(data)
		}

		var parsedData interface{}
		if err := json.Unmarshal(data, &parsedData); err == nil {
			wrappedData := jw.wrapResponse(parsedData, jw.statusCode)
			return jw.writeJSON(wrappedData)
		}

		textData := string(data)
		wrappedData := jw.wrapResponse(textData, jw.statusCode)
		return jw.writeJSON(wrappedData)
	}

	return jw.ResponseWriter.Write(data)
}

func (jw *jsonResponseWriter) wrapResponse(data interface{}, code int) interface{} {
	if code >= 400 {
		if jw.config.ErrorWrapper != nil {
			if err, ok := data.(error); ok {
				return jw.config.ErrorWrapper(err, jw.request)
			}
			// Convertir en erreur
			return jw.config.ErrorWrapper(fmt.Errorf("%v", data), jw.request)
		}
	}

	if jw.config.WrapSingleValues && jw.config.SuccessWrapper != nil {
		if _, ok := data.(responses.ApiResponse); !ok {
			return jw.config.SuccessWrapper(data, jw.request)
		}
	}

	return data
}

func (jw *jsonResponseWriter) setJSONHeaders() {
	jw.Header().Set("Content-Type", CONTENT_TYPE)
	jw.Header().Set("X-Content-Type-Options", "nosniff")
}

func (jw *jsonResponseWriter) isJSONValid(data []byte) bool {
	return json.Valid(data)
}

func (jw *jsonResponseWriter) isJSONContent() bool {
	contentType := jw.Header().Get("Content-Type")
	return strings.Contains(contentType, "application/json")
}

func JSONSerializationMiddleware() middlewares.MiddlewareFunc {
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
