package repository

import (
	"context"

	entity "github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/domain/entity"
)

// user repository
type UserRepository interface {
	WithTx(fn GormOption) error
	CreateUser(context.Context, *entity.User) (*entity.User, error)
	GetUserByID(context.Context, uint64) (*entity.User, error)
	GetUserByEmail(context.Context, string) (*entity.User, error)
}
