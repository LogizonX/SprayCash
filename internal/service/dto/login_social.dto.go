package dto

type LoginSocialDTO struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email" binding: "required"`
}
