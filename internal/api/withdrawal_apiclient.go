package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/model"
)

type WithdrawalApiclient interface {
	TradeWithdrawal(tradeWithdrawal model.TradeWithdrawal) (model.TradeWithdrawal, error)
}

type withdrawalApiclient struct{}

func NewWithdrawalApiclient() WithdrawalApiclient {
	return &withdrawalApiclient{}
}
func (s *withdrawalApiclient) TradeWithdrawal(tradeWithdrawal model.TradeWithdrawal) (model.TradeWithdrawal, error) {
	var tradeWithdrawalData model.TradeWithdrawal

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(tradeWithdrawal)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return tradeWithdrawal, err
	}

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, os.Getenv(common.ARC_META_INTEGRATOR_BASEURL)+"/Trade/Add", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return tradeWithdrawal, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return tradeWithdrawal, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return tradeWithdrawal, err
	}

	// Log the response status
	log.Println("Response status:", resp.Status)

	// Parse the response JSON into the struct
	err = json.Unmarshal(body, &tradeWithdrawalData)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		// log.Fatalf("Error unmarshaling JSON: %v", err)
		return tradeWithdrawal, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Println("Response:", tradeWithdrawalData)

	return tradeWithdrawalData, nil
}
