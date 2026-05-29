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
