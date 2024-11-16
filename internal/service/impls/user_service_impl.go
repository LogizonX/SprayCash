package impls

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/LoginX/SprayDash/config"
	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/repository"
	"github.com/LoginX/SprayDash/internal/service/dto"
	"github.com/LoginX/SprayDash/internal/utils"
	"github.com/LoginX/SprayDash/pkg/auth"
	"github.com/LoginX/SprayDash/pkg/common"
)

type UserServiceImpl struct {
	// depends on
	repo repository.UserRepository
}

func NewUserServiceImpl(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (s *UserServiceImpl) generateVirtualAccount(user *model.User) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	responseBody, err := common.GeneratePayazaVirtualAccount(user)
	if err != nil {
		fmt.Println("Error generating virtual account:", err)
		return
	}
	fmt.Println("Response body: ", responseBody)
	accountDetails := model.NewAccountDetails(responseBody.ResponseContent.VirtualAccountName, responseBody.ResponseContent.VirtualAccountNumber, responseBody.ResponseContent.VirtualProviderBankName, responseBody.ResponseContent.VirtualProviderBankCode)
	// update bankdetails
	err = s.repo.UpdateUserBankDetails(ctx, user.Email, accountDetails)
	if err != nil {
		fmt.Println("Error updating bank details:", err)
		return
	}

}

// implement interface methods

func (s *UserServiceImpl) Register(createUserDto dto.CreateUserDTO) (string, error) {
	// need to hash the password
	hashedPassword, hashErr := auth.HashPassword(createUserDto.Password)
	if hashErr != nil {
		log.Println("Error hashing password: ", hashErr)
		return "", hashErr
	}
	// check if user already exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, existErr := s.repo.GetUserByEmail(ctx, createUserDto.Email)
	if existErr == nil {
		log.Println("User already exists: ", existErr)
		return "", errors.New("user already exists")
	}

	newUser := model.NewUser(createUserDto.Name, createUserDto.Email, hashedPassword)

	// Call the CreateUser function with the context
	user, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		log.Println("Error creating user: ", err)
		if errors.Is(err, context.DeadlineExceeded) {
			return "", errors.New("request timed out")
		}
		return "", err
	}
	// get the bank details in a goroutine
	go s.generateVirtualAccount(user)
	// send a welcome email
	code, cErr := utils.GenerateAndCacheCode()
	if cErr != nil {
		log.Println("Error generating code: ", cErr)
	} else {

		// send email
		go utils.SendMail(user.Email, "Welcome to SprayDash", user.Name, fmt.Sprintf("%s", code))
	}

	return "User registered successfully", nil

}

func (s *UserServiceImpl) Login(loginDto dto.LoginDTO) (dto.LoginResponseDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// get user by the email
	user, err := s.repo.GetUserByEmail(ctx, loginDto.Email)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return dto.LoginResponseDTO{}, errors.New("request timed out")
		}
		// TODO: handle case of user does not exists
		return dto.LoginResponseDTO{}, err
	}
	// compare password
	if !auth.ComparePassword(user.Password, loginDto.Password) {
		return dto.LoginResponseDTO{}, errors.New("invalid credentials")
	}

	secret := []byte(config.GetEnv("JWT_SECRET", "somesecret"))
	exp, expErr := strconv.Atoi(config.GetEnv("JWT_EXP", "3600"))
	if expErr != nil {
		log.Println("Error converting JWT_EXP to int: ", expErr)
		return dto.LoginResponseDTO{}, expErr
	}
	token, tokenErr := auth.CreateJWT(secret, exp, user)
	if tokenErr != nil {
		log.Println("Error creating JWT: ", tokenErr)
		return dto.LoginResponseDTO{}, errors.New("error creating a token")
	}
	if user.AccountDetails.AccountNo == "" {
		go s.generateVirtualAccount(user)
	}
	// return token
	return dto.LoginResponseDTO{
		AccessToken: token["token"],
		ExpiresIn:   token["expiresAt"],
		TokenType:   "jwt",
	}, nil

}

