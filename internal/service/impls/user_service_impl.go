package impls

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/LoginX/SprayDash/config"
	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/repository"
	"github.com/LoginX/SprayDash/internal/service/dto"
	"github.com/LoginX/SprayDash/internal/utils"
	"github.com/LoginX/SprayDash/pkg/auth"
	"github.com/LoginX/SprayDash/pkg/common"
)

type UserServiceImpl struct {
	// depends on
	repo repository.UserRepository
}

func NewUserServiceImpl(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (s *UserServiceImpl) generateVirtualAccount(user *model.User) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	responseBody, err := common.GeneratePayazaVirtualAccount(user)
	if err != nil {
		fmt.Println("Error generating virtual account:", err)
		return
	}
	fmt.Println("Response body: ", responseBody)
	accountDetails := model.NewAccountDetails(responseBody.ResponseContent.VirtualAccountName, responseBody.ResponseContent.VirtualAccountNumber, responseBody.ResponseContent.VirtualProviderBankName, responseBody.ResponseContent.VirtualProviderBankCode)
	// update bankdetails
	err = s.repo.UpdateUserBankDetails(ctx, user.Email, accountDetails)
	if err != nil {
		fmt.Println("Error updating bank details:", err)
		return
	}

}

// implement interface methods

func (s *UserServiceImpl) Register(createUserDto dto.CreateUserDTO) (string, error) {
	// need to hash the password
	hashedPassword, hashErr := auth.HashPassword(createUserDto.Password)
	if hashErr != nil {
		log.Println("Error hashing password: ", hashErr)
		return "", hashErr
	}
	// check if user already exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, existErr := s.repo.GetUserByEmail(ctx, createUserDto.Email)
	if existErr == nil {
		log.Println("User already exists: ", existErr)
		return "", errors.New("user already exists")
	}

	newUser := model.NewUser(createUserDto.Name, createUserDto.Email, hashedPassword)

	// Call the CreateUser function with the context
	user, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		log.Println("Error creating user: ", err)
		if errors.Is(err, context.DeadlineExceeded) {
			return "", errors.New("request timed out")
		}
		return "", err
	}
	// get the bank details in a goroutine
	go s.generateVirtualAccount(user)
	// send a welcome email
	code, cErr := utils.GenerateAndCacheCode(newUser.Email)
	if cErr != nil {
		log.Println("Error generating code: ", cErr)
	} else {

		// send email
		go utils.SendMail(user.Email, "Welcome to SprayDash", user.Name, fmt.Sprintf("%d", code), "email_template")
	}

	return "User registered successfully", nil

}

func (s *UserServiceImpl) VerifyUser(pl dto.VerifyUserDTO) (string, error) {
	// get the user by the email
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.repo.GetUserByEmail(ctx, pl.Email)
	if err != nil {
		log.Println("Error getting user by email:", err)
		return "error", errors.New("user not found")
	}
	// get code
	code, cErr := utils.GetCachedCode(pl.Email)
	if cErr != nil {
		log.Println("Error getting cached code:", cErr)
		return "error", cErr
	}
	if pl.Code != code {
		return "error", errors.New("invalid code")

	}

	// update the user account
	updateMap := map[string]interface{}{
		"verified": true,
	}
	_, uErr := s.repo.UpdateUser(ctx, updateMap, pl.Email)
	if uErr != nil {
		log.Println("Error updating user account:", err)
		return "error", err
	}
	return "Email Verification successful", nil
}

func (s *UserServiceImpl) Login(loginDto dto.LoginDTO) (dto.LoginResponseDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// get user by the email
	user, err := s.repo.GetUserByEmail(ctx, loginDto.Email)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return dto.LoginResponseDTO{}, errors.New("request timed out")
		}
		// TODO: handle case of user does not exists
		return dto.LoginResponseDTO{}, err
	}
	// compare password
	if !auth.ComparePassword(user.Password, loginDto.Password) {
		return dto.LoginResponseDTO{}, errors.New("invalid credentials")
	}

	secret := []byte(config.GetEnv("JWT_SECRET", "somesecret"))
	exp, expErr := strconv.Atoi(config.GetEnv("JWT_EXP", "3600"))
	if expErr != nil {
		log.Println("Error converting JWT_EXP to int: ", expErr)
		return dto.LoginResponseDTO{}, expErr
	}
	token, tokenErr := auth.CreateJWT(secret, exp, user)
	if tokenErr != nil {
		log.Println("Error creating JWT: ", tokenErr)
		return dto.LoginResponseDTO{}, errors.New("error creating a token")
	}
	if user.AccountDetails.AccountNo == "" {
		go s.generateVirtualAccount(user)
	}
	// return token
	return dto.LoginResponseDTO{
		AccessToken: token["token"],
		ExpiresIn:   token["expiresAt"],
		TokenType:   "jwt",
	}, nil

}

