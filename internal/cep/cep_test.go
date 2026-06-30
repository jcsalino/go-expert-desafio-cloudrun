package cep

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	assert.True(t, IsValid("01001000"))
	assert.False(t, IsValid("1234567"))   // 7 dígitos
	assert.False(t, IsValid("123456789")) // 9 dígitos
	assert.False(t, IsValid("0100100a"))  // caractere inválido
	assert.False(t, IsValid(""))
}

func TestViaCEP_Find_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"cep":"01001-000","localidade":"São Paulo","uf":"SP"}`))
	}))
	defer server.Close()

	finder := &ViaCEP{BaseURL: server.URL, Client: server.Client()}
	city, err := finder.Find(context.Background(), "01001000")
	assert.NoError(t, err)
	assert.Equal(t, "São Paulo", city)
}

func TestViaCEP_Find_NotFound(t *testing.T) {
	// A ViaCEP responde 200 com "erro" no formato string para CEP inexistente.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"erro": "true"}`))
	}))
	defer server.Close()

	finder := &ViaCEP{BaseURL: server.URL, Client: server.Client()}
	_, err := finder.Find(context.Background(), "99999999")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestViaCEP_Find_InvalidZipcode(t *testing.T) {
	finder := NewViaCEP()
	_, err := finder.Find(context.Background(), "123")
	assert.ErrorIs(t, err, ErrInvalidZipcode)
}
