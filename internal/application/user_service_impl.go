package application

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
	"user-api-restful/internal/domain"

	"github.com/oklog/ulid/v2"
)

// UserServiceImpl es la implementación concreta de la interfaz UserService.
// Contiene la lógica de negocio y orquestación para la gestión de usuarios.
type UserServiceImpl struct {
	// Repo es el contract para la persistencia de usuarios.
	Repo domain.UserRepository
	// txPort es el contract para manejar los límites transaccionales.
	txPort domain.UserTransactionPort
}

// NewUserServiceImpl crea e inicializa un nuevo UserServiceImpl.
// Recibe los contratos (interfaces) de Repositorio y Transacción, siguiendo el
// patrón de Inyección de Dependencias.
func NewUserServiceImpl(repo domain.UserRepository, tx domain.UserTransactionPort) *UserServiceImpl {
	return &UserServiceImpl{Repo: repo, txPort: tx}
}

// Asegura que UserServiceImpl implemente la interfaz UserService en tiempo de compilación.
var _ UserService = (*UserServiceImpl)(nil)

// Create valida los datos de entrada, genera un ID único (ULID) y persiste
// el nuevo usuario dentro de una transacción.
func (u *UserServiceImpl) Create(user *domain.UserCreateRequest) (*domain.User, error) {
	var createdUser *domain.User

	// Ejecuta la lógica de creación de usuario dentro de una transacción.
	err := u.txPort.Execute(func(repo domain.UserRepository) error {

		// Mapeo del DTO de entrada a la entidad de dominio.
		newUser := domain.User{
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		}

		// Generación de un ULID (ID único, ordenable por tiempo).
		t := time.Now()
		entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
		newUser.ID = ulid.MustNew(ulid.Timestamp(t), entropy).String()

		// Persistencia del nuevo usuario.
		result := repo.Create(&newUser)

		if result != nil {
			log.Printf("Estamos en create, error: %v", result)
			return result
		}

		createdUser = &newUser
		return nil
	})

	if err != nil {
		// Mapea el error de persistencia a un error de dominio/aplicación.
		return nil, u.mapRepositoryError(err)
	}

	return createdUser, nil
}

// FindAll recupera todos los usuarios del repositorio.
func (u *UserServiceImpl) FindAll() (*[]domain.User, error) {
	users, err := u.Repo.FindAll()

	if err != nil {
		// Mapea el error antes de retornarlo.
		return nil, u.mapRepositoryError(err)
	}

	return users, nil
}

// FindById recupera un usuario por su ID.
func (u *UserServiceImpl) FindById(id string) (*domain.User, error) {
	user, err := u.Repo.FindById(id)

	if err != nil {
		// Mapea el error antes de retornarlo.
		return nil, u.mapRepositoryError(err)
	}

	return user, nil
}

// Update aplica los cambios a un usuario existente dentro de una transacción.
func (u *UserServiceImpl) Update(user *domain.User) (*domain.User, error) {

	err := u.txPort.Execute(func(repo domain.UserRepository) error {
		// El repositorio se encarga de la lógica de actualización.
		err := repo.Update(user)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, u.mapRepositoryError(err)
	}

	// Retorna el usuario actualizado.
	return user, nil
}

// Delete elimina un usuario del sistema por su ID, ejecutándose dentro de una transacción.
func (u *UserServiceImpl) Delete(id string) error {
	err := u.txPort.Execute(func(repo domain.UserRepository) error {
		// El repositorio se encarga de la lógica de eliminación.
		err := repo.Delete(id)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return u.mapRepositoryError(err)
	}

	return nil
}

// mapRepositoryError traduce los errores específicos del repositorio (como los de la BD)
// a errores estándar de la capa de aplicación/dominio, asegurando que la capa de
// presentación (e.g., HTTP handlers) no dependa de detalles de persistencia.
func (u *UserServiceImpl) mapRepositoryError(err error) error {
	// Errores de "Sentinel" (comparación con errors.Is)
	if errors.Is(err, domain.ErrUserNotFound) {
		return fmt.Errorf("user not found")
	}
	if errors.Is(err, domain.ErrUsernameInUse) {
		return fmt.Errorf("username already in use")
	}
	if errors.Is(err, domain.ErrEmailInUse) {
		return fmt.Errorf("email already in use")
	}

	// Errores dinámicos (comparación con errors.As)
	var errValue domain.ErrValueNotNullable
	if errors.As(err, &errValue) {
		return fmt.Errorf(errValue.Error())
	}

	// Error interno genérico (Wrapping)
	// El %w envuelve el error original, permitiendo la inspección posterior
	// con errors.Is/As (useful para debugging/logging).
	return fmt.Errorf("internal service error: %w", err)
}
