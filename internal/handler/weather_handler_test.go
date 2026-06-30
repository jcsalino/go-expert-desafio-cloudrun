package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"desafio-cloudrun/internal/cep"
	"desafio-cloudrun/internal/temperature"

	"github.com/stretchr/testify/assert"
)

// mocks
type mockCEPFinder struct {
	city string
	err  error
}

func (m mockCEPFinder) Find(ctx context.Context, c string) (string, error) {
	return m.city, m.err
}

type mockWeatherFinder struct {
	celsius float64
	err     error
}

func (m mockWeatherFinder) Temperature(ctx context.Context, city string) (float64, error) {
	return m.celsius, m.err
}

func doRequest(h *WeatherHandler, zipcode string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{cep}", h.Handle)
	req := httptest.NewRequest(http.MethodGet, "/"+zipcode, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestHandle_Success(t *testing.T) {
	h := NewWeatherHandler(
		mockCEPFinder{city: "São Paulo"},
		mockWeatherFinder{celsius: 28.5},
	)
	rec := doRequest(h, "01001000")

	assert.Equal(t, http.StatusOK, rec.Code)

	var body temperature.Temperature
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.InDelta(t, 28.5, body.Celsius, 0.001)
	assert.InDelta(t, 83.3, body.Fahrenheit, 0.001)
	assert.InDelta(t, 301.5, body.Kelvin, 0.001)
}

func TestHandle_InvalidZipcode(t *testing.T) {
	h := NewWeatherHandler(mockCEPFinder{}, mockWeatherFinder{})
	rec := doRequest(h, "123")

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid zipcode")
}

func TestHandle_NotFound(t *testing.T) {
	h := NewWeatherHandler(
		mockCEPFinder{err: cep.ErrNotFound},
		mockWeatherFinder{},
	)
	rec := doRequest(h, "99999999")

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "can not find zipcode")
}
