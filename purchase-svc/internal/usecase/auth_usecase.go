package usecase

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/dto"
)

type AuthUseCase interface {
	VerifyToken(context.Context, string) (*dto.VerifyTokenResponse, error)
}
