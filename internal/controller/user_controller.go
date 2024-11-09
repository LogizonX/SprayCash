package controller

import (
	"log"
	"net/http"

	"github.com/LoginX/SprayDash/internal/service"
	"github.com/LoginX/SprayDash/internal/service/dto"
	"github.com/LoginX/SprayDash/internal/utils"
	"github.com/gin-gonic/gin"
)

// userController depends on user service
type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/register", uc.Register)
	rg.POST("/auth/login", uc.Login)
}

func (uc *UserController) Register(c *gin.Context) {
	// bind payload
	var pl dto.CreateUserDTO

	if err := c.ShouldBindJSON(&pl); err != nil {
		// log the error
		log.Println(err)
		c.JSON(400, utils.Response(http.StatusBadRequest, nil, err.Error()))
		return
	}
	// send to service
	resp, rErr := uc.userService.Register(pl)
	if rErr != nil {
		// log the error
		log.Println("this is the error: ", rErr)
		// check for the possible errors
		if rErr.Error() == "user already exists" {
			c.JSON(http.StatusConflict, utils.Response(http.StatusConflict, nil, rErr.Error()))
			return
		}
		if rErr.Error() == "request timed out" {
			c.JSON(408, utils.Response(http.StatusRequestTimeout, nil, rErr.Error()))
			return
		}
		c.JSON(500, utils.Response(http.StatusInternalServerError, nil, rErr.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.Response(http.StatusCreated, nil, resp))
}

func (us *UserController) Login(c *gin.Context) {
	var pl dto.LoginDTO

	if err := c.ShouldBindJSON(&pl); err != nil {
		// log the error
		log.Println(err)
		c.JSON(400, utils.Response(http.StatusBadRequest, nil, err.Error()))
		return
	}
	// send to service
	resp, rErr := us.userService.Login(pl)
	if rErr != nil {
		// log the error
		log.Println("this is the error: ", rErr)
		// check for the possible errors
		if rErr.Error() == "invalid credentials" {
			c.JSON(http.StatusBadRequest, utils.Response(http.StatusBadRequest, nil, rErr.Error()))
			return
		}
		if rErr.Error() == "request timed out" {
			c.JSON(408, utils.Response(http.StatusRequestTimeout, nil, rErr.Error()))
			return
		}
		c.JSON(500, utils.Response(http.StatusInternalServerError, nil, rErr.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.Response(http.StatusOK, resp, "Login Successful"))
}
