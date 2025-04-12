package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	// Open-Meteo no requiere API key
	url := "https://api.open-meteo.com/v1/forecast?latitude=-33.4372&longitude=-70.6506&current=temperature_2m,weather_code"
	log.Printf("Haciendo petición a: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error en la petición HTTP: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error leyendo el cuerpo de la respuesta: %v", err)
		return nil, err
	}

	log.Printf("Respuesta de la API: %s", string(body))

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		log.Printf("Error decodificando JSON: %v", err)
		return nil, err
	}

	return &weather, nil
}

func getWeatherIcon(code int) string {
	// Mapeo de códigos de clima a iconos
	icons := map[int]string{
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

	if icon, ok := icons[code]; ok {
		return icon
	}
	return "❓" // Icono por defecto si no se encuentra el código
}

func main() {
	// Motor de plantillas HTML
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Ruta principal
	app.Get("/", func(c *fiber.Ctx) error {
		weather, err := getWeather()
		if err != nil {
			log.Printf("Error obteniendo el clima: %v", err)
			return c.Render("index", fiber.Map{
				"Temperature": "N/A",
				"WeatherIcon": "❓",
			})
		}

		icon := getWeatherIcon(weather.Current.WeatherCode)
		log.Printf("Temperatura: %.1f°C, Código del clima: %d", weather.Current.Temperature, weather.Current.WeatherCode)
		return c.Render("index", fiber.Map{
			"Temperature": fmt.Sprintf("%.1f", weather.Current.Temperature),
			"WeatherIcon": icon,
		})
	})

	log.Println("Servidor iniciado en http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}
