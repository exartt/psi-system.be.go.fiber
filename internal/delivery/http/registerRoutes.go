package http

import (
	"github.com/gofiber/fiber/v2"
	"psi-system.be.go.fiber/internal/domain/handlers"
	"psi-system.be.go.fiber/internal/domain/services"
	"psi-system.be.go.fiber/internal/infrastructure"
	"psi-system.be.go.fiber/internal/repositories"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/schedule/v1")
	handler := provide()
	api.Post("/create-appointment", handler.CreateAppointment)
}

func provide() *handlers.AppointmentHandler {
	appointmentRepo := repositories.NewAppointmentRepository(infrastructure.DB)
	appointmentService := services.NewAppointmentService(appointmentRepo)
	return handlers.NewAppointmentHandler(appointmentService)
}
