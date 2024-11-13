package repository

import (
	"context"

	"github.com/LoginX/SprayDash/internal/model"
)

type UserRepository interface {
	// GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreditUser(ctx context.Context, amount float64, userId string) error
	DebitUser(ctx context.Context, amount float64, userId string) error
	UpdateUserBankDetails(ctx context.Context, userId string, accountDetails *model.AccountDetails) error
}
