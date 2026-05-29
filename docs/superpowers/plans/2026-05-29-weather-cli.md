# App del clima (CLI en Go) — Plan de implementación

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Un CLI en Go (`clima <ciudad>`) que muestra el clima actual usando Open-Meteo (sin API key), encadenando geocoding → forecast.

**Architecture:** Un paquete `internal/weather` con un `Client` (URLs configurables para poder testear con httptest), tipos de dominio, y un traductor de códigos WMO a texto. `main.go` lee el argumento, encadena las dos llamadas e imprime el resultado.

**Tech Stack:** Go (net/http, encoding/json, httptest), API Open-Meteo.

---

## Estructura de archivos

```
weather-go/
  go.mod
  .gitignore
  main.go
  internal/weather/
    types.go
    conditions.go
    conditions_test.go
    client.go
    client_test.go
  README.md
```

Proyecto NUEVO en `/Users/mqr93ea/Repos/weather-go`. Comandos desde esa carpeta.

---

### Task 1: Inicializar el proyecto

**Files:**
- Create: `go.mod` (vía comando), `.gitignore`

- [ ] **Step 1: Inicializar el módulo Go y git**

Run desde `/Users/mqr93ea/Repos/weather-go`:
```bash
go mod init weather-go
git init
```
Expected: se crea `go.mod` con `module weather-go` y un repo git.

- [ ] **Step 2: Crear `.gitignore`**

```
# Binarios
/clima
weather-go
*.exe
```

- [ ] **Step 3: Commit**

```bash
git add -A && git -c user.name="Samuel" -c user.email="samuelcatalanz123@gmail.com" commit -m "chore: inicializar módulo weather-go"
```

---

### Task 2: Tipos y traductor de condiciones (TDD)

**Files:**
- Create: `internal/weather/types.go`, `internal/weather/conditions.go`
- Test: `internal/weather/conditions_test.go`

- [ ] **Step 1: Crear `internal/weather/types.go`**

```go
// Package weather fetches current weather from the Open-Meteo API.
package weather

// Place is a geocoded location.
type Place struct {
	Name    string
	Country string
	Lat     float64
	Lon     float64
}

// Weather is the current weather plus today's high/low for a place.
type Weather struct {
	TempC       float64
	MaxC        float64
	MinC        float64
	WindKmh     float64
	Code        int
	Description string
}
```

- [ ] **Step 2: Escribir el test que falla — `internal/weather/conditions_test.go`**

```go
package weather

import "testing"

func TestDescribeKnown(t *testing.T) {
	if got := describe(0); got != "Despejado" {
		t.Fatalf("describe(0)=%q, esperaba Despejado", got)
	}
	if got := describe(61); got != "Lluvia ligera" {
		t.Fatalf("describe(61)=%q, esperaba Lluvia ligera", got)
	}
}

func TestDescribeUnknown(t *testing.T) {
	if got := describe(1234); got != "Desconocido" {
		t.Fatalf("describe(1234)=%q, esperaba Desconocido", got)
	}
}
```

- [ ] **Step 3: Ejecutar y ver que falla**

Run: `go test ./internal/weather/`
Expected: FALLA (no existe `describe`).

- [ ] **Step 4: Crear `internal/weather/conditions.go`**

```go
package weather

// descriptions maps WMO weather codes to short Spanish text.
var descriptions = map[int]string{
	0:  "Despejado",
	1:  "Mayormente despejado",
	2:  "Parcialmente nublado",
	3:  "Nublado",
	45: "Niebla",
	48: "Niebla con escarcha",
	51: "Llovizna ligera",
	53: "Llovizna moderada",
	55: "Llovizna densa",
	61: "Lluvia ligera",
	63: "Lluvia moderada",
	65: "Lluvia fuerte",
	71: "Nieve ligera",
	73: "Nieve moderada",
	75: "Nieve fuerte",
	80: "Chubascos ligeros",
	81: "Chubascos moderados",
	82: "Chubascos violentos",
	95: "Tormenta",
	96: "Tormenta con granizo ligero",
	99: "Tormenta con granizo fuerte",
}

// describe returns a Spanish description for a WMO weather code, or
// "Desconocido" when the code is not recognized.
func describe(code int) string {
	if d, ok := descriptions[code]; ok {
		return d
	}
	return "Desconocido"
}
```

- [ ] **Step 5: Ejecutar y ver que pasa**

Run: `go test ./internal/weather/`
Expected: PASA (2 tests).

- [ ] **Step 6: Commit**

```bash
git add -A && git -c user.name="Samuel" -c user.email="samuelcatalanz123@gmail.com" commit -m "feat: tipos de dominio y traductor de condiciones (TDD)"
```

---

### Task 3: Cliente HTTP de Open-Meteo (TDD con httptest)

**Files:**
- Create: `internal/weather/client.go`
- Test: `internal/weather/client_test.go`

- [ ] **Step 1: Escribir el test que falla — `internal/weather/client_test.go`**

