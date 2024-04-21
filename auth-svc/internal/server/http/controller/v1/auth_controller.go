package v1

import (
	"net/http"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/usecase"
	"github.com/Chengxufeng1994/go-saga-example/common/response"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type AuthController struct {
	logger      *log.Entry
	authService usecase.AuthUseCase
}

func NewAuthController(authService usecase.AuthUseCase) *AuthController {
	logger := config.ContextLogger.WithField("controller", "user controller")

	return &AuthController{
		logger:      logger,
		authService: authService,
	}
}

func (ctrl *AuthController) SignUp(c *gin.Context) {
	var req dto.UserCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.WithError(err).Error("json marshal")
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	user, err := ctrl.authService.SignUp(c.Request.Context(), &req)
	if err != nil {
		ctrl.logger.WithError(err).Error("SignUp")
		c.AbortWithStatusJSON(http.StatusBadRequest,
			response.ErrorResponse{
				BaseResponse: &response.BaseResponse{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
				Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		response.SuccessResponse{
			BaseResponse: &response.BaseResponse{
				Code:    http.StatusOK,
				Message: "success",
			},
			Data: user})
}

func (ctrl *AuthController) SignIn(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.WithError(err).Error("json marshal")
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	resp, err := ctrl.authService.SignIn(c.Request.Context(), &req)
	if err != nil {
		ctrl.logger.WithError(err).Error("SignIn")
		c.AbortWithStatusJSON(http.StatusBadRequest,
			response.ErrorResponse{
				BaseResponse: &response.BaseResponse{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
				Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		response.SuccessResponse{
			BaseResponse: &response.BaseResponse{
				Code:    http.StatusOK,
				Message: "success",
			},
			Data: resp})
}

func (ctrl *AuthController) SignOut(c *gin.Context) {

	resp, err := ctrl.authService.SignOut(c.Request.Context())
	if err != nil {
		ctrl.logger.WithError(err).Error("SignOut")
		c.AbortWithStatusJSON(http.StatusBadRequest,
			response.ErrorResponse{
				BaseResponse: &response.BaseResponse{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
				Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		response.SuccessResponse{
			BaseResponse: &response.BaseResponse{
				Code:    http.StatusOK,
				Message: "success",
			},
			Data: resp})
}

func (ctrl *AuthController) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.WithError(err).Error("json marshal")
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	resp, err := ctrl.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		ctrl.logger.WithError(err).Error("SignOut")
		c.AbortWithStatusJSON(http.StatusBadRequest,
			response.ErrorResponse{
				BaseResponse: &response.BaseResponse{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
				Detail: err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		response.SuccessResponse{
			BaseResponse: &response.BaseResponse{
				Code:    http.StatusOK,
				Message: "success",
			},
			Data: resp})
}
