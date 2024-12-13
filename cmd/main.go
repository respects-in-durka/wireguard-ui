package main

import (
	"github.com/gofiber/fiber/v3"
	"log"
	"wireguard-ui/internal/http"
)

func main() {
	app := fiber.New()
	http.Configure(app)

	log.Fatal(app.Listen(":3000"))
}
