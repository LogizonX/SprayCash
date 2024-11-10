package utils

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Response(statusCode int, data any, message any) map[string]any {
	var status string
	switch {
	case statusCode >= 200 && statusCode <= 299:
		status = "success"
	case statusCode == 400:
		status = "error"
	case statusCode >= 300 && statusCode <= 399:
		status = "redirect"
	case statusCode == 404:
		status = "not found"
	case statusCode >= 405 && statusCode <= 499:
		status = "error"
	case statusCode == 401 || statusCode == 403:
		status = "unauthorized"
	case statusCode >= 500:
		status = "error"
		message = "This is from us!, please contact admin"
	default:
		status = "error"
		message = "This is from us!, please contact admin"
	}
	res := map[string]any{
		"status":      status,
		"data":        data,
		"message":     message,
		"status_code": statusCode,
	}
	return res

}

func GetTokenFromRequest(c *gin.Context) (string, error) {
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("authorization header not provided")
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		return "", errors.New("invalid authorization header format")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != "bearer" {
		return "", errors.New("unsupported authorization type")
	}

	accessToken := fields[1]
	return accessToken, nil
}

func GetUserFromContext(c *gin.Context) (any, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response(http.StatusForbidden, nil, "Unauthorized"))
	c.Abort()
}
