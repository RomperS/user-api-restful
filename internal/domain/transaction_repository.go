// Package domain contiene las estructuras de datos fundamentales (models/entities)
// y define los contracts (interfaces) para la lógica de negocio.
package domain

// UserTransactionPort define el contract para manejar transacciones
// a través de la capa de persistencia.
// Su propósito principal es asegurar que un conjunto de operaciones de repositorio
// se ejecuten de forma atómica (commit o rollback).
type UserTransactionPort interface {
	// Execute ejecuta la función 'fn' dentro de una única transacción.
	// La función 'fn' recibe una instancia de UserRepository que está
	// enlazada a la transacción actual. Si 'fn' retorna un error, la transacción
	// debe ser revertida (rollback); de lo contrario, se confirma (commit).
	Execute(fn func(repo UserRepository) error) error
}
