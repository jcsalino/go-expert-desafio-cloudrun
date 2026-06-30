package cep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

var (
	// ErrInvalidZipcode indica CEP fora do formato de 8 dígitos.
	ErrInvalidZipcode = errors.New("invalid zipcode")
	// ErrNotFound indica CEP válido no formato, mas inexistente.
	ErrNotFound = errors.New("can not find zipcode")
)

var cepRegex = regexp.MustCompile(`^[0-9]{8}$`)

// IsValid valida que o CEP tem exatamente 8 dígitos numéricos.
func IsValid(cep string) bool {
	return cepRegex.MatchString(cep)
}

// Finder abstrai a busca de cidade por CEP (permite mock nos testes).
type Finder interface {
	Find(ctx context.Context, cep string) (city string, err error)
}

// viaCEPResponse mapeia a resposta da ViaCEP.
//
// Obs.: para CEP inexistente a ViaCEP responde 200 com {"erro":"true"} (string)
// e sem "localidade". Em vez de tentar decodificar "erro" (cujo tipo varia entre
// versões da API), detectamos "não encontrado" pela ausência de localidade.
type viaCEPResponse struct {
	Localidade string `json:"localidade"`
}

// ViaCEP implementa Finder consultando https://viacep.com.br.
type ViaCEP struct {
	BaseURL string
	Client  *http.Client
}

func NewViaCEP() *ViaCEP {
	return &ViaCEP{
		BaseURL: "https://viacep.com.br/ws",
		Client:  http.DefaultClient,
	}
}

func (v *ViaCEP) Find(ctx context.Context, cep string) (string, error) {
	if !IsValid(cep) {
		return "", ErrInvalidZipcode
	}
	url := fmt.Sprintf("%s/%s/json/", v.BaseURL, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := v.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ErrNotFound
	}

	var data viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if data.Localidade == "" {
		return "", ErrNotFound
	}
	return data.Localidade, nil
}
