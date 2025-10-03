package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"user-api-restful/internal/application"
	"user-api-restful/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// Package http define los controladores (handlers) para la API REST.

// UserHandler maneja todas las peticiones HTTP relacionadas con la gestión de usuarios.
// Depende de la interfaz application.UserService para la lógica de negocio.
type UserHandler struct {
	userService application.UserService // Contract de la lógica de negocio.
	validator   *validator.Validate     // Instancia del validador para DTOs.
}

// NewUserHandler crea una nueva instancia de UserHandler con el servicio de usuario inyectado.
func NewUserHandler(service application.UserService) *UserHandler {
	return &UserHandler{
		userService: service,
		validator:   validator.New(),
	}
}

// CreateUser maneja la petición POST para crear un nuevo usuario.
// Se encarga de la deserialización (JSON), la validación del request body,
// el llamado al servicio y el mapeo de errores de dominio a respuestas HTTP.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) *HTTPError {
	var request domain.UserCreateRequest

	// 1. Deserialización JSON
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHTTPError(errors.New("invalid request body format"), http.StatusBadRequest)
	}

	// 2. Validación de la estructura
	err = h.validator.Struct(request)

	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldError := range validationErrors {
				switch fieldError.Tag() {
				case "required", "excludesall":
					return NewHTTPError(
						errors.New(fieldError.Field()+" is required and cannot be blank."),
						http.StatusBadRequest,
					)
				case "email":
					return NewHTTPError(
						errors.New("email format is invalid"),
						http.StatusBadRequest,
					)
				default:
					return NewHTTPError(
						errors.New("Validation failed on field: "+fieldError.Field()),
						http.StatusBadRequest,
					)
				}
			}
		}
		// Fallo inesperado durante la validación
		return NewHTTPError(errors.New(err.Error()), http.StatusInternalServerError)
	}

	// 3. Llamada al servicio de aplicación
	userResponse, err := h.userService.Create(&request)

	// 4. Mapeo de errores de dominio a HTTP Status Codes
	if err != nil {
		if errors.Is(err, domain.ErrIdInUse) || errors.Is(err, domain.ErrEmailInUse) || errors.Is(err, domain.ErrUsernameInUse) {
			// 409 Conflict para errores de unicidad/recurso existente.
			return NewHTTPError(errors.New(err.Error()), http.StatusConflict)
		}
		if errors.Is(err, domain.ErrValueNotNullable{}) {
			// 400 Bad Request para errores de validación (valores nulos/vacíos).
			return NewHTTPError(errors.New(err.Error()), http.StatusBadRequest)
		}
		// 500 Internal Server Error para cualquier otro fallo.
		return NewHTTPError(errors.New(err.Error()), http.StatusInternalServerError)
	}

	// 5. Respuesta exitosa (201 Created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(userResponse)
	if err != nil {
		return NewHTTPError(errors.New("error json encoding response"), http.StatusInternalServerError)
	}

	return nil
}

// FindAll maneja la petición GET para obtener todos los usuarios.
func (h *UserHandler) FindAll(w http.ResponseWriter, r *http.Request) *HTTPError {
	// Llamada al servicio
	userResponse, err := h.userService.FindAll()

	if err != nil {
		return NewHTTPError(errors.New(err.Error()), http.StatusInternalServerError)
	}

	// Respuesta exitosa (200 OK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(userResponse)
	if err != nil {
		return NewHTTPError(errors.New("error json encoding response"), http.StatusInternalServerError)
	}

	return nil
}

// FindById maneja la petición GET para obtener un usuario por ID.
func (h *UserHandler) FindById(w http.ResponseWriter, r *http.Request) *HTTPError {
	// 1. Extracción del parámetro de la URL
	id := chi.URLParam(r, "id")

	if id == "" {
		return NewHTTPError(errors.New("user ID is required in the request path or query"), http.StatusBadRequest)
	}

	// 2. Llamada al servicio
	userResponse, err := h.userService.FindById(id)

	// 3. Mapeo de errores
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// 404 Not Found para recurso no encontrado.
			return NewHTTPError(errors.New(err.Error()), http.StatusNotFound)
		}
		return NewHTTPError(errors.New(err.Error()), http.StatusInternalServerError)
	}

	// 4. Respuesta exitosa (200 OK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(userResponse)
	if err != nil {
		return NewHTTPError(errors.New("error json encoding response"), http.StatusInternalServerError)
	}

	return nil
}

// Update maneja la petición PUT para actualizar un usuario.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) *HTTPError {
	var request domain.User

	// 1. Deserialización JSON
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHTTPError(errors.New("invalid request body format"), http.StatusBadRequest)
	}

	// 2. Validación (más simple aquí)
	err = h.validator.Struct(request)
	if err != nil {
		return NewHTTPError(errors.New("validation failed on update fields"), http.StatusBadRequest)
	}

	// 3. Llamada al servicio
	userResponse, err := h.userService.Update(&request)

	// 4. Mapeo de errores
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return NewHTTPError(errors.New(err.Error()), http.StatusNotFound)
		}
		if errors.Is(err, domain.ErrEmailInUse) || errors.Is(err, domain.ErrUsernameInUse) {
			return NewHTTPError(errors.New(err.Error()), http.StatusConflict)
		}
		return NewHTTPError(errors.New(err.Error()), http.StatusInternalServerError)
	}

	// 5. Respuesta exitosa (200 OK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK for successful update

	err = json.NewEncoder(w).Encode(userResponse)
	if err != nil {
		return NewHTTPError(errors.New("error json encoding response"), http.StatusInternalServerError)
	}

	return nil
}

// Delete maneja la petición DELETE para eliminar un usuario por ID.
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) *HTTPError {
	// 1. Extracción del parámetro de la URL
	id := chi.URLParam(r, "id")

	if id == "" {
		return NewHTTPError(errors.New("user ID is required in the request path or query"), http.StatusBadRequest)
	}

	// 2. Llamada al servicio
	err := h.userService.Delete(id)

	// 3. Mapeo de errores
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// 404 Not Found si el recurso a eliminar no existe.
			return NewHTTPError(errors.New(err.Error()), http.StatusNotFound)
		}
		return NewHTTPError(errors.New(err.Error()), http.StatusInternalServerError)
	}

	// 4. Respuesta exitosa (204 No Content)
	w.WriteHeader(http.StatusNoContent)

	return nil
}
