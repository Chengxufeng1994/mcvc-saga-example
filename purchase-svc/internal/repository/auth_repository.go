package repository

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/domain"
)

type AuthRepository interface {
	VerifyToken(context.Context, string) (*domain.Auth, error)
}
