package service

import (
	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/service/dto"
)

// user service interface
type UserService interface {
	Login(loginDto dto.LoginDTO) (dto.LoginResponseDTO, error)
	Register(createUserDto dto.CreateUserDTO) (string, error)
	GetUserDetails(email string) (*model.User, error)
	LoginSocial(loginDto dto.LoginSocialDTO) (dto.LoginResponseDTO, error)
}
