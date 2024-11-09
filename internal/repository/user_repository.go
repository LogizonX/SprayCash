package repository

import (
	"context"

	"github.com/LoginX/SprayDash/internal/model"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
}
