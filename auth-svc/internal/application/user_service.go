package application

import (
	"context"
	"errors"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/repository"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/usecase"
	"github.com/Chengxufeng1994/go-saga-example/common/model"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	userRepository repository.UserRepository
	logger         *log.Entry
}

func NewUserService(userRepository repository.UserRepository) usecase.UserUseCase {
	return &UserService{
		userRepository: userRepository,
		logger:         config.ContextLogger.WithFields(log.Fields{"type": "service:UserService"}),
	}
}

// GetUserByID implements usecase.UserUseCase.
func (svc *UserService) GetUserByID(ctx context.Context, id uint64) (*dto.User, error) {
	user, err := svc.userRepository.GetUserByID(ctx, id)
	if err != nil {
		var nfErr *repository.ErrNotFound
		switch {
		case errors.As(err, &nfErr):
			return nil, model.NewAppError("GetUser", "app.user.get_by_id.error", nil, "").Wrap(err)
		default:
			return nil, model.NewAppError("GetUser", "app.user.get_by_id.error", nil, "").Wrap(err)
		}
	}

	return &dto.User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Active:      user.Active,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
	}, nil
}
