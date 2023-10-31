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
	personRepo := repositories.NewPersonRepository(infrastructure.DB)
	billToReceiveHandler := provideBillToReceive(personRepo)

	// scheduler
	schedule.Post("/create-appointment", appointmentHandler.CreateAppointment)
	schedule.Get("/list-appointments", appointmentHandler.GetAppointmentsByYear)
	schedule.Put("/update-appointment/:id/status/:status", appointmentHandler.UpdateAppointment)
	schedule.Put("/update-appointment/:id", appointmentHandler.UpdateAppointment)
	schedule.Put("/update-appointment/:id/cancel-appointment", appointmentHandler.CancelAppointment)
	schedule.Post("/confirm-appointment", billToReceiveHandler.CreateBill)

	// bill
	billToReceive := app.Group("/transactions/v1")
	billToReceive.Get("/list-bill", billToReceiveHandler.ListBillByType)
	billToReceive.Post("/create-bill", billToReceiveHandler.CreateBill)
	billToReceive.Get("/get-bill/:id", billToReceiveHandler.GetByID)
	billToReceive.Put("/update-bill", billToReceiveHandler.UpdateBill)
	billToReceive.Delete("/delete-bill/:id", billToReceiveHandler.Delete)
	billToReceive.Put("/update-bill/:id/confirm", billToReceiveHandler.ConfirmPaymentBill)
	billToReceive.Put("/update-bill/:id/remove", billToReceiveHandler.RemoveConfirmationPaymentBill)

	calendar := app.Group("/calendar/v1")
	googleConsumerHandler := provideGoogleConsumer()
	calendar.Get("/google-authenticate", googleConsumerHandler.RequestGoogleAuth)
	calendar.Post("/google-authenticate/callback", googleConsumerHandler.HandleGoogleCallback)

	patient := app.Group("/patient/v1")
	patientHandler := providePatient(personRepo)
	patient.Get("/list-patients", patientHandler.GetPatientsOptions)
	patient.Post("/create-patient", patientHandler.CreatePatient)
	patient.Put("/update-patient", patientHandler.UpdatePatient)
	patient.Delete("/delete-patient/:id", patientHandler.DeactivatePatient)
	patient.Get("/list-person-patient", patientHandler.GetPersonPatient)
	patient.Get("/get-patient/:id", patientHandler.GetPatient)

	psychologist := app.Group("/psychologist/v1")
	psychologistHandler := providePsychologist(personRepo)
	psychologist.Post("/create-psychologist", psychologistHandler.CreatePsychologist)

}

func providePsychologist(personRepo repositories.IPersonRepository) *handlers.PsychologistHandler {
	psychologistRepo := repositories.NewPsychologistRepository(infrastructure.DB)
	psychologistService := services.NewPsychologistService(personRepo, psychologistRepo)
	return handlers.NewPsychologistHandler(psychologistService)
}

func providePatient(personRepo repositories.IPersonRepository) *handlers.PatientHandler {
	patientRepo := repositories.NewPatientRepository(infrastructure.DB)
	patientService := services.NewPatientService(patientRepo, personRepo)
	return handlers.NewPatientHandler(patientService)
}

func provideGoogleConsumer() *handlers.GoogleConsumerHandler {
	iGConsumer := repositories.NewGCalendarRepository(infrastructure.DB)
	consumer := services.NewGCalendarService(iGConsumer)
	return handlers.NewGoogleConsumerHandler(consumer)
}

func provideAppointment() *handlers.AppointmentHandler {
	appointmentRepo := repositories.NewAppointmentRepository(infrastructure.DB)
	appointmentService := services.NewAppointmentService(appointmentRepo)
	return handlers.NewAppointmentHandler(appointmentService)
}

func provideBillToReceive(personRepo repositories.IPersonRepository) *handlers.BillToReceiveHandler {
	transactionsRepo := repositories.NewCashFlowRepository(infrastructure.DB)
	patientRepo := repositories.NewPatientRepository(infrastructure.DB)
	patientService := services.NewPatientService(patientRepo, personRepo)
	appointmentRepo := repositories.NewAppointmentRepository(infrastructure.DB)
	appointmentService := services.NewAppointmentService(appointmentRepo)
	billToReceiveService := services.NewBillToReceiveService(transactionsRepo, patientService, appointmentService)
	return handlers.NewBillToReceiveHandler(billToReceiveService)
}
