package utils

import (
	"context"
	"time"
)

type CacheService interface {
	Set(ctx context.Context, email string, code int, expiration time.Duration) error
	Get(ctx context.Context, email string) (int, error)
}

type Mailer interface {
	SendMail(to, subject, userName, message, templateName string) error
}

type CodeGenerator interface {
	GenerateCode() int
	GenerateInviteCode() string
	GenerateReferenceCode() string
}

type AzureService interface {
	UploadFileToAzureBlob(file []byte, fileName string, containerName string) (string, error)
}
