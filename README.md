# Desafio Clima por CEP — Cloud Run

Desafio do curso **Go Expert**. Sistema em Go que recebe um CEP, descobre a cidade (ViaCEP) e retorna a temperatura atual em Celsius, Fahrenheit e Kelvin (WeatherAPI). Publicado no Google Cloud Run.

## URL (Cloud Run)

**https://desafio-cloudrun-329065667660.us-central1.run.app**

Exemplo de uso:

```sh
curl https://desafio-cloudrun-329065667660.us-central1.run.app/01001000
```

## Contrato

`GET /{cep}` — CEP com 8 dígitos.

| Situação           | Status | Corpo                       |
| ------------------ | ------ | --------------------------- |
| Sucesso            | 200    | `{"temp_C","temp_F","temp_K"}` |
| CEP inválido       | 422    | `invalid zipcode`           |
| CEP não encontrado | 404    | `can not find zipcode`      |

Sucesso (200):

```json
{ "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5 }
```

## Requisitos

- Go 1.25+ (ou Docker)
- Uma API key da [WeatherAPI](https://www.weatherapi.com/) (free tier)

## Rodar os testes

```sh
go test ./...
```

## Rodar localmente via Docker

```sh
docker build -t desafio-cloudrun .
docker run --rm -p 8080:8080 -e WEATHER_API_KEY=sua_chave desafio-cloudrun
```

Testar:

```sh
curl http://localhost:8080/01001000
```

## Deploy no Cloud Run

```sh
gcloud run deploy desafio-cloudrun \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=sua_chave
```
