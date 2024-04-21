// redis
package redis

import (
	"context"
	"fmt"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/repository"
	libredis "github.com/Chengxufeng1994/go-saga-example/common/redis"
	"github.com/Chengxufeng1994/go-saga-example/common/token"
)

type TokenRepository struct {
	rc libredis.RedisCache
}

func (r *TokenRepository) buildAccessTokenKey(userID uint64) string {
	return fmt.Sprintf("user:%d:access_token", userID)
}

func (r *TokenRepository) buildRefreshTokenKey(userID uint64) string {
	return fmt.Sprintf("user:%d:refresh_token", userID)
}

func NewTokenRepository(rc libredis.RedisCache) repository.TokenRepository {
	return &TokenRepository{
		rc: rc,
	}
}

// GetAccessToken implements repository.TokenRepository.
func (r *TokenRepository) GetAccessToken(ctx context.Context, userId uint64) (string, error) {
	key := r.buildAccessTokenKey(userId)
	var val string
	ok, err := r.rc.Get(ctx, key, &val)
	if !ok || err != nil {
		return "", err
	}

	return val, nil
}

// StoreAccessToken implements repository.TokenRepository.
func (r *TokenRepository) StoreAccessToken(ctx context.Context, claims *token.Claims, value string, ttl int) error {
	key := r.buildAccessTokenKey(claims.UserID)
	return r.rc.Set(ctx, key, value, ttl)
}

// RemoveAccessToken implements repository.TokenRepository.
func (r *TokenRepository) RemoveAccessToken(ctx context.Context, userID uint64) error {
	key := r.buildAccessTokenKey(userID)
	return r.rc.Del(ctx, key)
}

// StoreRefreshToken implements repository.TokenRepository.
func (r *TokenRepository) StoreRefreshToken(ctx context.Context, claims *token.Claims, value string, ttl int) error {
	key := r.buildRefreshTokenKey(claims.UserID)
	return r.rc.Set(ctx, key, value, ttl)
}

// RemoveRefreshToken implements repository.TokenRepository.
func (r *TokenRepository) RemoveRefreshToken(ctx context.Context, userID uint64) error {
	key := r.buildRefreshTokenKey(userID)
	return r.rc.Del(ctx, key)
}
