package application

import (
	"context"
	"errors"
	"time"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/domain/entity"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/repository"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/usecase"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/utils"
	"github.com/Chengxufeng1994/go-saga-example/common/constant"
	"github.com/Chengxufeng1994/go-saga-example/common/model"
	"github.com/Chengxufeng1994/go-saga-example/common/token"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type AuthService struct {
	logger              *log.Entry
	userRepository      repository.UserRepository
	tokenRepository     repository.TokenRepository
	tokenEnhancer       token.Enhancer
	accessTokenExpires  int
	refreshTokenExpires int
}

func NewAuthService(
	userRepository repository.UserRepository,
	tokenRepository repository.TokenRepository,
	tokenEnhancer token.Enhancer,
	accessTokenExpires, refreshTokenExpires int,
) usecase.AuthUseCase {
	return &AuthService{
		logger:              config.ContextLogger.WithFields(log.Fields{"type": "service:AuthService"}),
		userRepository:      userRepository,
		tokenRepository:     tokenRepository,
		tokenEnhancer:       tokenEnhancer,
		accessTokenExpires:  accessTokenExpires,
		refreshTokenExpires: refreshTokenExpires,
	}
}

// SignUp implements usecase.AuthUseCase.
func (svc *AuthService) SignUp(ctx context.Context, req *dto.UserCreationRequest) (*dto.User, error) {
	hashedPassword, err := utils.HashedPassword(req.Password)
	if err != nil {
		return nil, usecase.NewAppError("HashedPassword", "app.user.create.error", nil, "").Wrap(err)
	}

	user, err := svc.userRepository.CreateUser(ctx, &entity.User{
		Active:      req.Active,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
	})

	if err != nil {
		return nil, usecase.NewAppError("CreateUser", "app.user.create.error", nil, "").Wrap(err)
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

func (svc *AuthService) createAccessToken(claims *token.Claims) (string, error) {
	claims.IsRefresh = false
	return svc.tokenEnhancer.Sign(claims)
}

func (svc *AuthService) createRefreshToken(claims *token.Claims) (string, error) {
	claims.IsRefresh = true
	return svc.tokenEnhancer.Sign(claims)
}

func (svc *AuthService) createTokenPair(claims *token.Claims) (string, string, error) {
	now := time.Now()
	accessExpiresAt := now.Add(time.Duration(svc.accessTokenExpires) * time.Second)
	refreshExpiresAt := now.Add(time.Duration(svc.refreshTokenExpires) * time.Second)
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	accessToken, err := svc.createAccessToken(claims)
	if err != nil {
		return "", "", err
	}

	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(refreshExpiresAt)
	refreshToken, err := svc.createRefreshToken(claims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// SignIn implements usecase.AuthUseCase.
func (svc *AuthService) SignIn(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	existed, err := svc.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, usecase.NewAppError("SignIn", "app.user.get_user_by_email.error", nil, "").Wrap(err)
	}

	if err := utils.ComparePassword([]byte(req.Password), []byte(existed.Password)); err != nil {
		return nil, usecase.NewAppError("SignIn", "app.auth.invalid_password.error", nil, "").Wrap(err)
	}

	claims := &token.Claims{
		UserID: existed.ID,
	}

	accessToken, refreshToken, err := svc.createTokenPair(claims)
	if err != nil {
		return nil, usecase.NewAppError("SignIn", "app.auth.gen_token.error", nil, "").Wrap(err)
	}

	err = svc.tokenRepository.StoreAccessToken(ctx, claims, accessToken, svc.accessTokenExpires)
	if err != nil {
		return nil, usecase.NewAppError("SignIn", "app.auth.store_token.error", nil, "").Wrap(err)
	}

	err = svc.tokenRepository.StoreRefreshToken(ctx, claims, refreshToken, svc.refreshTokenExpires)
	if err != nil {
		return nil, usecase.NewAppError("SignIn", "app.auth.store_token.error", nil, "").Wrap(err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// SignOut implements usecase.AuthUseCase.
func (svc *AuthService) SignOut(ctx context.Context) (string, error) {
	val := ctx.Value(constant.CtxUserKey)
	userId := val.(uint64)
	_ = svc.tokenRepository.RemoveAccessToken(ctx, userId)
	_ = svc.tokenRepository.RemoveRefreshToken(ctx, userId)
	return "Ok", nil
}

// VerifyToken implements usecase.AuthUseCase.
func (svc *AuthService) VerifyToken(ctx context.Context, req *dto.VerifyTokenRequest) (*dto.VerifyTokenResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("This is sing in event")
	defer span.End()

	claims, err := svc.tokenEnhancer.Verify(req.AccessToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, model.NewAppError("VerifyToken", "app.auth.expired_token.error", nil, "").Wrap(err)
		default:
			return nil, model.NewAppError("VerifyToken", "app.auth.verify_token.error", nil, "").Wrap(err)
		}
	}

	return &dto.VerifyTokenResponse{UserId: claims.UserID}, nil
}

// RefreshToken implements usecase.AuthUseCase.
func (svc *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	claims, err := svc.tokenEnhancer.Verify(req.RefreshToken)
	if err != nil {
		return nil, model.NewAppError("RefreshToken", "app.auth.verify_token.error", nil, "").Wrap(err)
	}
	existed, err := svc.userRepository.GetUserByID(ctx, claims.UserID)
	if err != nil {
		var nfErr *repository.ErrNotFound
		if errors.As(err, &nfErr) {
			return nil, model.NewAppError("RefreshToken", "app.user.get_by_id.error", nil, "").Wrap(err)
		}

		return nil, model.NewAppError("RefreshToken", "app.user.get_by_id.error", nil, "").Wrap(err)
	}
	if existed == nil {
		return nil, model.NewAppError("RefreshToken", "app.user.get_by_id.error", nil, "").Wrap(err)
	}

	accessToken, refreshToken, err := svc.createTokenPair(claims)
	if err != nil {
		return nil, model.NewAppError("RefreshToken", "app.auth.gen_token.error", nil, "").Wrap(err)
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
