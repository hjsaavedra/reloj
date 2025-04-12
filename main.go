package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type WeatherResponse struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
		WeatherCode int     `json:"weather_code"`
	} `json:"current"`
}

func getWeather() (*WeatherResponse, error) {
	// Coordenadas de Santo Domingo, Chile
	url := "https://api.open-meteo.com/v1/forecast?latitude=-33.63&longitude=-71.63&current=temperature_2m,weather_code&timezone=America%2FSantiago"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, err
	}

	return &weather, nil
}

func getWeatherIcon(code int) string {
	// Mapeo de códigos de clima a emojis
	weatherIcons := map[int]string{
		0:  "☀️", // Clear sky
		1:  "🌤️", // Mainly clear
		2:  "⛅",  // Partly cloudy
		3:  "☁️", // Overcast
		45: "🌫️", // Fog
		48: "🌫️", // Depositing rime fog
		51: "🌧️", // Light drizzle
		53: "🌧️", // Moderate drizzle
		55: "🌧️", // Dense drizzle
		61: "🌧️", // Slight rain
		63: "🌧️", // Moderate rain
		65: "🌧️", // Heavy rain
		71: "🌨️", // Slight snow
		73: "🌨️", // Moderate snow
		75: "🌨️", // Heavy snow
		77: "🌨️", // Snow grains
		80: "🌧️", // Slight rain showers
		81: "🌧️", // Moderate rain showers
		82: "🌧️", // Violent rain showers
		85: "🌨️", // Slight snow showers
		86: "🌨️", // Heavy snow showers
		95: "⛈️", // Thunderstorm
		96: "⛈️", // Thunderstorm with slight hail
		99: "⛈️", // Thunderstorm with heavy hail
	}

	if icon, ok := weatherIcons[code]; ok {
		return icon
	}
	return "❓"
}

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		weather, err := getWeather()
		if err != nil {
			fmt.Printf("Error getting weather: %v\n", err)
			return c.Render("index", fiber.Map{
				"Temperature": "N/A",
				"WeatherIcon": "❓",
			})
		}

		return c.Render("index", fiber.Map{
			"Temperature": fmt.Sprintf("%.1f", weather.Current.Temperature),
			"WeatherIcon": getWeatherIcon(weather.Current.WeatherCode),
		})
	})

	app.Listen(":3000")
}
