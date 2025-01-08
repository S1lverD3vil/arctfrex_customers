package deposit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type DepositApiclient interface {
	TradeDeposit(tradeDeposit TradeDeposit) (TradeDeposit, error)
}

type depositApiclient struct{}

func NewDepositApiclient() DepositApiclient {
	return &depositApiclient{}
}
func (s *depositApiclient) TradeDeposit(tradeDeposit TradeDeposit) (TradeDeposit, error) {
	var tradeDepositData TradeDeposit

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(tradeDeposit)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return tradeDeposit, err
	}

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, "https://meta-integrator-arctfrex.ngrok.app/api/Trade/Add", bytes.NewBuffer(jsonBody))
	// req, err := http.NewRequest(http.MethodPost, "https://enabled-simply-moth.ngrok-free.app/api/Trade/Add", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return tradeDeposit, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return tradeDeposit, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return tradeDeposit, err
	}

	// Log the response status
	log.Println("Response status:", resp.Status)

	// Parse the response JSON into the struct
	err = json.Unmarshal(body, &tradeDepositData)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		// log.Fatalf("Error unmarshaling JSON: %v", err)
		return tradeDeposit, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Println("Response:", tradeDepositData)

	return tradeDepositData, nil
}
