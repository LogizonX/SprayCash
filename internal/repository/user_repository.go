package repository

import (
	"context"

	"github.com/LoginX/SprayDash/internal/model"
)

type UserRepository interface {
	// GetUserByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreditUser(ctx context.Context, amount float64, userEmail string, walletHistory *model.WalletHistory) error
	DebitUser(ctx context.Context, amount float64, userEmail string, walletHistory *model.WalletHistory) error
	UpdateUserBankDetails(ctx context.Context, userEmail string, accountDetails *model.AccountDetails) error
	GetUserByVirtualAccount(ctx context.Context, virtualAccount string) (*model.User, error)
	CreateWalletHistory(ctx context.Context, walletHistory *model.WalletHistory) (*model.WalletHistory, error)
	CreateNewFundsTracking(ctx context.Context, fundsTracking *model.FundsTracking) (*model.FundsTracking, error)
	UpdateUser(ctx context.Context, updateMap map[string]interface{}, email string) (*model.User, error)
}
