package server

import (
	"net/http"
	"time"
)

type JSONSerializationConfig struct {
	// PrettyPrint active l'indentation JSON
	PrettyPrint bool

	// IncludeTimestamp ajoute automatiquement un timestamp
	IncludeTimestamp bool

	// WrapSingleValues encapsule les valeurs simples dans APIResponse
	WrapSingleValues bool

	// ErrorWrapper fonction personnalisée pour wrapper les erreurs
	ErrorWrapper func(error, *http.Request) interface{}

	// SuccessWrapper fonction personnalisée pour wrapper les succès
	SuccessWrapper func(interface{}, *http.Request) interface{}

	// ContentType définit le Content-Type (défaut: application/json)
	ContentType string
}

var DefaultJSONConfig = JSONSerializationConfig{
	PrettyPrint:      false,
	IncludeTimestamp: true,
	WrapSingleValues: true,
	ContentType:      "application/json; charset=utf-8",

	ErrorWrapper: func(err error, r *http.Request) interface{} {
		return ErrorResponse{
			Success:   false,
			Message:   err.Error(),
			Timestamp: time.Now().Unix(),
		}
	},

	SuccessWrapper: func(data interface{}, r *http.Request) interface{} {
		return SuccessResponse{
			Success:   true,
			Data:      data,
			Timestamp: time.Now().Unix(),
		}
	},
}
