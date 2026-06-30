package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"desafio-cloudrun/internal/cep"
	"desafio-cloudrun/internal/temperature"
	"desafio-cloudrun/internal/weather"
)

// WeatherHandler orquestra: valida o CEP, busca a cidade e a temperatura,
// e devolve as conversões no formato do contrato.
type WeatherHandler struct {
	CEPFinder     cep.Finder
	WeatherFinder weather.Finder
}

func NewWeatherHandler(cepFinder cep.Finder, weatherFinder weather.Finder) *WeatherHandler {
	return &WeatherHandler{
		CEPFinder:     cepFinder,
		WeatherFinder: weatherFinder,
	}
}

func (h *WeatherHandler) Handle(w http.ResponseWriter, r *http.Request) {
	zipcode := r.PathValue("cep")

	if !cep.IsValid(zipcode) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := h.CEPFinder.Find(r.Context(), zipcode)
	if err != nil {
		if errors.Is(err, cep.ErrInvalidZipcode) {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, cep.ErrNotFound) {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	celsius, err := h.WeatherFinder.Temperature(r.Context(), city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(temperature.FromCelsius(celsius))
}
