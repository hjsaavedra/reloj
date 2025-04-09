package main

import (
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Motor de plantillas HTML
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Ruta principal
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	
    // Usa el puerto asignado por Railway
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000" // por si estás en local
    }

    app.Listen(":" + port)
}
