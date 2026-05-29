# Diseño: App del clima en Go (CLI)

**Fecha:** 2026-05-29
**Estado:** Aprobado para escribir el plan de implementación
**Autor del proyecto:** Samuel (2º proyecto de portafolio)

## Objetivo

Un programa de línea de comandos (CLI) en Go que, dado el nombre de una ciudad,
muestra el clima actual consumiendo la API gratuita **Open-Meteo** (sin registro
ni API key). Demuestra consumo de APIs de terceros, parseo de JSON, encadenar
peticiones y pruebas con servidores HTTP simulados.

## Decisiones tomadas (brainstorming)

| Tema        | Decisión                                                   |
| ----------- | ---------------------------------------------------------- |
| Forma       | **CLI**: `clima <ciudad>` imprime el tiempo en la terminal |
| Datos       | **Open-Meteo** (gratis, sin registro ni llave)             |
| Flujo       | 2 llamadas: geocoding (ciudad→coords) → forecast (coords→clima) |
| Muestra     | Temp. actual, máx/mín del día, estado (texto ES), viento   |
| Repo        | Proyecto y repositorio **nuevos** (`weather-go`)           |

## Cómo funciona (encadenar dos APIs, ambas sin llave)

1. **Geocoding:** `GET https://geocoding-api.open-meteo.com/v1/search?name=<ciudad>&count=1&language=es`
   → devuelve `results[0]` con `latitude`, `longitude`, `name`, `country`.
2. **Forecast:** `GET https://api.open-meteo.com/v1/forecast?latitude=<lat>&longitude=<lon>&current=temperature_2m,weather_code,wind_speed_10m&daily=temperature_2m_max,temperature_2m_min&timezone=auto`
   → devuelve `current` (temperatura, código de clima, viento) y `daily`
   (máx/mín de hoy).

El `weather_code` (un número, estándar WMO) se traduce a texto en español
(p. ej. `0`→"Despejado", `2`→"Parcialmente nublado", `61`→"Lluvia ligera").

## Estructura del código

```
weather-go/
  go.mod
  main.go                      Punto de entrada: lee el argumento (ciudad),
                               llama a la capa weather, imprime el resultado.
  internal/weather/
    client.go                  Client con BaseURLs configurables; Geocode(ciudad)
                               (Place) y Forecast(lat,lon) (Weather).
    client_test.go             Pruebas con httptest (servidores HTTP simulados).
    conditions.go              describe(code int) string — código WMO → texto ES.
    conditions_test.go         Pruebas de la traducción.
    types.go                   Tipos: Place, Weather (y structs internos de parseo).
  README.md                    Descripción, uso, ejemplo, stack.
```

- **client.go:** un `Client` con campos `geoBaseURL`, `forecastBaseURL` y un
  `*http.Client`. `NewClient()` usa las URLs reales de Open-Meteo; en las
  pruebas se inyectan URLs de un `httptest.Server`. Métodos:
  - `Geocode(ctx, city string) (Place, error)` — error si no hay resultados.
  - `Forecast(ctx, lat, lon float64) (Weather, error)`.
- **types.go:** `Place{Name, Country string; Lat, Lon float64}` y
  `Weather{TempC, MaxC, MinC, WindKmh float64; Code int}`.
- **conditions.go:** `describe(code int) string` con un `map[int]string` y un
  texto por defecto ("Desconocido") para códigos no mapeados.
- **main.go:** valida que haya un argumento; encadena `Geocode` →`Forecast`;
  formatea e imprime; mapea errores a mensajes claros y código de salida ≠ 0.

## Datos que se muestran

```
🌤️  <Ciudad>, <País>
    Temperatura: <X>°C  (máx <Y>° / mín <Z>°)
    Estado: <texto en español>
    Viento: <W> km/h
```

## Manejo de errores (mensajes claros, salida ≠ 0)

- Sin argumento → `Uso: clima <ciudad>` y salida 1.
- Ciudad no encontrada (geocoding sin resultados) → `No se encontró la ciudad "<x>"`.
- Fallo de red / respuesta no-200 / JSON inválido → mensaje entendible.

## Pruebas

- **conditions_test.go:** `describe` devuelve el texto correcto para códigos
  conocidos y "Desconocido" para uno no mapeado.
- **client_test.go:** con `httptest.Server` devolviendo JSON de ejemplo,
  `Geocode` parsea bien un resultado y devuelve error con `results` vacío;
  `Forecast` parsea `current`/`daily` correctamente. (Sin internet real.)
- `go build ./...`, `go vet ./...` y `go test ./...` limpios.

## Fuera de alcance (YAGNI)

Pronóstico de varios días, salida en otros formatos (JSON/colores avanzados),
caché, interfaz gráfica o web, y selección entre varias ciudades homónimas
(usamos el primer resultado).

## Criterios de éxito

1. `go run . Guatemala` imprime el clima actual de la ciudad.
2. Una ciudad inexistente y la falta de argumento dan mensajes claros y salida ≠ 0.
3. Las pruebas pasan sin necesidad de internet (usando httptest).
4. El proyecto queda en su propio repositorio, listo para GitHub.
