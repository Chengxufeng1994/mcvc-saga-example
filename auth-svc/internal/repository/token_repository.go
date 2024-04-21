package repository

import (
	"context"

	"github.com/Chengxufeng1994/go-saga-example/common/token"
)

// token repository
type TokenRepository interface {
	GetAccessToken(context.Context, uint64) (string, error)
	StoreAccessToken(context.Context, *token.Claims, string, int) error
	RemoveAccessToken(context.Context, uint64) error
	StoreRefreshToken(context.Context, *token.Claims, string, int) error
	RemoveRefreshToken(context.Context, uint64) error
}
