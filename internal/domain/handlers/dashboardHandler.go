package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
	"psi-system.be.go.fiber/internal/domain/model/dashboard"
	"psi-system.be.go.fiber/internal/domain/services"
	"psi-system.be.go.fiber/internal/domain/utils"
	"time"
)

type DashboardHandler struct {
	TransactionService services.TransactionService
	PatientService     services.PatientService
	AppointmentService services.AppointmentService
	Logger             *logrus.Logger
}

func NewDashboardHandler(tService services.TransactionService, pService services.PatientService, aService services.AppointmentService) *DashboardHandler {
	return &DashboardHandler{
		TransactionService: tService,
		PatientService:     pService,
		AppointmentService: aService,
	}
}

func (h *DashboardHandler) GetDashboardData(c *fiber.Ctx) error {
	psychologistID, err := utils.GetPsychologistIDFromContext(c)
	var filteredDateInitial, filteredDateFinal time.Time
	var format = "02/01/2006"
	dataInicial := c.Query("dataInicial")
	dataFinal := c.Query("dataFinal")

	if dataInicial != "" {
		filteredDateInitial, err = time.Parse(format, dataInicial)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ocorreu um erro inesperado, por favor, tente novamente mais tarde."})
		}
	} else {
		filteredDateInitial = time.Now().AddDate(-1, 0, 0)
	}

	if dataFinal != "" {
		filteredDateFinal, err = time.Parse(format, dataFinal)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ocorreu um erro inesperado, por favor, tente novamente mais tarde."})
		}
	} else {
		filteredDateFinal = time.Now().AddDate(1, 0, 0)
	}

	cashFlows, err := h.TransactionService.GetCashFlowByDate(psychologistID, filteredDateInitial, filteredDateFinal)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Ocorreu um erro inesperado, por favor, tente novamente mais tarde."})
	}

	dashboardData := CalculateCashFlowValues(cashFlows)

	dashboardData.NewPatients, _ = h.PatientService.CountNewPatients(psychologistID, filteredDateInitial, filteredDateFinal)

	dashboardData.Appointments, _ = h.AppointmentService.CountAppointmentsByDate(psychologistID, filteredDateInitial, filteredDateFinal)

	return c.Status(fiber.StatusOK).JSON(dashboardData)

}
func CalculateCashFlowValues(cashFlows []cashflow.CashFlow) dashboard.Data {
	result := dashboard.Data{}

	for _, cf := range cashFlows {
		if cf.TransactionType == "RECEIVABLE" {
			result.Brute += cf.Value
		} else if cf.TransactionType == "PAYABLE" {
			result.Paid += cf.Value
		}
	}

	result.Liquid = result.Brute - result.Paid
	return result
}
