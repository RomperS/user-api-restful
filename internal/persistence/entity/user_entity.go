package entity

import "user-api-restful/internal/domain"

// Package database contiene las estructuras (Entities) específicas de la base de datos
// y las utilidades de mapeo necesarias para la persistencia.

// UserEntity representa la estructura de la tabla de usuarios en la base de datos PostgreSQL.
// Utiliza tags de GORM para definir el esquema y las restricciones (primary key, unique index, not blank).
type UserEntity struct {
	ID       string `json:"id" gorm:"primary_key"`
	Name     string `json:"name" gorm:"not blank"`
	Username string `json:"username" gorm:"uniqueIndex:idx_username,not blank"`
	Email    string `json:"email" gorm:"uniqueIndex:idx_email,not blank"`
}

// ToEntity convierte una entidad de dominio (*domain.User) a una entidad de persistencia (UserEntity).
// Esto se utiliza antes de escribir datos en la base de datos.
func ToEntity(user *domain.User) UserEntity {
	if user == nil {
		return UserEntity{}
	}

	return UserEntity{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	}
}

// FromEntity convierte una entidad de persistencia (*UserEntity) a una entidad de dominio (domain.User).
// Esto se utiliza después de leer datos de la base de datos.
func FromEntity(entity *UserEntity) domain.User {
	return domain.User{
		ID:       entity.ID,
		Name:     entity.Name,
		Username: entity.Username,
		Email:    entity.Email,
	}
}
