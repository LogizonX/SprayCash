package dto

type VerifyUserDTO struct {
	Email string `json:"email" binding:"required"`
	Code  int    `json:"code" binding:"required"`
}
