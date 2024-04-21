package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/domain/entity"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrFailedToCreateUser = errors.New("failed to create user")

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// WithTx implements repository.UserRepository.
func (u *UserRepository) WithTx(fn repository.GormOption) error {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error

}

// CreateUser implements repository.UserRepository.
func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	result := r.db.WithContext(ctx).Model(&entity.User{}).Create(user)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			switch pgErr.Code {
			case repository.ForeignKeyViolation:
			case repository.UniqueViolation:
				if pgErr.ConstraintName == "uni_users_email" {
					return nil, repository.NewErrInvalidInput("User", "email", user.Email).Wrap(result.Error)
				}
			}
		}

		return nil, repository.NewErrFailedCreate("User")
	}

	return user, nil
}

// GetUserByID implements repository.UserRepository.
func (r *UserRepository) GetUserByID(ctx context.Context, id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, repository.NewErrNotFound("User", fmt.Sprintf("%d", id))
		}

		return nil, err
	}

	return &user, nil
}

// GetUserByEmail implements repository.UserRepository.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, repository.NewErrNotFound("User", fmt.Sprintf("email=%s", email))
		}

		return nil, err
	}

	return &user, nil
}
