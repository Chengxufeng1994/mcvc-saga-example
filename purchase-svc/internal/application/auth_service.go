package application

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/common/model"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/repository"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/usecase"
	log "github.com/sirupsen/logrus"
)

type AuthService struct {
	logger         *log.Entry
	authRepository repository.AuthRepository
}

func NewAuthService(logger *config.Logger, authRepository repository.AuthRepository) usecase.AuthUseCase {
	return &AuthService{
		logger:         logger.ContextLogger.WithFields(log.Fields{"type": "service:AuthService"}),
		authRepository: authRepository,
	}
}

// VerifyToken implements usecase.AuthUseCase.
func (s *AuthService) VerifyToken(ctx context.Context, accessToken string) (*dto.VerifyTokenResponse, error) {
	res, err := s.authRepository.VerifyToken(ctx, accessToken)
	if err != nil {
		return nil, model.NewAppError("VerifyToken", "app.auth.verify_token.error", nil, "")
	}

	return &dto.VerifyTokenResponse{
		UserId:    res.UserId,
		IsExpired: res.IsExpired,
	}, nil
}
