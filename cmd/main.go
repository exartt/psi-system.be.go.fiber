package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"psi-system.be.go.fiber/internal/delivery/http"
	"psi-system.be.go.fiber/internal/infrastructure"
	"psi-system.be.go.fiber/internal/infrastructure/database"
	"psi-system.be.go.fiber/pkg"
)

func main() {
	pkg.LoadEnv()
	infrastructure.ConnectDB()
	database.Migrate()

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())

	http.RegisterRoutes(app)
	log.Fatal(app.Listen(":3030"))
}
