// Package domain contiene las estructuras de datos fundamentales (models/entities)
// y define los contracts (interfaces) para la lógica de negocio.
package domain

// UserRepository define el contract para la persistencia de datos de usuario.
// Esta interface desacopla la lógica de negocio del almacenamiento de datos
// (como una base de datos o un servicio externo).
type UserRepository interface {
	// Create inserta un nuevo User en el almacenamiento.
	// Retorna un error si la operación falla (e.g., conflicto de ID o conexión).
	Create(user *User) error
	// FindAll recupera todos los usuarios del almacenamiento.
	FindAll() (*[]User, error)
	// FindById recupera un User por su identificador único (ID).
	// Retorna nil si no se encuentra el usuario.
	FindById(id string) (*User, error)
	// Update aplica los cambios a un User existente en el almacenamiento.
	// Retorna un error si la operación falla (e.g., el usuario no existe).
	Update(user *User) error
	// Delete elimina un User permanentemente del almacenamiento usando su ID.
	Delete(id string) error
}
