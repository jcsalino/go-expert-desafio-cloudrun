package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Finder abstrai a consulta de temperatura por cidade (permite mock nos testes).
type Finder interface {
	Temperature(ctx context.Context, city string) (celsius float64, err error)
}

// weatherAPIResponse mapeia o campo usado da resposta da WeatherAPI (current.json).
type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// WeatherAPI implementa Finder usando https://www.weatherapi.com (endpoint current.json).
type WeatherAPI struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewWeatherAPI(apiKey string) *WeatherAPI {
	return &WeatherAPI{
		APIKey:  apiKey,
		BaseURL: "http://api.weatherapi.com/v1",
		Client:  http.DefaultClient,
	}
}

func (wa *WeatherAPI) Temperature(ctx context.Context, city string) (float64, error) {
	endpoint := fmt.Sprintf(
		"%s/current.json?key=%s&q=%s&aqi=no",
		wa.BaseURL, wa.APIKey, url.QueryEscape(city),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}
	resp, err := wa.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather api retornou status %d", resp.StatusCode)
	}

	var data weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Current.TempC, nil
}
