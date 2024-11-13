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

type ResponseContent struct {
	VirtualAccountName      string `json:"virtual_account_name"`
	VirtualProviderBankName string `json:"virtual_provider_bank_name"`
	VirtualProviderBankCode string `json:"virtual_provider_bank_code"`
	VirtualAccountNumber    string `json:"virtual_account_number"`
}

type ResponseBody struct {
	ResponseContent ResponseContent `json:"response_content"`
	ResponseCode    int             `json:"response_code"`
	ResponseMessage string          `json:"response_message"`
}
