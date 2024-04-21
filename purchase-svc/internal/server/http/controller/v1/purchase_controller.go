package v1

import (
	"errors"
	"net/http"

	"github.com/Chengxufeng1994/go-saga-example/common/constant"
	"github.com/Chengxufeng1994/go-saga-example/common/response"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/purchase-svc/internal/usecase"
	"github.com/gin-gonic/gin"
)

var (
	// ErrInvalidParam is invalid parameter error
	ErrInvalidParam = errors.New("invalid parameter")
	// ErrUnauthorized is unauthorized error
	ErrUnauthorized = errors.New("unauthorized")
	// ErrServer is server error
	ErrServer = errors.New("server error")
)

type PurchaseController struct {
	purchaseService usecase.PurchaseUseCase
}

func NewPurchaseController(purchaseService usecase.PurchaseUseCase) *PurchaseController {
	return &PurchaseController{
		purchaseService: purchaseService,
	}
}

func (ctrl *PurchaseController) CreatePurchase(c *gin.Context) {
	var req dto.PurchaseCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := response.ErrorResponse{
			BaseResponse: &response.BaseResponse{
				Code:    http.StatusUnauthorized,
				Message: ErrInvalidParam.Error(),
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	userId, ok := c.Request.Context().Value(constant.CtxUserKey).(uint64)
	if !ok {
		resp := response.ErrorResponse{
			BaseResponse: &response.BaseResponse{
				Code:    http.StatusUnauthorized,
				Message: ErrUnauthorized.Error(),
			},
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
		return
	}

	data, err := ctrl.purchaseService.CreatePurchase(c.Request.Context(), userId, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse{
		BaseResponse: &response.BaseResponse{
			Code:    http.StatusOK,
			Message: "Ok",
		},
		Data: data,
	})
}

func (ctrl *PurchaseController) GetResult(c *gin.Context) {
	c.JSON(http.StatusOK, response.OkMessage)
}
