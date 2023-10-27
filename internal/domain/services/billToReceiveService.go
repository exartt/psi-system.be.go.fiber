package services

import (
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
	"psi-system.be.go.fiber/internal/domain/model/person"
	"psi-system.be.go.fiber/internal/repositories"
	"time"
)

type BillToReceiveService interface {
	CreateBillToReceive(billDTO *cashflow.BillToPayDTO) error
	ListBillToReceive() ([]*cashflow.CashFlow, error)
	UpdateBillToReceive(billDTO *cashflow.BillToPayDTO) error
}

type billToReceiveService struct {
	repo               repositories.BillToReceiveRepository
	patientService     PatientService
	appointmentService AppointmentService
}

func NewBillToReceiveService(repo repositories.BillToReceiveRepository, patientService PatientService, appointmentService AppointmentService) BillToReceiveService {
	return &billToReceiveService{repo: repo, patientService: patientService, appointmentService: appointmentService}
}

func (s *billToReceiveService) ListBillToReceive() ([]*cashflow.CashFlow, error) {
	return s.repo.ListBillToReceive()
}

func (s *billToReceiveService) CreateBillToReceive(billDTO *cashflow.BillToPayDTO) error {
	patient, err := s.patientService.GetByID(billDTO.PatientID)
	if err != nil {
		return logAndReturnError("CreateBillToReceive", "Error fetching patient", err)
	}

	value, err := calculateValue(patient)
	if err != nil {
		return logAndReturnError("CreateBillToReceive", "Error retrieving the session price", err)
	}

	if value != 0 {
		cashFlow, err := prepareCashFlow(value, billDTO)
		if err != nil {
			return logAndReturnError("CreateBillToReceive", "Error preparing cash flow", err)
		}

		if err = s.repo.CreateBillToReceive(cashFlow); err != nil {
			return logAndReturnError("CreateBillToReceive", "Database error", err)
		}

		logInfo("CreateBillToReceive", "Bill to receive saved successfully")

	}

	if err = s.appointmentService.UpdateStatusAppointment(billDTO.AppointmentID, enums.Atendido); err != nil {
		return logAndReturnError("CreateBillToReceive", "Error updating appointment status", err)
	}

	logInfo("CreateBillToReceive", "Appointment status updated successfully")

	return nil
}

func (s *billToReceiveService) UpdateBillToReceive(billDTO *cashflow.BillToPayDTO) error {
	// Fetch existing bill to receive by ID
	existingBill, err := s.repo.GetByID(uint64(billDTO.ID))
	if err != nil {
		return logAndReturnError("UpdateBillToReceive", "Error fetching existing bill", err)
	}
	if existingBill == nil {
		return logAndReturnError("UpdateBillToReceive", "Bill not found", err)
	}

	existingBill.Description = billDTO.Description

	if err = s.repo.UpdateBillToReceive(existingBill); err != nil {
		return logAndReturnError("UpdateBillToReceive", "Database error", err)
	}

	logInfo("UpdateBillToReceive", "Bill to receive updated successfully")
	return nil
}

func prepareCashFlow(value float64, billDTO *cashflow.BillToPayDTO) (*cashflow.CashFlow, error) {

	return &cashflow.CashFlow{
		PsychologistID:  billDTO.PsychologistID,
		PatientId:       billDTO.PatientID,
		AppointmentID:   billDTO.AppointmentID,
		TransactionType: enums.RECEIVABLE.String(),
		Value:           value,
		Description:     billDTO.Description,
		RecordDate:      time.Now(),
		Status:          enums.PENDING.String(),
	}, nil
}

func calculateValue(patient *person.Patient) (float64, error) {
	if !patient.IsPlan {
		return patient.SessionPrice, nil
	}
	return 0, nil
}
