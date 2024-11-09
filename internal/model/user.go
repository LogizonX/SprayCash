package model

type User struct {
	Id            string `bson:"_id,omitempty" json:"id"`
	Name          string `bson:"name" json:"name"`
	Email         string `bson:"email" json:"email"`
	Password      string `bson:"password"  json:"password"`
	WalletBalance string `bson:"wallet_balance" json:"wallet_balance"`
}