func (s *UserServiceImpl) GetUserDetails(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserServiceImpl) LoginSocial(pl dto.LoginSocialDTO) (dto.LoginResponseDTO, error) {
	email := pl.Email
	// get user by the email
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := s.repo.GetUserByEmail(ctx, email)
	// user does not exists create one
	if err != nil {
		// create user
		hashedPassword, hashErr := auth.HashPassword(pl.Email)
		if hashErr != nil {
			log.Println("Error hashing password: ", hashErr)
			return dto.LoginResponseDTO{}, hashErr
		}
		newUser := model.NewUser(pl.Name, pl.Email, hashedPassword)
		_, err := s.repo.CreateUser(ctx, newUser)
		if err != nil {
			log.Println("Error creating user: ", err)
			if errors.Is(err, context.DeadlineExceeded) {
				return dto.LoginResponseDTO{}, errors.New("request timed out")
			}
			return dto.LoginResponseDTO{}, err
		}
	}
	secret := []byte(config.GetEnv("JWT_SECRET", "somesecret"))
	exp, expErr := strconv.Atoi(config.GetEnv("JWT_EXP", "3600"))
	if expErr != nil {
		log.Println("Error converting JWT_EXP to int: ", expErr)
		return dto.LoginResponseDTO{}, expErr
	}
	token, tokenErr := auth.CreateJWT(secret, exp, user)
	if tokenErr != nil {
		log.Println("Error creating JWT: ", tokenErr)
		return dto.LoginResponseDTO{}, errors.New("error creating a token")
	}
	if user.AccountDetails.AccountNo == "" {
		go s.generateVirtualAccount(user)
	}
	return dto.LoginResponseDTO{
		AccessToken: token["token"],
		ExpiresIn:   token["expiresAt"],
		TokenType:   "jwt",
	}, nil

}

func (s *UserServiceImpl) PayazaWebhook(pl *dto.Transaction) (string, error) {
	fmt.Println("Incoming Data!")
	// get the user by the email
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// get the transaction reference and query transaction status
	data := ""
	tranRef := pl.TransactionReference
	url := fmt.Sprintf("https://api.payaza.africa/live/merchant-collection/transfer_notification_controller/transaction-query?transaction_reference=%s", tranRef)
	req, err := http.NewRequest("GET", url, strings.NewReader(data))
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Payaza %s", config.GetEnv("PAYAZA_API_KEY", "somkey")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TenantID", "test")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()
	res, rErr := ioutil.ReadAll(resp.Body)
	if rErr != nil {
		log.Println("Error reading response body:", rErr)
		return "", rErr
	}
	responseBody := new(dto.TransactionResponse)
	err = json.Unmarshal(res, responseBody)
	fmt.Println(responseBody)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return "", err
	}
	// check the status of the transaction
	if responseBody.Data.TransactionStatus == "Completed" && responseBody.Success {
		// get user by the virtual account number
		virtualAccountNumber := pl.VirtualAccountNumber
		user, err := s.repo.GetUserByVirtualAccount(ctx, virtualAccountNumber)
		if err != nil {
			log.Println("Error getting user by virtual account:", err)
			return "", err
		}

		newWalletHistory := model.NewWalletHistory(user.Email, float64(responseBody.Data.AmountReceived), user.WalletBalance, user.WalletBalance+float64(responseBody.Data.AmountReceived), responseBody.Data.TransactionReference)

		// credit the user account
		creditErr := s.repo.CreditUser(ctx, float64(responseBody.Data.AmountReceived), user.Email, newWalletHistory)
		if creditErr != nil {
			log.Println("Error crediting user account:", creditErr)
			return "", creditErr
		}
		// send a mail to the user
		// TODO: change email template for this
		message := "Your Wallet has been credited with " + strconv.FormatFloat(float64(responseBody.Data.AmountReceived), 'f', 2, 64)
		go utils.SendMail(user.Email, "Fund Successful", user.Name, message, "fund_success_template")
		return "success", nil
	}
	return "failed", errors.New("transaction failed")
}

