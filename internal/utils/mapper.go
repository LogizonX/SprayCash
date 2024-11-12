package utils

type UserResponse struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	WalletBalance float64 `json:"wallet_balance"`
}

func NewUserResponse(id string, name string, email string, wallet_balance float64) *UserResponse {
	return &UserResponse{
		Id:            id,
		Name:          name,
		Email:         email,
		WalletBalance: wallet_balance,
	}
}
