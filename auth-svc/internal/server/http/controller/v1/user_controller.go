package v1

import (
	"net/http"
	"strconv"

	"github.com/Chengxufeng1994/go-saga-example/auth-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/auth-svc/internal/usecase"
	"github.com/Chengxufeng1994/go-saga-example/common/response"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type UserController struct {
	logger      *log.Entry
	userService usecase.UserUseCase
}

func NewUserController(userService usecase.UserUseCase) *UserController {
	logger := config.ContextLogger.WithField("type", "controller:UserController")

	return &UserController{
		logger:      logger,
		userService: userService,
	}
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	val := c.Param("id")
	id, err := strconv.Atoi(val)
	if err != nil {
		ctrl.logger.WithError(err).Error("strconv atoi")
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	user, err := ctrl.userService.GetUserByID(c.Request.Context(), uint64(id))
	if err != nil {
		ctrl.logger.WithError(err).Error("GetUserByID")
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
