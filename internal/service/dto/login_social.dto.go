package dto

type LoginSocialDTO struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email" binding: "required"`
	Username string `json:"username" binding: "required"`
}
