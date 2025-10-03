package database

import (
	"errors"
	"log"
	"user-api-restful/internal/domain"
	"user-api-restful/internal/persistence/entity"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Package database contiene las implementaciones de los contratos de repositorio (UserRepository)
// y transacción (UserTransactionPort) utilizando GORM y PostgreSQL.
//
// PostgresRepository implementa las interfaces domain.UserRepository y
// domain.UserTransactionPort para la persistencia de usuarios en PostgreSQL.
type PostgresRepository struct {
	db *gorm.DB
}

// PgErrorData es una estructura auxiliar para manejar y tipificar errores de PostgreSQL.
type PgErrorData struct {
	Code       string
	Constraint string
}

// extractPgError intenta extraer información estructurada (código y constraint)
// de errores específicos de PostgreSQL (pq.Error o pgconn.PgError).
// Retorna nil si no es un error de PG o si es un error de GORM conocido.
func extractPgError(err error) *PgErrorData {
	if err == nil {
		return nil
	}

	// Intenta extraer errores del driver pq (usado a menudo con GORM).
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		constraintName := pqErr.Constraint
		if constraintName == "" {
			constraintName = pqErr.Table
		}

		return &PgErrorData{
			Code:       string(pqErr.Code),
			Constraint: constraintName}
	}

	// Intenta extraer errores del driver pgx (otra opción popular).
	var pgconnErr *pgconn.PgError
	if errors.As(err, &pgconnErr) {
		return &PgErrorData{
			Code:       pgconnErr.Code,
			Constraint: pgconnErr.ConstraintName}
	}

	// Ignora errores de GORM que no requieren mapeo a PgErrorData.
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, gorm.ErrInvalidData) {
		return nil
	}

	return nil
}

// NewPostgresRepository crea una nueva instancia del repositorio, inyectando la conexión a GORM.
func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create inserta un nuevo usuario. Mapea errores de unicidad (23505) y
// not-null (23502) de PostgreSQL a los errores de dominio.
func (p *PostgresRepository) Create(user *domain.User) error {
	userEntity := entity.ToEntity(user)

	result := p.db.Create(&userEntity).Error

	if result != nil {
		err := extractPgError(result)
		if (err != nil) && (err.Code == "23505") { // Código de violación de Unique/Primary Key
			switch err.Constraint {
			case "users_pkey":
				return domain.ErrIdInUse
			case "idx_username":
				return domain.ErrUsernameInUse
			case "idx_email":
				return domain.ErrEmailInUse
			}
		}
		if (err != nil) && (err.Code == "23502") { // Código de violación Not Null
			column := ""
			switch err.Constraint {
			case "id", "username", "email":
				column = err.Constraint
			default:
				column = "a column"
			}
			return domain.ErrValueNotNullable{Value: column}
		}
		// Cualquier otro error de persistencia se mapea como error interno.
		return domain.ErrInternalServer{Value: result.Error()}
	}

	return nil
}

// FindAll recupera todos los registros de usuario y los mapea a entidades de dominio.
func (p *PostgresRepository) FindAll() (*[]domain.User, error) {
	var userEntities []entity.UserEntity

	err := p.db.Find(&userEntities).Error

	if err != nil {
		return nil, domain.ErrInternalServer{Value: err.Error()}
	}

	// Mapeo de entidades de persistencia a entidades de dominio (Domain Entities).
	users := make([]domain.User, len(userEntities))

	for i, targetEntity := range userEntities {
		users[i] = entity.FromEntity(&targetEntity)
	}

	return &users, nil
}

// FindById recupera un usuario por su ID. Mapea gorm.ErrRecordNotFound a domain.ErrUserNotFound.
func (p *PostgresRepository) FindById(id string) (*domain.User, error) {
	var userEntity entity.UserEntity

	err := p.db.Where("id = ?", id).First(&userEntity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrInternalServer{Value: err.Error()}
	}

	user := entity.FromEntity(&userEntity)

	return &user, nil
}

// Update actualiza un usuario existente. Realiza un mapeo similar a Create
// para errores de unicidad y not-null.
func (p *PostgresRepository) Update(user *domain.User) error {
	userEntity := entity.ToEntity(user)

	// Usa Model y Updates para actualizar solo los campos provistos y basándose en el ID.
	result := p.db.Model(&entity.UserEntity{ID: userEntity.ID}).Updates(userEntity)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}

		// Mapeo de errores de PostgreSQL.
		err := extractPgError(result.Error)
		if (err != nil) && (err.Code == "23505") {
			switch err.Constraint {
			case "users_pkey":
				return domain.ErrIdInUse
			case "idx_username":
				return domain.ErrUsernameInUse
			case "idx_email":
				return domain.ErrEmailInUse
			}
		}

		if (err != nil) && (err.Code == "23502") {
			column := ""
			switch err.Constraint {
			case "id", "username", "email":
				column = err.Constraint
			default:
				column = "a column"
			}
			return domain.ErrValueNotNullable{Value: column}
		}
		return domain.ErrInternalServer{Value: result.Error.Error()}
	}

	// Si GORM no reportó error, pero ninguna fila fue afectada, significa que el usuario no existía.
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete elimina un usuario por su ID.
func (p *PostgresRepository) Delete(id string) error {
	userToDelete := entity.UserEntity{ID: id}

	result := p.db.Delete(&userToDelete)

	if result.Error != nil {
		return domain.ErrInternalServer{Value: result.Error.Error()}
	}

	// Si RowsAffected es cero, el usuario no fue encontrado para eliminar.
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Execute implementa el UserTransactionPort, ejecutando la función de dominio
// dentro de una transacción de GORM.
func (p *PostgresRepository) Execute(fn func(repo domain.UserRepository) error) error {
	var capturedDomainError error

	// Inicia una transacción de GORM.
	txErr := p.db.Transaction(func(tx *gorm.DB) error {
		// Crea una nueva instancia de repositorio que usa la transacción (txRepo).
		txRepo := &PostgresRepository{db: tx}

		// Ejecuta la lógica de negocio, pasando el repositorio transaccional.
		txResultErr := fn(txRepo)

		if txResultErr != nil {
			// Captura el error de dominio para retornarlo posteriormente,
			// forzando un Rollback al retornar el error aquí.
			capturedDomainError = txResultErr
			return txResultErr
		}

		// Si no hay error, GORM hace Commit.
		return nil
	})

	if txErr != nil {
		log.Printf("[Transaction Failed] Database Error: %v", txErr)

		// Si la transacción falló debido a un error de dominio, retorna ese error.
		if capturedDomainError != nil {
			return capturedDomainError
		}

		// Si falló por un error de conexión o base de datos, retorna un error de transacción.
		return domain.ErrTransactionFailed{Value: txErr.Error()}
	}

	return nil
}
