package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"psi-system.be.go.fiber/config"
	"psi-system.be.go.fiber/internal/delivery/http"
	"psi-system.be.go.fiber/internal/infrastructure"
	"psi-system.be.go.fiber/internal/infrastructure/database"
	"psi-system.be.go.fiber/internal/jobs"
	"psi-system.be.go.fiber/internal/repositories"
	"psi-system.be.go.fiber/pkg"
)

func main() {
	pkg.LoadEnv()
	infrastructure.ConnectDB()
	database.Migrate()

	transactionsRepo := repositories.NewCashFlowRepository(infrastructure.DB)
	jobs.ScheduleJob(transactionsRepo)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())

	app.Use(config.JWTMiddleware())

	http.RegisterRoutes(app)
	log.Fatal(app.Listen(":3030"))
}
