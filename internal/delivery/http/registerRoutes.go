package http

import (
	"github.com/gofiber/fiber/v2"
	"psi-system.be.go.fiber/internal/domain/handlers"
	"psi-system.be.go.fiber/internal/domain/services"
	"psi-system.be.go.fiber/internal/infrastructure"
	"psi-system.be.go.fiber/internal/repositories"
)

func RegisterRoutes(app *fiber.App) {
	schedule := app.Group("/schedule/v1")
	appointmentHandler := provideAppointment()
	billToReceiveHandler := provideBillToReceive()
	schedule.Post("/create-appointment", appointmentHandler.CreateAppointment)
	schedule.Get("/list-appointments", appointmentHandler.GetAppointmentsByYear)
	schedule.Put("/update-appointment/:id/status/:status", appointmentHandler.UpdateAppointment)
	schedule.Put("/update-appointment/:id", appointmentHandler.UpdateAppointment)
	schedule.Put("/update-appointment/:id/cancel-appointment", appointmentHandler.CancelAppointment)
	schedule.Post("/confirm-appointment", billToReceiveHandler.CreateBillToReceive)

	calendar := app.Group("/calendar/v1")
	googleConsumerHandler := provideGoogleConsumer()
	calendar.Get("/google-authenticate", googleConsumerHandler.RequestGoogleAuth)

	patient := app.Group("/patient/v1")
	patientHandler := providePatient()
	patient.Get("/list-patients", patientHandler.GetPatientsOptions)
}

func providePatient() *handlers.PatientHandler {
	patientRepo := repositories.NewPatientRepository(infrastructure.DB)
	patientService := services.NewPatientService(patientRepo)
	return handlers.NewPatientHandler(patientService)
}

func provideGoogleConsumer() *handlers.GoogleConsumerHandler {
	gconsumer := services.NewGoogleConsumerService()
	return handlers.NewGoogleConsumerHandler(gconsumer)
}

func provideAppointment() *handlers.AppointmentHandler {
	appointmentRepo := repositories.NewAppointmentRepository(infrastructure.DB)
	appointmentService := services.NewAppointmentService(appointmentRepo)
	return handlers.NewAppointmentHandler(appointmentService)
}

func provideBillToReceive() *handlers.BillToReceiveHandler {
	billToReceiveRepo := repositories.NewCashFlowRepository(infrastructure.DB)
	patientRepo := repositories.NewPatientRepository(infrastructure.DB)
	patientService := services.NewPatientService(patientRepo)
	appointmentRepo := repositories.NewAppointmentRepository(infrastructure.DB)
	appointmentService := services.NewAppointmentService(appointmentRepo)
	billToReceiveService := services.NewBillToReceiveService(billToReceiveRepo, patientService, appointmentService)
	return handlers.NewBillToReceiveHandler(billToReceiveService)
}
