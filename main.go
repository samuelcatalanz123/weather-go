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
