package domain

import "errors"

// Package domain contiene las estructuras de datos fundamentales (models/entities)
// y define los contracts (interfaces) y los errores específicos de la lógica de negocio.
//
// ERRORES PREDEFINIDOS
var (
	// ErrUserNotFound indica que un usuario solicitado no fue encontrado.
	ErrUserNotFound = errors.New("user not found")

	// ErrUsernameInUse indica que un nombre de usuario ya está asignado.
	ErrUsernameInUse = errors.New("username already in use")
	// ErrEmailInUse indica que una dirección de correo electrónico ya está registrada.
	ErrEmailInUse = errors.New("email already in use")
	// ErrIdInUse indica que un identificador proporcionado ya está en uso.
	ErrIdInUse = errors.New("id already in use")
)

// ErrValueNotNullable representa un error cuando se intenta dejar nulo
// un campo que requiere un valor.
type ErrValueNotNullable struct {
	Value string
}

// Error implementa la interface error para ErrValueNotNullable.
func (e ErrValueNotNullable) Error() string {
	return e.Value + " is not nullable"
}

// ErrInternalServer representa un fallo inesperado del servidor.
// No debe ser retornado directamente a un cliente, sino logueado.
type ErrInternalServer struct {
	Value string
}

// Error implementa la interface error para ErrInternalServer.
func (e ErrInternalServer) Error() string {
	return "internal server error: " + e.Value
}

// ErrTransactionFailed representa un fallo al intentar completar una Unit of Work.
// Esto suele ser el resultado de un rollback de la base de datos.
type ErrTransactionFailed struct {
	Value string
}

// Error implementa la interface error para ErrTransactionFailed.
func (e ErrTransactionFailed) Error() string {
	return "error transaction failed: " + e.Value
}
