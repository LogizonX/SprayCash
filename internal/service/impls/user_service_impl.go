package impls

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/LoginX/SprayDash/config"
	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/repository"
	"github.com/LoginX/SprayDash/internal/service/dto"
	"github.com/LoginX/SprayDash/internal/utils"
	"github.com/LoginX/SprayDash/pkg/auth"
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
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		url := "https://router.prod.payaza.africa/api/request/secure/payloadhandler"
		userFirstName := strings.Fields(user.Name)[0]
		userLastName := strings.Fields(user.Name)[1]
		payload := map[string]interface{}{
			"service_type": "Account",
			"service_payload": map[string]interface{}{
				"request_application":      "Payaza",
				"application_module":       "USER_MODULE",
				"application_version":      "1.0.0",
				"request_class":            "CreateReservedAccountForCustomers",
				"customer_first_name":      userFirstName,
				"customer_last_name":       userLastName,
				"customer_email":           user.Email,
				"customer_phone":           "09012345673",
				"virtual_account_provider": "Premiumtrust",
			},
		}
		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Authorization", "Payaza UFo3OC1QS1RFU1QtMjM5NTNCMUUtQ0RGQS00OTFCLUIxMEMtN0IzRjEyMzlBRjI5")
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()
		print(resp.Body)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		respBody := []byte(string(body))
		var responseBody utils.ResponseBody
		jErr := json.Unmarshal(respBody, &responseBody)
		if jErr != nil {
			fmt.Println("Error unmarshalling JSON:", jErr)
			return
		}
		accountDetails := model.NewAccountDetails(responseBody.ResponseContent.VirtualAccountName, responseBody.ResponseContent.VirtualAccountNumber, responseBody.ResponseContent.VirtualProviderBankName, responseBody.ResponseContent.VirtualProviderBankCode)
		// update bankdetails
		err = s.repo.UpdateUserBankDetails(ctx, user.Email, accountDetails)
		if err != nil {
			fmt.Println("Error updating bank details:", err)
			return
		}
		fmt.Println("Response:", resp.Status)
		fmt.Println("Response Body:", string(body))

	}()
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
	return dto.LoginResponseDTO{
		AccessToken: token["token"],
		ExpiresIn:   token["expiresAt"],
		TokenType:   "jwt",
	}, nil

}
