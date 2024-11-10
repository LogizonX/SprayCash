package impls

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/LoginX/SprayDash/config"
	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/repository"
	"github.com/LoginX/SprayDash/internal/service/dto"
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
	_, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		log.Println("Error creating user: ", err)
		if errors.Is(err, context.DeadlineExceeded) {
			return "", errors.New("request timed out")
		}
		return "", err
	}
	return "User registered successfully", nil

}

func (s *UserServiceImpl) Login(loginDto dto.LoginDTO) (dto.LoginResponseDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// get user by the email
	user, err := s.repo.GetUserByEmail(ctx, loginDto.Email)
	fmt.Println(loginDto.Email)
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
