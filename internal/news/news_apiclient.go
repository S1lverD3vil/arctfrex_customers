package news

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

type NewsApiclient interface {
	GetLatestNews() (*NewsResponse, error)
	GetLatestNewsBulletin() (*NewsBulletinResponse, error)
}

type newsApiclient struct{}

func NewNewsApiClient() NewsApiclient {
	return &newsApiclient{}
}

func (na *newsApiclient) GetLatestNews() (*NewsResponse, error) {
	newsResponse := NewsResponse{
		Data: []NewsData{},
	}
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// Create a new request
	req, err := http.NewRequest(common.HTTP_METHOD_GET, "http://localhost:3000/arctfrex/latest/news", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "https://forexnewsapi.com/api/v1", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Add query parameters
	query := url.Values{}
	// query.Add("token", "cxx1uuydcvf84j7l5i9wbzmvic9r4pnscmqf0mvl") // Add your API key
	// query.Add("items", "3")                                        // Add your API key
	// query.Add("currencypair", "AUD-USD")
	//query.Add("currencypair", "XAU-USD,USD-JPY,AUD-USD")
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

	//log.Println("News Client Response status:", resp.Status)

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(body, &newsResponse)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)

		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	//log.Println("Response status:", &newsResponse)

	return &newsResponse, nil
}

func (na *newsApiclient) GetLatestNewsBulletin() (*NewsBulletinResponse, error) {
	newsBulletinResponse := NewsBulletinResponse{
		Results: []Article{},
	}
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// Create a new request
	req, err := http.NewRequest(common.HTTP_METHOD_GET, "http://localhost:3000/arctfrex/latest/bulletin", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "https://newsdata.io/api/1/news", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Add query parameters
	query := url.Values{}
	// query.Add("apikey", "pub_5597818e1829bb6f9cb8c6ed47ded39faa364") // Add your API key
	// query.Add("country", "id,us")                                    // Add your API key
	// query.Add("category", "business ")
	// //query.Add("currencypair", "XAU-USD,USD-JPY,AUD-USD")
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

	//log.Println("News Client Response status:", resp.Status)
	//log.Println("News Client Response body:", body)

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(body, &newsBulletinResponse)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)

		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	//log.Println("Response status:", &newsResponse)

	return &newsBulletinResponse, nil
}