func (s *UserServiceImpl) DisburseFunds(msg model.MessageData, inviteCode string) (string, error) {
	// get the user by the email
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	receiverUser, err := s.repo.GetUserByEmail(ctx, msg.ReceiverEmail)
	if err != nil {
		log.Println("Error getting user by email:", err)
		return "error", err
	}
	senderUser, err := s.repo.GetUserByEmail(ctx, msg.SenderEmail)
	if err != nil {
		log.Println("Error getting user by email:", err)
		return "error", err
	}
	tranRef := utils.GenerateReferenceCode()
	newWalletHistoryReceiver := model.NewWalletHistory(receiverUser.Email, float64(msg.Amount), receiverUser.WalletBalance, receiverUser.WalletBalance+float64(msg.Amount), tranRef)
	// debit the sender account
	newWalletHistorySender := model.NewWalletHistory(senderUser.Email, float64(msg.Amount), senderUser.WalletBalance, senderUser.WalletBalance-float64(msg.Amount), tranRef)
	// check if the sender has enough funds
	if senderUser.WalletBalance < float64(msg.Amount) {
		return "error", errors.New("insufficient funds")
	}
	debitErr := s.repo.DebitUser(ctx, float64(msg.Amount), senderUser.Email, newWalletHistorySender)
	if debitErr != nil {
		log.Println("Error debiting user account:", debitErr)
		return "error", debitErr
	}
	// credit the user account
	creditErr := s.repo.CreditUser(ctx, float64(msg.Amount), receiverUser.Email, newWalletHistoryReceiver)
	if creditErr != nil {
		log.Println("Error crediting user account:", creditErr)
		newFunds := model.NewFundsTracking(receiverUser.Email, receiverUser.Name, senderUser.Email, senderUser.Name, float64(msg.Amount), "failed", tranRef, inviteCode)
		_, err = s.repo.CreateNewFundsTracking(ctx, newFunds)
		if err != nil {
			log.Println("Error creating funds tracking record:", err)
		}
		// refund the sender
		newWalletHistorySender = model.NewWalletHistory(senderUser.Email, float64(msg.Amount), senderUser.WalletBalance, senderUser.WalletBalance+float64(msg.Amount), tranRef)
		_ = s.repo.CreditUser(ctx, float64(msg.Amount), senderUser.Email, newWalletHistorySender)
		return "error", creditErr
	}
	// create new successful transaction
	newFunds := model.NewFundsTracking(receiverUser.Email, receiverUser.Name, senderUser.Email, senderUser.Name, float64(msg.Amount), "success", tranRef, inviteCode)
	_, err = s.repo.CreateNewFundsTracking(ctx, newFunds)
	if err != nil {
		log.Println("Error creating funds tracking record:", err)
	}
	return "success", nil

}

func (s *UserServiceImpl) PayazaTestFundAccount(pl *dto.TestFundDTO) (*dto.TestFundAccountResponse, error) {
	// fund payload
	payload := map[string]interface{}{
		"service_type": "Account",
		"service_payload": map[string]interface{}{
			"request_application":              "Payaza",
			"application_module":               "USER_MODULE",
			"application_version":              "1.0.0",
			"request_class":                    "MerchantFundTestVirtualAccount",
			"account_number":                   pl.AccountNumber,
			"initiation_transaction_reference": utils.GenerateReferenceCode(),
			"transaction_amount":               pl.Amount,
			"currency":                         "NGN",
			"source_account_number":            "4859693408",
			"source_account_name":              "John Doe",
			"source_bank_name":                 "Eastern Bank",
		},
	}
	// marshal the payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshalling data: ", err)
		return nil, err
	}
	// send the request
	url := "https://router.prod.payaza.africa/api/request/secure/payloadhandler"
	req, rErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if rErr != nil {
		log.Println("Error occurred setting request: ", rErr)
		return nil, rErr
	}
	// set the request headers
	req.Header.Set("Authorization", fmt.Sprintf("Payaza %s", config.GetEnv("PAYAZA_API_KEY", "somkey")))
	req.Header.Set("X-TenantID", "test")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	payazaResp := new(dto.FundAccountPayazaResponse)
	err = json.Unmarshal(body, payazaResp)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return nil, err
	}
	return &dto.TestFundAccountResponse{
		Success: true,
		Message: payazaResp.ResponseMessage,
	}, nil

}
