package controller

import (
	"net/http"

	"github.com/LoginX/SprayDash/internal/utils"
	"github.com/gin-gonic/gin"
)

type RootController struct{}

func NewRootController() *RootController {
	return &RootController{}
}

func (rc *RootController) RegisterRoutes(rg *gin.Engine) {
	rg.GET("/", rc.Home)
	rg.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, utils.Response(http.StatusNotFound, nil, "Route not found"))
	})
}

func (rc *RootController) Home(ctx *gin.Context) {
	// root route for health

	ctx.JSON(http.StatusOK, utils.Response(http.StatusOK, nil, "SprayDash is up and running"))

}
