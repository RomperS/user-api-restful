// Package domain contiene las estructuras de datos fundamentales (models/entities)
// para la aplicación.
package domain

// User representa la entidad principal de un usuario en el sistema.
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserCreateRequest es la estructura utilizada para recibir datos
// al crear un nuevo usuario.
type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,excludesall= "`
	Username string `json:"username" validate:"required,excludesall= "`
	Email    string `json:"email" validate:"required,excludesall= ,email"`
}

// UserResponse es la estructura utilizada para enviar de vuelta los datos
// de un usuario al cliente.
type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
