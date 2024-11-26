package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/service/dto"
	"github.com/LoginX/SprayDash/internal/service/impls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	// dependencies setup
	mockRepo := &MockUserRepository{}
	mockCache := &MockRedisCacheService{}
	mockMailer := &mailerService{}
	mockCodeGenerator := &codeGeneratorService{}
	userService := impls.NewUserServiceImpl(mockRepo, mockCache, mockMailer, mockCodeGenerator)
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
		mockCache.On("Set", mock.Anything, createUserDto.Email, mock.AnythingOfType("int"), mock.AnythingOfType("time.Duration")).Return(1234, nil)
		mockCodeGenerator.On("GenerateCode").Return(1234)
		mockMailer.On("SendMail", createUserDto.Email, "Welcome to SprayDash", createUserDto.Name, "Your account has been created successfully.", "welcome_email").Return(nil)

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

type MockRedisCacheService struct {
	mock.Mock
}

func (m *MockRedisCacheService) Set(ctx context.Context, email string, code int, expiration time.Duration) error {
	return nil
}

func (m *MockRedisCacheService) Get(ctx context.Context, email string) (int, error) {
	return 0, nil
}

type mailerService struct {
	mock.Mock
}

func (m *mailerService) SendMail(recipient string, subject string, username string, message string, template_name string) error {
	return nil
}

type codeGeneratorService struct {
	mock.Mock
}

func (c *codeGeneratorService) GenerateCode() int {
	return 1234
}

func (c *codeGeneratorService) GenerateInviteCode() string {
	return "1234"
}

func (c *codeGeneratorService) GenerateReferenceCode() string {
	return "1234"
}
