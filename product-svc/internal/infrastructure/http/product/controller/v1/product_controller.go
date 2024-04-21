package v1

import (
	"net/http"
	"strconv"

	"github.com/Chengxufeng1994/go-saga-example/common/response"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/config"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/dto"
	"github.com/Chengxufeng1994/go-saga-example/product-svc/internal/usecase"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type ProductController struct {
	logger         *log.Entry
	productService usecase.ProductUseCase
}

func NewProductController(productService usecase.ProductUseCase) *ProductController {
	return &ProductController{
		logger: config.ContextLogger.WithFields(log.Fields{
			"type": "controller:ProductController",
		}),
		productService: productService,
	}
}

func (h *ProductController) CreateProduct(c *gin.Context) {
	var req dto.ProductCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("json marshal")
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	res, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("SignUp")
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

func (h *ProductController) ListProducts(c *gin.Context) {
	res, err := h.productService.ListProducts(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("ListProducts")
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

func (h *ProductController) GetProduct(c *gin.Context) {
	val := c.Param("id")
	id, err := strconv.Atoi(val)
	if err != nil {
		h.logger.WithError(err).Error("strconv atoi")
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	res, err := h.productService.GetProduct(c.Request.Context(), uint64(id))
	if err != nil {
		h.logger.WithError(err).Error("GetUserByID")
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
