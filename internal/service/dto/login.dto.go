package dto

type LoginDTO struct {
	Email    string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
