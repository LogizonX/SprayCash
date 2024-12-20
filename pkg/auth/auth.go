package auth

import (
	"log"
	"strconv"
	"time"

	"github.com/LoginX/SprayDash/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

var UserKey contextKey = "userEmail"

func CreateJWT(secret []byte, expiration int, user *model.User) (map[string]string, error) {
	expAt := time.Now().Add(time.Duration(expiration) * time.Second).Unix()
	log.Println("expAt: ", expAt)
	// create te token claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userEmail": user.Email,
		"exp":       expAt,
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}
	data := map[string]string{
		"token":     tokenString,
		"expiresAt": strconv.FormatInt(expAt, 10),
	}
	return data, nil
}
