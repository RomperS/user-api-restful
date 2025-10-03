package application

import "user-api-restful/internal/domain"

// Package application define las interfaces y estructuras de los servicios
// de la aplicación que contienen la lógica de negocio principal.
//
// UserService define el contract para las operaciones de negocio relacionadas
// con la gestión de usuarios. Actúa como orquestador entre el puerto de entrada
// (e.g., HTTP handler) y la capa de dominio/persistencia.
type UserService interface {
	// Create valida los datos de entrada y persiste un nuevo usuario.
	// Retorna la entidad User creada y puede retornar errores como
	// ErrUsernameInUse o ErrEmailInUse.
	Create(user *domain.UserCreateRequest) (*domain.User, error)
	// FindAll recupera la lista completa de todos los usuarios.
	FindAll() (*[]domain.User, error)
	// FindById recupera un usuario específico utilizando su ID.
	// Retorna ErrUserNotFound si el usuario no existe.
	FindById(id string) (*domain.User, error)
	// Update aplica los cambios al usuario proporcionado.
	// Retorna ErrUserNotFound si el usuario a actualizar no existe.
	Update(user *domain.User) (*domain.User, error)
	// Delete elimina un usuario del sistema por su ID.
	Delete(id string) error
}
