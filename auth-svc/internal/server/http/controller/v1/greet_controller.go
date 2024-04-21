package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GreetController struct{}

func NewGreetController() *GreetController {
	return &GreetController{}
}

func (ctrl *GreetController) SayHello(c *gin.Context) {
	c.JSON(http.StatusOK, "Hello")
}
