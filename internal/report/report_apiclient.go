package report

import (
	"arctfrex-customers/internal/common"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ReportApiClient interface {
	GetAccountsManifest() (*AccountGetManifestResponse, error)
	GetGroupUserLogins(GroupUserLoginsRequest) (*GroupUserLoginsResponse, error)
}

type reportApiClient struct{}

func NewReportApiClient() ReportApiClient {
	return &reportApiClient{}
}

func (s *reportApiClient) GetAccountsManifest() (*AccountGetManifestResponse, error) {
	accountGetManifestResponse := AccountGetManifestResponse{}
	// Quotes: []LiveMarketUpdatesQuote{},
	//}
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// Create a new request
	req, err := http.NewRequest(common.HTTP_METHOD_GET, os.Getenv(common.ARC_META_INTEGRATOR_BASEURL)+"/clients/AccountGet", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "http://localhost:3000/arctfrex/api/clients/AccountGet", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "https://marketdata.tradermade.com/api/v1/live", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Add query parameters
	query := url.Values{}
	// query.Add("api_key", "D_8qKh8z8wd7a-LFIetc") // Add your API key
	// query.Add("currency", "AUDUSD,USDKRW,USDJPY,JPYKRW")
	query.Add("login", "2175778")
	req.URL.RawQuery = query.Encode()

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	log.Println("Response status:", resp.Status)

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(body, &accountGetManifestResponse)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)

		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &accountGetManifestResponse, nil
}

func (s *reportApiClient) GetGroupUserLogins(groupUserLoginsRequest GroupUserLoginsRequest) (*GroupUserLoginsResponse, error) {
	groupUserLoginsResponse := GroupUserLoginsResponse{}
	// Quotes: []LiveMarketUpdatesQuote{},
	//}
	// Convert the request body to JSON
	jsonBody, err := json.Marshal(groupUserLoginsRequest)
	if err != nil {
		fmt.Println("Error marshalling request body:", err)
		return &groupUserLoginsResponse, err
	}

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// Create a new request
	req, err := http.NewRequest(common.HTTP_METHOD_POST, os.Getenv(common.ARC_META_INTEGRATOR_BASEURL)+"/clients/Logins", bytes.NewBuffer(jsonBody))
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "http://localhost:3000/arctfrex/api/clients/AccountGet", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "https://marketdata.tradermade.com/api/v1/live", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// // Add query parameters
	// query := url.Values{}
	// // query.Add("api_key", "D_8qKh8z8wd7a-LFIetc") // Add your API key
	// // query.Add("currency", "AUDUSD,USDKRW,USDJPY,JPYKRW")
	// //query.Add("login", "2175778")
	// req.URL.RawQuery = query.Encode()

	// Add headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	log.Println("Response status:", resp.Status)

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(body, &groupUserLoginsResponse)
	if err != nil {
		log.Println("Error unmarshalling JSON groupUserLoginsResponse:", err)

		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &groupUserLoginsResponse, nil
}
