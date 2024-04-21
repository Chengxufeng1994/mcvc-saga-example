package usecase

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/dto"
)

// UserUseCase defines user data related interface
type AuthUseCase interface {
	SignUp(context.Context, *dto.UserCreationRequest) (*dto.User, error)
	SignIn(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	SignOut(context.Context) (string, error)
	VerifyToken(context.Context, *dto.VerifyTokenRequest) (*dto.VerifyTokenResponse, error)
	RefreshToken(context.Context, *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
}
