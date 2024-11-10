package model

type User struct {
	Id            string  `bson:"_id,omitempty" json:"id"`
	Name          string  `bson:"name" json:"name"`
	Email         string  `bson:"email" json:"email"`
	Password      string  `bson:"password"  json:"-"`
	WalletBalance float64 `bson:"wallet_balance" json:"wallet_balance"`
}

func NewUser(name string, email string, password string) *User {
	return &User{
		Name:          name,
		Email:         email,
		Password:      password,
		WalletBalance: 0.0,
	}
}
