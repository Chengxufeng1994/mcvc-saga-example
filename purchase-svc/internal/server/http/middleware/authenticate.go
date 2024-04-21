package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Chengxufeng1994/go-saga-example/common/constant"
	"github.com/Chengxufeng1994/go-saga-example/common/model"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/usecase"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

type JwtAuthenticator struct {
	logger      *log.Entry
	authService usecase.AuthUseCase
}

func NewJwtAuthenticator(logger *config.Logger, authService usecase.AuthUseCase) *JwtAuthenticator {
	return &JwtAuthenticator{
		logger: logger.ContextLogger.WithFields(log.Fields{
			"type": "middleware:JwtAuthenticator",
		}),
		authService: authService,
	}
}

func (authenticator *JwtAuthenticator) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := extractToken(c.Request)
		if accessToken == "" {
			err := model.NewAppError("AuthMiddleware", "app.auth.invalid_token.error", nil, "")
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		verifyTokenResponse, err := authenticator.authService.VerifyToken(c.Request.Context(), accessToken)
		if err != nil {
			authenticator.logger.WithError(err).Error("verify token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		if verifyTokenResponse.IsExpired {
			authenticator.logger.WithError(err).Error("expired token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), constant.CtxUserKey, verifyTokenResponse.UserId))
		c.Next()
	}
}