func (s *UserServiceImpl) GetUserDetails(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserServiceImpl) LoginSocial(pl dto.LoginSocialDTO) (dto.LoginResponseDTO, error) {
	email := pl.Email
	// get user by the email
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.repo.GetUserByEmail(ctx, email)
	// user does not exists create one
	if err != nil {
		// create user
		hashedPassword, hashErr := auth.HashPassword(pl.Email)
		if hashErr != nil {
			log.Println("Error hashing password: ", hashErr)
			return dto.LoginResponseDTO{}, hashErr
		}
		newUser := model.NewUser(pl.Name, pl.Email, hashedPassword)
		_, err := s.repo.CreateUser(ctx, newUser)
		if err != nil {
			log.Println("Error creating user: ", err)
			if errors.Is(err, context.DeadlineExceeded) {
				return dto.LoginResponseDTO{}, errors.New("request timed out")
			}
			return dto.LoginResponseDTO{}, err
		}
	}
	secret := []byte(config.GetEnv("JWT_SECRET", "somesecret"))
	exp, expErr := strconv.Atoi(config.GetEnv("JWT_EXP", "3600"))
	if expErr != nil {
		log.Println("Error converting JWT_EXP to int: ", expErr)
		return dto.LoginResponseDTO{}, expErr
	}
	token, tokenErr := auth.CreateJWT(secret, exp, user)
	if tokenErr != nil {
		log.Println("Error creating JWT: ", tokenErr)
		return dto.LoginResponseDTO{}, errors.New("error creating a token")
	}
	if user.AccountDetails.AccountNo == "" {
		go s.generateVirtualAccount(user)
	}
	return dto.LoginResponseDTO{
		AccessToken: token["token"],
		ExpiresIn:   token["expiresAt"],
		TokenType:   "jwt",
	}, nil

}

func (s *UserServiceImpl) PayazaWebhook(pl *dto.Transaction) (string, error) {
	// get the user by the email
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// get the transaction reference and query transaction status
	tranRef := pl.TransactionReference
	url := fmt.Sprintf("https://api.payaza.africa/live/payaza-account/api/v1/mainaccounts/merchant/transaction/%s", tranRef)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Payaza %s", config.GetEnv("PAYAZA_API_KEY", "somkey")))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()
	res, rErr := ioutil.ReadAll(resp.Body)
	if rErr != nil {
		log.Println("Error reading response body:", rErr)
		return "", rErr
	}
	// convert the response to a byte runes
	resBytes := []byte(string(res))
	responseBody := new(dto.TransactionResponse)
	// unmarshal the response
	err = json.Unmarshal(resBytes, responseBody)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return "", err
	}
	// check the status of the transaction
	if responseBody.Data.TransactionStatus == "NIP_SUCCESS" && responseBody.Status == true {
		// get user by the virtual account number
		virtualAccountNumber := pl.VirtualAccountNumber
		user, err := s.repo.GetUserByVirtualAccount(ctx, virtualAccountNumber)
		if err != nil {
			log.Println("Error getting user by virtual account:", err)
			return "", err
		}

		newWalletHistory := model.NewWalletHistory(user.Email, float64(responseBody.Data.TransactionAmount), user.WalletBalance, user.WalletBalance+float64(responseBody.Data.TransactionAmount), responseBody.Data.TransactionReference)

		// credit the user account
		creditErr := s.repo.CreditUser(ctx, float64(responseBody.Data.TransactionAmount), user.Email, newWalletHistory)
		if creditErr != nil {
			log.Println("Error crediting user account:", creditErr)
			return "", creditErr
		}
		// send a mail to the user
		// TODO: change email template for this
		go utils.SendMail(user.Email, "Fund Successful", user.Name, "Your transaction was successful")
		return "success", nil
	}
	return "failed", errors.New("transaction failed")

}
