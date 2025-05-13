package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/model"
)

type MarketApiclient interface {
	// GetLatestMarketPrice() ([]ForexPriceResponse, error)
	GetLatestMarketPrice() (model.ArcMetaIntegratorPrice, error)
	GetLiveMarketUpdates() (*model.LiveMarketUpdatesResponse, error)
}

type marketApiclient struct{}

func NewMarketApiClient() MarketApiclient {
	return &marketApiclient{}
}

// func (s *marketApiclient) GetLatestMarketPrice() ([]ForexPriceResponse, error) {
func (s *marketApiclient) GetLatestMarketPrice() (model.ArcMetaIntegratorPrice, error) {
	// var forexPriceData []ForexPriceResponse
	var forexPriceData model.ArcMetaIntegratorPrice
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// Create a new request
	req, err := http.NewRequest(common.HTTP_METHOD_GET, os.Getenv(common.ARC_META_INTEGRATOR_BASEURL)+"/prices/get?symbol=XAUUSD.pkn,XAGUSD.pkn,USDJPY.pkn,JPK.pk,HKK.pk,GBPUSD.pkn,EURUSD.pkn,EURCHF.pkn,EURCAD.pkn,EURAUD.pkn,DJ.pk,CLSK.pkn,CHFJPY.pkn,AUDUSD.pkn,AUDNZD.pkn,AUDJPY.pkn", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, os.Getenv(common.ARC_META_INTEGRATOR_BASEURL)+"/prices?symbol=AUDJPY.fl,CHFJPY.fl,EURCAD.fl,EURGBP.fl,EURUSD.fl", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "https://enabled-simply-moth.ngrok-free.app/api/prices?symbol=AUDJPY.fl,CHFJPY.fl,EURCAD.fl,EURGBP.fl,EURUSD.fl", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "https://enabled-simply-moth.ngrok-free.app/api/prices?symbol=AUDJPY.fl,CHFJPY.fl,EURCAD.fl,EURGBP.fl,EURUSD.fl", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "https://enabled-simply-moth.ngrok-free.app/api/prices?symbol=AUDCAD,AUDCHF,AUDNZD,AUDUSD,EURUSD,GBPUSD,USDCAD,USDCHF,USDJPY", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "http://localhost:3000/arctfrex/latest/price", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "https://real-time-stock-finance-quote.p.rapidapi.com/forex/latest/AUDUSD,USDKRW,USDJPY,JPYKRW", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return model.ArcMetaIntegratorPrice{}, err
	}

	// // Add headers to the request
	// req.Header.Add("x-rapidapi-host", "real-time-stock-finance-quote.p.rapidapi.com")
	// req.Header.Add("x-rapidapi-key", "6deada9422msh4ea787410bc53bep1f5fabjsn9e80aeb0234f")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return model.ArcMetaIntegratorPrice{}, err
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return model.ArcMetaIntegratorPrice{}, err
	}

	log.Println("Response status:", resp.Status)

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(body, &forexPriceData)
	log.Println("Response:", forexPriceData)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)

		return model.ArcMetaIntegratorPrice{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// if len(forexPriceData) < 1 {
	// if forexPriceData == nil {
	// 	return nil, errors.New("not found data forex")
	// }

	return forexPriceData, nil
}

func (s *marketApiclient) GetLiveMarketUpdates() (*model.LiveMarketUpdatesResponse, error) {
	liveMarketUpdatesResponse := model.LiveMarketUpdatesResponse{
		Quotes: []model.LiveMarketUpdatesQuote{},
	}
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	// Create a new request
	req, err := http.NewRequest(common.HTTP_METHOD_GET, "http://localhost:3000/arctfrex/live/market/price", nil)
	// req, err := http.NewRequest(common.HTTP_METHOD_GET, "https://marketdata.tradermade.com/api/v1/live", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Add query parameters
	query := url.Values{}
	// query.Add("api_key", "D_8qKh8z8wd7a-LFIetc") // Add your API key
	// query.Add("currency", "AUDUSD,USDKRW,USDJPY,JPYKRW")
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

	//log.Println("Response status:", resp.Status)

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(body, &liveMarketUpdatesResponse)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)

		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &liveMarketUpdatesResponse, nil
}
