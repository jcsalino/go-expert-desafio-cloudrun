package temperature

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromCelsius(t *testing.T) {
	tests := []struct {
		name    string
		celsius float64
		want    Temperature
	}{
		{"exemplo do contrato", 28.5, Temperature{Celsius: 28.5, Fahrenheit: 83.3, Kelvin: 301.5}},
		{"zero", 0, Temperature{Celsius: 0, Fahrenheit: 32, Kelvin: 273}},
		{"negativo", -10, Temperature{Celsius: -10, Fahrenheit: 14, Kelvin: 263}},
		{"cem", 100, Temperature{Celsius: 100, Fahrenheit: 212, Kelvin: 373}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromCelsius(tt.celsius)
			assert.InDelta(t, tt.want.Celsius, got.Celsius, 0.001)
			assert.InDelta(t, tt.want.Fahrenheit, got.Fahrenheit, 0.001)
			assert.InDelta(t, tt.want.Kelvin, got.Kelvin, 0.001)
		})
	}
}
