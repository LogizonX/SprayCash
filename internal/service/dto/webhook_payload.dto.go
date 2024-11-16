package dto

type ReceivedFrom struct {
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
}

type Transaction struct {
	TransactionReference string       `json:"transaction_reference"`
	TransactionStatus    string       `json:"transaction_status"`
	VirtualAccountNumber string       `json:"virtual_account_number"`
	TransactionFee       float64      `json:"transaction_fee"`
	AmountReceived       float64      `json:"amount_received"`
	InitiatedDate        string       `json:"initiated_date"`
	CurrentStatusDate    string       `json:"current_status_date"`
	ReceivedFrom         ReceivedFrom `json:"received_from"`
	Channel              string       `json:"channel"`
	CurrencyCode         string       `json:"currency_code"`
	Branch               bool         `json:"branch"`
	SessionID            string       `json:"session_id"`
	Status               string       `json:"status"`
}

type Duration struct{}

type TransactionData struct {
	TransactionReference         string  `json:"transaction_reference"`
	AmountReceived               float64 `json:"amount_received"`
	TransactionFee               float64 `json:"transaction_fee"`
	TransactionStatus            string  `json:"transaction_status"`
	SenderName                   string  `json:"sender_name"`
	SenderAccountNumber          string  `json:"sender_account_number"`
	SourceBankName               *string `json:"source_bank_name"`
	InitiatedDate                string  `json:"initiated_date"`
	CurrentStatusDate            string  `json:"current_status_date"`
	Currency                     string  `json:"currency"`
	SessionID                    string  `json:"session_id"`
	MerchantTransactionReference string  `json:"merchant_transaction_reference"`
	TransactionType              string  `json:"transaction_type"`
	VirtualAccountNumber         string  `json:"virtual_account_number"`
	StatusReason                 string  `json:"status_reason"`
}

type TransactionResponse struct {
	Message string          `json:"message"`
	Success bool            `json:"success"`
	Data    TransactionData `json:"data"`
}

type TestFundDTO struct {
	Amount        float64 `json:"amount"`
	AccountNumber string  `json:"account_number"`
}

type TestFundAccountResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type FundAccountPayazaResponse struct {
	ResponseCode         int    `json:"response_code"`
	ResponseMessage      string `json:"response_message"`
	TransactionReference string `json:"transactionreference"`
}
