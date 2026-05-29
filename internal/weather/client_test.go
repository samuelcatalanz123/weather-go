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
