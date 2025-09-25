package responses

import (
	"encoding/json"
	"errors"
)

type ApiResponse interface {
	JsonString() ([]byte, error)
	Code() int
	ContentType() string
	IsDataEmpty() bool
}

type Response struct {
	code        int
	data        interface{}
	contentType string
}

func (r *Response) IsDataEmpty() bool {
	return r.data == nil
}

func (r *Response) ContentType() string {
	return r.contentType
}

func (r *Response) Code() int {
	return r.code
}

func (r *Response) Data() interface{} {
	return r.data
}

func (r *Response) JsonString() ([]byte, error) {
	JSONMarshalled, err := json.Marshal(r.data)
	if err != nil {
		return nil, errors.New("failed to serialize response to JSON")
	}
	return JSONMarshalled, nil
}

func createResponse(code int, data interface{}) ApiResponse {
	return &Response{
		code:        code,
		data:        data,
		contentType: "application/json; charset=utf-8",
	}
}

func Ok(data interface{}) ApiResponse {
	return createResponse(200, data)
}

func Created(data interface{}) ApiResponse {
	return createResponse(201, data)
}

func NoContent() ApiResponse {
	return createResponse(204, nil)
}

func BadRequest(message string) ApiResponse {
	return createResponse(400, map[string]string{"error": message})
}

func NotAuthorized(message string) ApiResponse {
	return createResponse(401, map[string]string{"error": message})
}

func NotFound(message string) ApiResponse {
	return createResponse(404, map[string]string{"error": message})
}

func Forbiden(message string) ApiResponse {
	return createResponse(403, map[string]string{"error": message})
}

func InternalError(message string) ApiResponse {
	return createResponse(500, map[string]string{"error": message})
}
