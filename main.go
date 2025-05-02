package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

const (
	weatherAPIKey      = "b2303b8f32f74b5fbf8152337250205"
	location           = "Santo Domingo, Chile"
	weatherUpdateEvery = 3 * time.Hour
)

type WeatherAPIResponse struct {
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
		} `json:"condition"`
	} `json:"current"`
}

var (
	weatherData     = make(map[string]interface{})
	weatherDataLock sync.RWMutex
)

func main() {
	// Configurar el motor de plantillas
	engine := html.New("./views", ".html")

	// Configurar la app Fiber
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Iniciar la actualización periódica del clima
	go updateWeatherPeriodically()

	// Ruta principal
	app.Get("/", func(c *fiber.Ctx) error {
		now := time.Now()

		weatherDataLock.RLock()
		defer weatherDataLock.RUnlock()

		// Si no hay datos de clima (primer inicio), obtener algunos
		if len(weatherData) == 0 {
			clima, iconoTemperatura, iconoClima, err := getClimaReal()
			if err != nil {
				clima = "Datos no disponibles"
				iconoTemperatura = getIconoClimaSimulado()
				iconoClima = ""
			}
			weatherData["Clima"] = clima
			weatherData["IconoTemperatura"] = iconoTemperatura
			weatherData["IconoClima"] = iconoClima
		}

		// Datos para pasar a la plantilla
		data := map[string]interface{}{
			"Dia":              formatDia(now),
			"Fecha":            formatFecha(now),
			"Hora":             formatHora(now),
			"Clima":            weatherData["Clima"],
			"IconoTemperatura": weatherData["IconoTemperatura"],
			"IconoClima":       weatherData["IconoClima"],
		}

		return c.Render("index", data)
	})

	// Ruta API para obtener datos actualizados del clima
	app.Get("/api/weather", func(c *fiber.Ctx) error {
		weatherDataLock.RLock()
		defer weatherDataLock.RUnlock()
		return c.JSON(weatherData)
	})

	// Iniciar el servidor
	fmt.Println("Servidor iniciado en http://localhost:3000")
	app.Listen(":3000")
}

// updateWeatherPeriodically actualiza los datos del clima cada 3 horas
func updateWeatherPeriodically() {
	for {
		clima, iconoTemperatura, iconoClima, err := getClimaReal()
		if err != nil {
			fmt.Printf("Error actualizando clima: %v\n", err)
			// Usar datos simulados como fallback
			clima = "Datos no disponibles"
			iconoTemperatura = getIconoClimaSimulado()
			iconoClima = ""
		}

		weatherDataLock.Lock()
		weatherData["Clima"] = clima
		weatherData["IconoTemperatura"] = iconoTemperatura
		weatherData["IconoClima"] = iconoClima
		weatherDataLock.Unlock()

		time.Sleep(weatherUpdateEvery)
	}
}

// getClimaReal obtiene datos actuales del clima desde WeatherAPI
func getClimaReal() (string, string, string, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=es",
		weatherAPIKey, url.QueryEscape(location))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", "", "", fmt.Errorf("error de conexión: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", "", fmt.Errorf("error API (código %d)", resp.StatusCode)
	}

	var data WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", "", fmt.Errorf("error decodificando respuesta: %v", err)
	}

	clima := data.Current.Condition.Text
	return clima, "🌡️" + fmt.Sprintf(" %.1f°c", data.Current.TempC), data.Current.Condition.Icon, nil
}

// formatDia devuelve la fecha formateada en español
func formatDia(t time.Time) string {
	dias := []string{"Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"}
	return dias[t.Weekday()]
}

func formatFecha(t time.Time) string {
	meses := []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
	return fmt.Sprintf("%d de %s", t.Day(), meses[t.Month()-1])
}

// formatHora devuelve la hora en formato HH:MM
func formatHora(t time.Time) string {
	return t.Format("15:04")
}

// getIconoClimaSimulado devuelve un icono aleatorio (usado como fallback)
func getIconoClimaSimulado() string {
	temperaturas := []int{18, 19, 20, 21, 22, 23, 24, 25}
	temp := temperaturas[rand.Intn(len(temperaturas))]

	iconos := []string{"☀️", "⛅", "☁️", "🌧️"}
	return iconos[rand.Intn(len(iconos))] + " " + fmt.Sprintf("%d°C", temp)

}
