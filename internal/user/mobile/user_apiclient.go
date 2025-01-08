package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type UserApiclient interface {
	ClientAdd(clientAdd ClientAdd) (ClientAdd, error)
	DemoAccountTopUp(demoAccountTopUp DemoAccountTopUp) (DemoAccountTopUp, error)
}

type userApiclient struct{}

func NewUserApiclient() UserApiclient {
	return &userApiclient{}
}

func (s *userApiclient) ClientAdd(clientAdd ClientAdd) (ClientAdd, error) {
	var clientAddData ClientAdd

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(clientAdd)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return clientAdd, err
	}

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, "https://meta-integrator-arctfrex.ngrok.app/api/Clients/Add", bytes.NewBuffer(jsonBody))
	// req, err := http.NewRequest(http.MethodPost, "https://enabled-simply-moth.ngrok-free.app/api/Clients/Add", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return clientAdd, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return clientAdd, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return clientAdd, err
	}

	// Log the response status
	log.Println("Response status:", resp.Status)

	// Parse the response JSON into the struct
	err = json.Unmarshal(body, &clientAddData)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		// log.Fatalf("Error unmarshaling JSON: %v", err)
		return clientAdd, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Println("Response:", clientAddData)

	return clientAddData, nil
}

func (s *userApiclient) DemoAccountTopUp(demoAccountTopUp DemoAccountTopUp) (DemoAccountTopUp, error) {
	var demoAccountTopUpData DemoAccountTopUp

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(demoAccountTopUp)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return demoAccountTopUp, err
	}

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, "https://meta-integrator-arctfrex.ngrok.app/api/Trade/Add", bytes.NewBuffer(jsonBody))
	// req, err := http.NewRequest(http.MethodPost, "https://enabled-simply-moth.ngrok-free.app/api/Trade/Add", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return demoAccountTopUp, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return demoAccountTopUp, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return demoAccountTopUp, err
	}

	// Log the response status
	log.Println("Response status:", resp.Status)

	// Parse the response JSON into the struct
	err = json.Unmarshal(body, &demoAccountTopUpData)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		// log.Fatalf("Error unmarshaling JSON: %v", err)
		return demoAccountTopUp, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Println("Response:", demoAccountTopUpData)

	return demoAccountTopUpData, nil
}
