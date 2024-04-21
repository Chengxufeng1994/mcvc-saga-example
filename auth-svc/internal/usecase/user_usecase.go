package usecase

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/dto"
)

// UserUseCase defines user data related interface
type UserUseCase interface {
	GetUserByID(context.Context, uint64) (*dto.User, error)
}
