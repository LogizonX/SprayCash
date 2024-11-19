package model

import "time"

type User struct {
	Id             string         `bson:"_id,omitempty" json:"id"`
	Name           string         `bson:"name" json:"name"`
	Email          string         `bson:"email" json:"email"`
	Password       string         `bson:"password"  json:"-"`
	WalletBalance  float64        `bson:"wallet_balance" json:"wallet_balance"`
	Verified       bool           `bson:"verified" json:"verified"`
	AccountDetails AccountDetails `bson:"account_details" json:"account_details"`
}

type WalletHistory struct {
	Id             string    `bson:"_id,omitempty" json:"id"`
	UserId         string    `bson:"user_id" json:"user_id"`
	Amount         float64   `bson:"amount" json:"amount"`
	PreviousAmount float64   `bson:"previous_amount" json:"previous_amount"`
	AfterAmount    float64   `bson:"after_amount" json:"after_amount"`
	TransactionRef string    `bson:"transaction_ref" json:"transaction_ref"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
}

type FundsTracking struct {
	Id             string    `bson:"_id,omitempty" json:"id"`
	ReceiverEmail  string    `bson:"receiver_email" json:"receiver_email"`
	ReceiverName   string    `bson:"receiver_name" json:"receiver_name"`
	SenderEmail    string    `bson:"sender_email" json:"sender_email"`
	SenderName     string    `bson:"sender_name" json:"sender_name"`
	Amount         float64   `bson:"amount" json:"amount"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	Status         string    `bson:"status" json:"status"`
	TransactionRef string    `bson:"transaction_ref" json:"transaction_ref"`
	InviteCode     string    `bson:"invite_code" json:"invite_code"`
}

func NewFundsTracking(receiverEmail string, receiverName string, senderEmail string, senderName string, amount float64, status string, transactionRef string, inviteCode string) *FundsTracking {
	return &FundsTracking{
		ReceiverEmail:  receiverEmail,
		ReceiverName:   receiverName,
		SenderEmail:    senderEmail,
		SenderName:     senderName,
		Amount:         amount,
		CreatedAt:      time.Now(),
		Status:         status,
		TransactionRef: transactionRef,
		InviteCode:     inviteCode,
	}
}

func NewWalletHistory(userId string, amount float64, previousAmount float64, afterAmount float64, transactionRef string) *WalletHistory {
	return &WalletHistory{
		UserId:         userId,
		Amount:         amount,
		PreviousAmount: previousAmount,
		AfterAmount:    afterAmount,
		TransactionRef: transactionRef,
		CreatedAt:      time.Now(),
	}
}

func NewUser(name string, email string, password string) *User {
	return &User{
		Name:          name,
		Email:         email,
		Password:      password,
		WalletBalance: 0.0,
		Verified:      false,
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
