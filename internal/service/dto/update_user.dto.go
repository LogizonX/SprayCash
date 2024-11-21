package dto

type UpdateUserDTO struct {
	Verified bool   `json:"verified"`
	Name     string `json:"name"`
}
