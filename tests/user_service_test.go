package test

import (
	"context"
	"errors"
	"testing"

	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/service/dto"
	"github.com/LoginX/SprayDash/internal/service/impls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	// dependencies setup
	mockRepo := &MockUserRepository{}
	userService := impls.NewUserServiceImpl(mockRepo)
	mockUtils := &MockUtils{}
	// arrange
	createUserDto := dto.CreateUserDTO{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword := "hashedPassword"
	newUser := model.NewUser(createUserDto.Name, createUserDto.Email, hashedPassword)

	t.Run("successful registration", func(t *testing.T) {
		// expect
		mockRepo.On("GetUserByEmail", mock.Anything, createUserDto.Email).Return(newUser, errors.New("user not found")).Once()
		mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).Return(newUser, nil)
		mockUtils.On("GenerateAndCacheCode", createUserDto.Email).Return(1234, nil)

		// act
		message, err := userService.Register(createUserDto)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, "User registered successfully", message)
		assert.Equal(t, newUser.Name, createUserDto.Name)
		mockRepo.AssertExpectations(t)
	})

}

type MockUserRepository struct {
	mock.Mock
}



func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*model.User), args.Error(1)

}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*model.User), args.Error(1)

}

func (m *MockUserRepository) CreditUser(ctx context.Context, amount float64, userEmail string, walletHistory *model.WalletHistory) error {
	return nil
}

func (m *MockUserRepository) DebitUser(ctx context.Context, amount float64, userEmail string, walletHistory *model.WalletHistory) error {
	return nil
}

func (m *MockUserRepository) UpdateUserBankDetails(ctx context.Context, userEmail string, accountDetails *model.AccountDetails) error {
	return nil
}

func (m *MockUserRepository) GetUserByVirtualAccount(ctx context.Context, virtualAccount string) (*model.User, error) {
	return nil, nil
}

func (m *MockUserRepository) CreateWalletHistory(ctx context.Context, walletHistory *model.WalletHistory) (*model.WalletHistory, error) {
	return nil, nil
}

func (m *MockUserRepository) CreateNewFundsTracking(ctx context.Context, fundsTracking *model.FundsTracking) (*model.FundsTracking, error) {
	return nil, nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, updateMap map[string]interface{}, email string) (*model.User, error) {
	return nil, nil
}

type MockUtils struct {
	mock.Mock
}

func (m *MockUtils) GenerateAndCacheCode(email string) (int, error) {
	args := m.Called(email)
	return args.Int(0), args.Error(1)
}