```go
package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(srv *httptest.Server) *Client {
	return &Client{geoBaseURL: srv.URL, forecastBaseURL: srv.URL, http: srv.Client()}
}

func TestGeocodeFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":[{"name":"Ciudad de Guatemala","country":"Guatemala","latitude":14.64,"longitude":-90.51}]}`))
	}))
	defer srv.Close()

	place, err := newTestClient(srv).Geocode(context.Background(), "Guatemala")
	if err != nil {
		t.Fatal(err)
	}
	if place.Name != "Ciudad de Guatemala" || place.Country != "Guatemala" {
		t.Fatalf("place inesperado: %+v", place)
	}
	if place.Lat != 14.64 || place.Lon != -90.51 {
		t.Fatalf("coords inesperadas: %+v", place)
	}
}

func TestGeocodeNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":[]}`))
	}))
	defer srv.Close()

	if _, err := newTestClient(srv).Geocode(context.Background(), "Xyz"); err == nil {
		t.Fatal("esperaba error de ciudad no encontrada")
	}
}

func TestForecastParses(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"current":{"temperature_2m":24.0,"weather_code":2,"wind_speed_10m":12.0},"daily":{"temperature_2m_max":[27.0],"temperature_2m_min":[18.0]}}`))
	}))
	defer srv.Close()

	wth, err := newTestClient(srv).Forecast(context.Background(), 14.64, -90.51)
	if err != nil {
		t.Fatal(err)
	}
	if wth.TempC != 24.0 || wth.MaxC != 27.0 || wth.MinC != 18.0 || wth.WindKmh != 12.0 {
		t.Fatalf("weather inesperado: %+v", wth)
	}
	if wth.Code != 2 || wth.Description != "Parcialmente nublado" {
		t.Fatalf("descripción inesperada: %+v", wth)
	}
}
```

- [ ] **Step 2: Ejecutar y ver que falla**

Run: `go test ./internal/weather/`
Expected: FALLA (no existe `Client`).

- [ ] **Step 3: Crear `internal/weather/client.go`**

```go
package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client talks to the Open-Meteo geocoding and forecast APIs.
type Client struct {
	geoBaseURL      string
	forecastBaseURL string
	http            *http.Client
}

// NewClient returns a Client pointing at the real Open-Meteo endpoints.
func NewClient() *Client {
	return &Client{
		geoBaseURL:      "https://geocoding-api.open-meteo.com/v1/search",
		forecastBaseURL: "https://api.open-meteo.com/v1/forecast",
		http:            &http.Client{Timeout: 10 * time.Second},
	}
}

// Geocode resolves a city name to a Place, returning an error when the city
// is not found.
func (c *Client) Geocode(ctx context.Context, city string) (Place, error) {
	q := url.Values{}
	q.Set("name", city)
	q.Set("count", "1")
	q.Set("language", "es")

	var body struct {
		Results []struct {
			Name      string  `json:"name"`
			Country   string  `json:"country"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}
	if err := c.getJSON(ctx, c.geoBaseURL+"?"+q.Encode(), &body); err != nil {
		return Place{}, err
	}
	if len(body.Results) == 0 {
		return Place{}, fmt.Errorf("no se encontró la ciudad %q", city)
	}
	r := body.Results[0]
	return Place{Name: r.Name, Country: r.Country, Lat: r.Latitude, Lon: r.Longitude}, nil
}

// Forecast returns the current weather and today's high/low at the given
// coordinates.
func (c *Client) Forecast(ctx context.Context, lat, lon float64) (Weather, error) {
	q := url.Values{}
	q.Set("latitude", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Set("longitude", strconv.FormatFloat(lon, 'f', -1, 64))
	q.Set("current", "temperature_2m,weather_code,wind_speed_10m")
	q.Set("daily", "temperature_2m_max,temperature_2m_min")
	q.Set("timezone", "auto")

	var body struct {
		Current struct {
			Temperature float64 `json:"temperature_2m"`
			WeatherCode int     `json:"weather_code"`
			WindSpeed   float64 `json:"wind_speed_10m"`
		} `json:"current"`
		Daily struct {
			TempMax []float64 `json:"temperature_2m_max"`
			TempMin []float64 `json:"temperature_2m_min"`
		} `json:"daily"`
	}
	if err := c.getJSON(ctx, c.forecastBaseURL+"?"+q.Encode(), &body); err != nil {
		return Weather{}, err
	}

	w := Weather{
		TempC:       body.Current.Temperature,
		WindKmh:     body.Current.WindSpeed,
		Code:        body.Current.WeatherCode,
		Description: describe(body.Current.WeatherCode),
	}
	if len(body.Daily.TempMax) > 0 {
		w.MaxC = body.Daily.TempMax[0]
	}
	if len(body.Daily.TempMin) > 0 {
		w.MinC = body.Daily.TempMin[0]
	}
	return w, nil
}

// getJSON performs a GET request and decodes a successful JSON response into dst.
func (c *Client) getJSON(ctx context.Context, endpoint string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("error de red: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("respuesta inesperada del servicio (código %d)", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		return fmt.Errorf("respuesta no válida: %w", err)
	}
	return nil
}
```

- [ ] **Step 4: Ejecutar y ver que pasa**

Run: `go test ./internal/weather/`
Expected: PASA (5 tests en total).

- [ ] **Step 5: Commit**

```bash
git add -A && git -c user.name="Samuel" -c user.email="samuelcatalanz123@gmail.com" commit -m "feat: cliente Open-Meteo (geocode + forecast) con pruebas httptest"
```

---

### Task 4: Entrypoint `main.go`

**Files:**
- Create: `main.go`

- [ ] **Step 1: Crear `main.go`**

```go
// Command clima prints the current weather for a city using the Open-Meteo API.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"weather-go/internal/weather"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: clima <ciudad>")
	}
	city := strings.Join(args, " ")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := weather.NewClient()

	place, err := client.Geocode(ctx, city)
	if err != nil {
		return err
	}
	w, err := client.Forecast(ctx, place.Lat, place.Lon)
	if err != nil {
		return err
	}

	fmt.Printf("🌤️  %s, %s\n", place.Name, place.Country)
	fmt.Printf("    Temperatura: %.0f°C  (máx %.0f° / mín %.0f°)\n", w.TempC, w.MaxC, w.MinC)
	fmt.Printf("    Estado: %s\n", w.Description)
	fmt.Printf("    Viento: %.0f km/h\n", w.WindKmh)
	return nil
}
```

- [ ] **Step 2: Verificar compilación y comportamiento sin argumento**

Run desde la raíz del proyecto:
```bash
go build ./... && echo "build OK"
go run . ; echo "código de salida: $?"
```
Expected: `build OK`; al correr sin argumento imprime `Error: uso: clima <ciudad>` y código de salida `1`.

- [ ] **Step 3: Commit**

```bash
git add -A && git -c user.name="Samuel" -c user.email="samuelcatalanz123@gmail.com" commit -m "feat: entrypoint CLI que encadena geocode y forecast"
```

---

### Task 5: README y verificación final

**Files:**
- Create: `README.md`

- [ ] **Step 1: Crear `README.md`**

```markdown
# clima — App del clima (CLI en Go)

