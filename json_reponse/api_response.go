package jsonreponse

import (
	"encoding/json"
	"time"
)

type ApiResponse interface {
	JsonString() string
}

type ErrorResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message"`
	Timestamp    int64  `json:"timestamp"`
}

type SuccessResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

func (s *SuccessResponse) JsonString() string {
	JSONMarshalled, _ := json.Marshal(s)
	return string(JSONMarshalled)
}

func (e *ErrorResponse) JsonString() string {
	JSONMarshalled, _ := json.Marshal(e)
	return string(JSONMarshalled)
}

func Ok(data interface{}) ApiResponse {
	return &SuccessResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func Ko(message string) ApiResponse {
	return &ErrorResponse{
		Success:      false,
		ErrorMessage: message,
		Timestamp:    time.Now().Unix(),
	}
}
