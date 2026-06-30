package main

import (
	"log"
	"net/http"
	"os"

	"desafio-cloudrun/internal/cep"
	"desafio-cloudrun/internal/handler"
	"desafio-cloudrun/internal/weather"
)

func main() {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY não definida")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // padrão do Cloud Run
	}

	weatherHandler := handler.NewWeatherHandler(
		cep.NewViaCEP(),
		weather.NewWeatherAPI(apiKey),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{cep}", weatherHandler.Handle)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	// Raiz com instrução de uso (evita 404 ao abrir a URL base no navegador).
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"usage":"GET /{cep} — ex: /01001000"}` + "\n"))
	})

	log.Printf("servidor rodando na porta %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
