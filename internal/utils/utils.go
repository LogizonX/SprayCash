package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/LoginX/SprayDash/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gopkg.in/gomail.v2"
)

type RedisCacheService struct {
	client *redis.Client
}

func NewRedisCacheService() *RedisCacheService {
	client := redis.NewClient(&redis.Options{
		Addr:      config.GetEnv("REDIS_URL", "localhost:6379"),
		DB:        0,
		Password:  config.GetEnv("REDIS_ACCESS", ""),
		TLSConfig: &tls.Config{},
	})
	return &RedisCacheService{client: client}
}

func (r *RedisCacheService) Set(ctx context.Context, email string, code int, expiration time.Duration) error {
	return r.client.Set(ctx, email, code, expiration).Err()
}

func (r *RedisCacheService) Get(ctx context.Context, email string) (int, error) {
	code, err := r.client.Get(ctx, email).Int()
	if err != nil {
		return 0, err
	}
	return code, nil
}

type MailerService struct{}

func NewMailerService() *MailerService {
	return &MailerService{}
}

func (ms *MailerService) SendMail(recipient string, subject string, username string, message string, template_name string) error {
	templ, err := os.ReadFile(fmt.Sprintf("internal/utils/templates/%s.html", template_name))
	if err != nil {
		log.Println("Error reading email template:", err)
		return err
	}
	t, err := template.New("email").Parse(string(templ))
	if err != nil {
		log.Println("Error parsing email template:", err)
		return err
	}
	// mail data
	mailData := map[string]interface{}{
		"Username": username,
		"Message":  message,
	}
	// write the maildata to bytes buffer
	buf := new(bytes.Buffer)
	err = t.Execute(buf, mailData)
	if err != nil {
		log.Println("Error executing email template:", err)
		return err
	}
	m := gomail.NewMessage()

	// Set email headers
	m.SetHeader("From", "sainthaywon80@gmail.com")
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)

	// Set the HTML body
	m.SetBody("text/html", buf.String())
	smtpHost := config.GetEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := 465
	smtpUser := config.GetEnv("SMTP_USER", "protected@gmail.com")
	smtpPass := config.GetEnv("SMTP_PWD", "protected")

	// Create a new SMTP dialer
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Send the email and handle errors
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	// Success message
	fmt.Println("Email sent successfully!")

	return nil

}

type CodeGeneratorService struct{}

func NewCodeGeneratorService() *CodeGeneratorService {
	return &CodeGeneratorService{}
}

func (cg *CodeGeneratorService) GenerateCode() int {
	// get random four digit code
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(9000) + 1000
	return code
}

func (cg *CodeGeneratorService) GenerateInviteCode() string {
	// use string and timestamp to generate randon invite code
	rand.Seed(time.Now().UnixNano())
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	ranNumber := rand.Intn(100)
	return timestamp + strconv.Itoa(ranNumber)
}

func (cg *CodeGeneratorService) GenerateReferenceCode() string {
	// generate randm aplhanumeric code
	rand.Seed(time.Now().UnixNano())
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateAndCacheCode(cacheService CacheService, codeGenerator CodeGenerator, email string) (int, error) {
	// get random four digit code
	code := codeGenerator.GenerateCode()
	// cache the code
	ctx := context.Background()
	err := cacheService.Set(ctx, email, code, 15*time.Minute)
	if err != nil {
		log.Println("Error caching code:", err)
		return 0, err
	}
	return code, nil

}

func GetCachedCode(cacheService CacheService, email string) (int, error) {
	ctx := context.Background()
	code, err := cacheService.Get(ctx, email)
	if err != nil {
		log.Println("Error getting cached code:", err)
		return 0, err
	}
	return code, nil
}

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
