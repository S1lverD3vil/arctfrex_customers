package report

import (
	"arctfrex-customers/internal/common"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type ReportApiClient interface {
	GetAccountsManifest() (*AccountGetManifestResponse, error)
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
	req, err := http.NewRequest(common.HTTP_METHOD_GET, "https://meta-integrator-arctfrex.ngrok.app/api/clients/AccountGet", nil)
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
