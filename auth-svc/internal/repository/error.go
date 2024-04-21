package repository

import (
	"fmt"

	"gorm.io/gorm"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

// ErrInvalidInput indicates an error that has occurred due to an invalid input.
type ErrInvalidInput struct {
	Entity  string // The entity which was sent as the input.
	Field   string // The field of the entity which was invalid.
	Value   any    // The actual value of the field.
	wrapped error  // The original error
}

func NewErrInvalidInput(entity, field string, value any) *ErrInvalidInput {
	return &ErrInvalidInput{
		Entity: entity,
		Field:  field,
		Value:  value,
	}
}

func (e *ErrInvalidInput) Error() string {
	if e.wrapped != nil {
		return fmt.Sprintf("invalid input: entity: %s field: %s value: %s error: %s", e.Entity, e.Field, e.Value, e.wrapped)
	}

	return fmt.Sprintf("invalid input: entity: %s field: %s value: %s", e.Entity, e.Field, e.Value)
}

func (e *ErrInvalidInput) Wrap(err error) *ErrInvalidInput {
	e.wrapped = err
	return e
}

func (e *ErrInvalidInput) Unwrap() error {
	return e.wrapped
}

func (e *ErrInvalidInput) InvalidInputInfo() (entity string, field string, value any) {
	entity = e.Entity
	field = e.Field
	value = e.Value
	return
}

// ErrConflict indicates a conflict that occurred.
type ErrConflict struct {
	Resource string // The resource which created the conflict.
	err      error  // Internal error.
	meta     string // Any additional metadata.
}

func NewErrConflict(resource string, err error, meta string) *ErrConflict {
	return &ErrConflict{
		Resource: resource,
		err:      err,
		meta:     meta,
	}
}

func (e *ErrConflict) Error() string {
	msg := e.Resource + "exists " + e.meta
	if e.err != nil {
		msg += " " + e.err.Error()
	}
	return msg
}

func (e *ErrConflict) Unwrap() error {
	return e.err
}

// IsErrConflict allows easy type assertion without adding store as a dependency.
func (e *ErrConflict) IsErrConflict() bool {
	return true
}

// ErrNotFound indicates that a resource was not found
type ErrNotFound struct {
	resource string
	ID       string
	wrapped  error
}

func NewErrNotFound(resource, id string) *ErrNotFound {
	return &ErrNotFound{
		resource: resource,
		ID:       id,
	}
}

func (e *ErrNotFound) Wrap(err error) *ErrNotFound {
	e.wrapped = err
	return e
}

func (e *ErrNotFound) Error() string {
	if e.wrapped != nil {
		return fmt.Sprintf("resource: %s id: %s error: %s", e.resource, e.ID, e.wrapped)
	}

	return fmt.Sprintf("resource: %s id: %s", e.resource, e.ID)
}

// ErrFailedCreate indicates an error that has occurred due to an failed create.
type ErrFailedCreate struct {
	Entity  string // The entity which was sent as the input.
	wrapped error  // The original error
}

func NewErrFailedCreate(entity string) *ErrFailedCreate {
	return &ErrFailedCreate{
		Entity: entity,
	}
}

func (e *ErrFailedCreate) Error() string {
	if e.wrapped != nil {
		return fmt.Sprintf("failed create: entity: %s error: %s", e.Entity, e.wrapped)
	}

	return fmt.Sprintf("failed create: entity: %s", e.Entity)
}

func (e *ErrFailedCreate) Wrap(err error) *ErrFailedCreate {
	e.wrapped = err
	return e
}

func (e *ErrFailedCreate) Unwrap() error {
	return e.wrapped
}

func (e *ErrFailedCreate) FailedCreateInfo() (entity string) {
	entity = e.Entity
	return
}
