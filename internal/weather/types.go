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
