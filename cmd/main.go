package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
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
	//app.Use(cors.New(cors.Config{
	//	AllowCredentials: true,
	//}))
	http.RegisterRoutes(app)
	log.Fatal(app.Listen(":3020"))
}
