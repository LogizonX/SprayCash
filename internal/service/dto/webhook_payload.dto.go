package dto

type ReceivedFrom struct {
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
}

type Transaction struct {
	TransactionReference string       `json:"transaction_reference"`
	TransactionStatus    string       `json:"transaction_status"`
	VirtualAccountNumber string       `json:"virtual_account_number"`
	TransactionFee       int          `json:"transaction_fee"`
	AmountReceived       int          `json:"amount_received"`
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
	TransactionDateTime  string   `json:"transactionDateTime"`
	TransactionReference string   `json:"transactionReference"`
	CreditAccount        string   `json:"creditAccount"`
	BankCode             string   `json:"bankCode"`
	BeneficiaryName      string   `json:"beneficiaryName"`
	TransactionAmount    int      `json:"transactionAmount"`
	Fee                  int      `json:"fee"`
	SessionID            string   `json:"sessionId"`
	TransactionStatus    string   `json:"transactionStatus"`
	Narration            string   `json:"narration"`
	TransactionType      string   `json:"transactionType"`
	ResponseMessage      string   `json:"responseMessage"`
	ResponseCode         string   `json:"responseCode"`
	Currency             string   `json:"currency"`
	BalanceBefore        float64  `json:"balanceBefore"`
	BalanceAfter         float64  `json:"balanceAfter"`
	Duration             Duration `json:"duration"`
}

type TransactionResponse struct {
	Message    string          `json:"message"`
	Status     bool            `json:"status"`
	RetryCount int             `json:"retry_count"`
	Data       TransactionData `json:"data"`
}
