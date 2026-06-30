package temperature

// Temperature representa a temperatura nas três escalas, já no formato do contrato.
type Temperature struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

// FromCelsius converte uma temperatura em Celsius para as três escalas.
//
//	Fahrenheit = C * 1.8 + 32
//	Kelvin     = C + 273
func FromCelsius(celsius float64) Temperature {
	return Temperature{
		Celsius:    celsius,
		Fahrenheit: celsius*1.8 + 32,
		Kelvin:     celsius + 273,
	}
}
