package dto

type LoginResponseDTO struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   string `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}
