package http

import (
	"encoding/json"
	"net/http"
)

// Package http define los controladores (handlers) y utilidades específicas
// de la capa de presentación HTTP.

// ErrorResponse es la estructura estándar utilizada para retornar información
// detallada de un error al cliente.
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// HandlerFunc define la firma de los handlers de nuestra aplicación.
// A diferencia de http.HandlerFunc, retorna un *HTTPError para el manejo explícito de errores.
type HandlerFunc func(http.ResponseWriter, *http.Request) *HTTPError

// HTTPError es una estructura que envuelve un error estándar de Go junto con
// el código de estado HTTP apropiado para el cliente.
type HTTPError struct {
	Error  error
	Status int
}

// NewHTTPError es la función constructora para crear una nueva instancia de HTTPError.
func NewHTTPError(err error, status int) *HTTPError {
	return &HTTPError{Error: err, Status: status}
}

// ErrorHandlerWrapper es un decorador de handler que implementa un manejo
// de errores centralizado.
//
// Transforma una HandlerFunc que retorna un error de la aplicación a una
// http.HandlerFunc estándar, interceptando el error de la aplicación y
// generando una respuesta JSON consistente para el cliente.
func ErrorHandlerWrapper(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Llama al handler real de la aplicación.
		err := handler(w, r)

		// Si el handler retorna un error, se procesa aquí.
		if err != nil {
			statusCode := http.StatusInternalServerError
			if err.Status != 0 {
				statusCode = err.Status
			}

			// Construye la respuesta de error JSON.
			response := ErrorResponse{
				Status:  statusCode,
				Message: err.Error.Error(), // Usa el mensaje del error envuelto.
			}

			// Establece las cabeceras y escribe el código de estado.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)

			// Codifica y escribe el cuerpo de la respuesta JSON.
			_ = json.NewEncoder(w).Encode(response)
		}
	}
}
