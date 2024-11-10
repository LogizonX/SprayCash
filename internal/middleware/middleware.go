package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/LoginX/SprayDash/config"

	"github.com/LoginX/SprayDash/internal/service"
	"github.com/LoginX/SprayDash/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token for the context
		tokenString, err := utils.GetTokenFromRequest(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Response(http.StatusUnauthorized, nil, err.Error()))
			return
		}
		// decode the token
		token, pErr := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			// check the signing of ok
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			// return the signing secret
			return []byte(config.GetEnv("JWT_SECRET", "somesecret")), nil
		})
		if pErr != nil {
			log.Println(pErr.Error())
			utils.Forbidden(c)
			return
		}
		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			log.Println("got here for ok")
			utils.Forbidden(c)
			return
		}
		// get the userEmail from the token claim
		userEmail, ok := (*claims)["userEmail"].(string)
		if !ok {
			log.Println("got here for userEmail")
			utils.Forbidden(c)
			return
		}
		// get the user from the database
		user, err := userService.GetUserDetails(userEmail)
		if err != nil {
			log.Println("got here for user")
			utils.Forbidden(c)
			return
		}

		// set the user in the context
		c.Set("user", user)
		c.Next()

	}
}
