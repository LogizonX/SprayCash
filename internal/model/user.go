package model

type User struct {
	Id             string         `bson:"_id,omitempty" json:"id"`
	Name           string         `bson:"name" json:"name"`
	Email          string         `bson:"email" json:"email"`
	Password       string         `bson:"password"  json:"-"`
	WalletBalance  float64        `bson:"wallet_balance" json:"wallet_balance"`
	AccountDetails AccountDetails `bson:"account_details" json:"account_details"`
}

func NewUser(name string, email string, password string) *User {
	return &User{
		Name:          name,
		Email:         email,
		Password:      password,
		WalletBalance: 0.0,
	}
}

type AccountDetails struct {
	AccountName string `bson:"account_name" json:"account_name"`
	AccountNo   string `bson:"account_no" json:"account_no"`
	BankName    string `bson:"bank_name" json:"bank_name"`
	BankCode    string `bson:"bank_code" json:"bank_code"`
}

func NewAccountDetails(accountName string, accountNo string, bankName string, bankCode string) *AccountDetails {
	return &AccountDetails{
		AccountName: accountName,
		AccountNo:   accountNo,
		BankName:    bankName,
		BankCode:    bankCode,
	}
}
