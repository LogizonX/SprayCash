package controller

import (
	"log"
	"net/http"

	"github.com/LoginX/SprayDash/internal/middleware"
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
	rg.POST("/auth/social-auth", uc.LoginSocial)
	rg.GET("/users", middleware.AuthMiddleware(uc.userService), uc.FetchUserDetails)
	rg.POST("/payaza/webhook", uc.PayazaWebhook)
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

// FetchUserDetails from the authmiddleware
func (uc *UserController) FetchUserDetails(c *gin.Context) {
	// fetch the user details from the context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		log.Println("this is the error: ", err)
		c.JSON(404, utils.Response(404, nil, "user not found"))
	}
	c.JSON(200, utils.Response(200, user, "User details fetched successfully"))

}

func (uc *UserController) LoginSocial(c *gin.Context) {
	// bind payload
	var pl dto.LoginSocialDTO

	if err := c.ShouldBindJSON(&pl); err != nil {
		// log the error
		log.Println(err)
		c.JSON(400, utils.Response(http.StatusBadRequest, nil, err.Error()))
		return
	}
	// send to service
	resp, rErr := uc.userService.LoginSocial(pl)
	if rErr != nil {
		// log the error
		log.Println("this is the error: ", rErr)
		if rErr.Error() == "request timed out" {
			c.JSON(408, utils.Response(http.StatusRequestTimeout, nil, rErr.Error()))
			return
		}
		c.JSON(500, utils.Response(http.StatusInternalServerError, nil, rErr.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.Response(http.StatusOK, resp, "Login Successful"))
}

func (uc *UserController) GenerateDynamicAccount(c *gin.Context) {

}

func (uc *UserController) PayazaWebhook(c *gin.Context) {
	pl := new(dto.Transaction)
	if err := c.ShouldBindJSON(&pl); err != nil {
		log.Println(err)
		c.JSON(400, utils.Response(http.StatusBadRequest, nil, err.Error()))
		return
	}
	msg, rErr := uc.userService.PayazaWebhook(pl)
	if rErr != nil {
		log.Println("this is the error: ", rErr)
		if rErr.Error() == "request timed out" {
			c.JSON(408, utils.Response(http.StatusRequestTimeout, nil, rErr.Error()))
			return
		}
		c.JSON(500, utils.Response(http.StatusInternalServerError, nil, rErr.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.Response(http.StatusOK, nil, msg))

}