> 🌐 Languages: **Español** · English below

Programa de línea de comandos que muestra el clima actual de cualquier ciudad,
usando la API gratuita [Open-Meteo](https://open-meteo.com) (sin API key).

## Uso

```bash
go run . Guatemala
```

Ejemplo de salida:

```
🌤️  Ciudad de Guatemala, Guatemala
    Temperatura: 24°C  (máx 27° / mín 18°)
    Estado: Parcialmente nublado
    Viento: 12 km/h
```

La ciudad puede tener varias palabras: `go run . San Salvador`.

## Cómo funciona

1. **Geocodificación:** convierte el nombre de la ciudad en coordenadas
   (Open-Meteo Geocoding API).
2. **Pronóstico:** con esas coordenadas, obtiene el clima actual y la máxima/
   mínima del día (Open-Meteo Forecast API).

## Stack

- **Go** (net/http, encoding/json) — sin dependencias externas.
- **Open-Meteo** como fuente de datos (gratis, sin registro).
- **Pruebas** con `net/http/httptest` (no requieren internet).

## Pruebas

```bash
go test ./...
```

## Compilar un binario

```bash
go build -o clima .
./clima Guatemala
```

---

## English

A command-line app that shows the current weather for any city using the free
[Open-Meteo](https://open-meteo.com) API (no API key). It geocodes the city
name to coordinates, then fetches the current weather and today's high/low.
Run with `go run . <city>`. Built in Go with no external dependencies; tested
with `net/http/httptest`.
```

- [ ] **Step 2: Verificación final**

Run desde la raíz del proyecto:
```bash
go build ./... && go vet ./... && go test ./... && echo "TODO OK"
```
Expected: `TODO OK` (5 pruebas pasan; build y vet limpios).

- [ ] **Step 3: Prueba real (requiere internet)**

```bash
go run . Guatemala
```
Expected: imprime el clima actual de la Ciudad de Guatemala. (Si no hay
internet, dará un error de red claro — eso también es correcto.)

- [ ] **Step 4: Commit final**

```bash
git add -A && git -c user.name="Samuel" -c user.email="samuelcatalanz123@gmail.com" commit -m "docs: README del proyecto clima"
```

---

## Notas de verificación (self-review del plan)

- **Cobertura del spec:** módulo + git (Task 1), tipos + condiciones (Task 2),
  cliente geocode/forecast con httptest (Task 3), entrypoint que encadena las
  dos APIs + errores (Task 4), README + verificación + prueba real (Task 5).
- **Sin placeholders:** todo el código está completo.
- **Consistencia de tipos/firmas:** `Place{Name,Country,Lat,Lon}`,
  `Weather{TempC,MaxC,MinC,WindKmh,Code,Description}`, `NewClient()`,
  `Client.Geocode(ctx,string) (Place,error)`, `Client.Forecast(ctx,float64,
  float64) (Weather,error)`, `describe(int) string`. Los tests usan los campos
  internos `geoBaseURL/forecastBaseURL/http` del `Client` (test interno, mismo
  paquete). Mismos nombres en todas las tareas.
- **Atribución:** commits a nombre de Samuel + (se añadirá co-autoría de Claude
  al finalizar, igual que el otro proyecto, si el usuario lo desea).
```
