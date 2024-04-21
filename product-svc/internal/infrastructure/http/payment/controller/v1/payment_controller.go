package v1

import (
	"net/http"
	"strconv"

	"github.com/Chengxufeng1994/go-saga-example/common/constant"
	"github.com/Chengxufeng1994/go-saga-example/common/response"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/config"
	infrahttp "github.com/Chengxufeng1994/go-saga-example/product-svc/internal/infrastructure/http"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/usecase"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type PaymentController struct {
	logger         *log.Entry
	paymentService usecase.PaymentUseCase
}

func NewPaymentController(paymentService usecase.PaymentUseCase) *PaymentController {
	return &PaymentController{
		logger: config.ContextLogger.WithFields(log.Fields{
			"type": "controller:PaymentController",
		}),
		paymentService: paymentService,
	}
}

func (h *PaymentController) GetPayment(c *gin.Context) {
	userId, ok := c.Request.Context().Value(constant.CtxUserKey).(uint64)
	if !ok {
		resp := response.ErrorResponse{
			BaseResponse: &response.BaseResponse{
				Code:    http.StatusUnauthorized,
				Message: infrahttp.ErrUnauthorized.Error(),
			},
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
		return
	}

	val := c.Param("id")
	id, err := strconv.Atoi(val)
	if err != nil {
		h.logger.WithError(err).Error("strconv atoi")
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	res, err := h.paymentService.GetPayment(c.Request.Context(), userId, uint64(id))
	if err != nil {
		h.logger.WithError(err).Error("GetPayment")
		c.AbortWithStatusJSON(http.StatusBadRequest,
			response.ErrorResponse{
				BaseResponse: &response.BaseResponse{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
				Detail: err.Error()})
		return
	}
	c.JSON(http.StatusOK, generateResponse(res))
}

func generateResponse(res any) response.SuccessResponse {
	return response.SuccessResponse{
		BaseResponse: &response.BaseResponse{
			Code:    http.StatusOK,
			Message: "success",
		},
		Data: res}

}
