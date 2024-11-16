package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/LoginX/SprayDash/config"
	"github.com/LoginX/SprayDash/internal/model"
	"github.com/LoginX/SprayDash/internal/utils"
)

func GeneratePayazaVirtualAccount(user *model.User) (*utils.ResponseBody, error) {
	url := "https://router.prod.payaza.africa/api/request/secure/payloadhandler"
	userFirstName := strings.Fields(user.Name)[0]
	userLastName := ""
	if len(strings.Fields(user.Name)) > 1 {
		userLastName = strings.Fields(user.Name)[1]
	}
	payload := map[string]interface{}{
		"service_type": "Account",
		"service_payload": map[string]interface{}{
			"request_application":      "Payaza",
			"application_module":       "USER_MODULE",
			"application_version":      "1.0.0",
			"request_class":            "CreateReservedAccountForCustomers",
			"customer_first_name":      userFirstName,
			"customer_last_name":       userLastName,
			"customer_email":           user.Email,
			"customer_phone":           "09012345673",
			"virtual_account_provider": "Premiumtrust",
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Payaza %s", config.GetEnv("PAYAZA_API_KEY", "somkey")))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()
	print(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	respBody := []byte(string(body))
	var responseBody utils.ResponseBody
	jErr := json.Unmarshal(respBody, &responseBody)
	if jErr != nil {
		fmt.Println("Error unmarshalling JSON:", jErr)
		return nil, jErr
	}
	return &responseBody, nil
}
