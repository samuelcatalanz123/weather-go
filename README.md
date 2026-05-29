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
