package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeatherAPI_Temperature_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/current.json", r.URL.Path)
		assert.Equal(t, "test-key", r.URL.Query().Get("key"))
		assert.Equal(t, "São Paulo", r.URL.Query().Get("q"))
		w.Write([]byte(`{"location":{"name":"Sao Paulo"},"current":{"temp_c":28.5}}`))
	}))
	defer server.Close()

	api := &WeatherAPI{APIKey: "test-key", BaseURL: server.URL, Client: server.Client()}
	celsius, err := api.Temperature(context.Background(), "São Paulo")
	assert.NoError(t, err)
	assert.InDelta(t, 28.5, celsius, 0.001)
}

func TestWeatherAPI_Temperature_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	api := &WeatherAPI{APIKey: "test-key", BaseURL: server.URL, Client: server.Client()}
	_, err := api.Temperature(context.Background(), "Cidade Inexistente")
	assert.Error(t, err)
}
