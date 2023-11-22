package services

import (
	"errors"
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
	"psi-system.be.go.fiber/internal/domain/model/person"
	"psi-system.be.go.fiber/internal/domain/utils"
	"psi-system.be.go.fiber/internal/repositories"
	"time"
)

type TransactionService interface {
	CreateBill(billDTO *cashflow.BillDTO, psychologistID uint) error
	ListBillByType(psychologistID uint, typeTransaction enums.TransactionType) ([]cashflow.BillToReceiveTable, error)
	UpdateBill(billDTO *cashflow.BillDTO) error
	StatusUpdateBill(billDTO *cashflow.CashFlow) error
	GetByID(id uint64) (*cashflow.CashFlow, error)
	Delete(id uint64) error
	ListCashFlow(psychologistID uint) ([]cashflow.Table, error)
	GetCashFlowByDate(psychologistID uint, filteredDateInitial time.Time, filteredDateFinal time.Time) ([]cashflow.CashFlow, error)
	GetStatusBills(psychologistID uint) ([]cashflow.StatusBill, error)
}

type billToReceiveService struct {
	repo               repositories.TransactionRepository
	patientService     PatientService
	appointmentService AppointmentService
}

func NewBillToReceiveService(repo repositories.TransactionRepository, patientService PatientService, appointmentService AppointmentService) TransactionService {
	return &billToReceiveService{repo: repo, patientService: patientService, appointmentService: appointmentService}
}

func (s *billToReceiveService) ListBillByType(psychologistID uint, typeTransaction enums.TransactionType) ([]cashflow.BillToReceiveTable, error) {
	return s.repo.ListBillByType(psychologistID, typeTransaction)
}

func (s *billToReceiveService) CreateBill(billDTO *cashflow.BillDTO, psychologistID uint) error {
	transactionType, err := utils.CastToTransactioTypeEnum(billDTO.TransactionType)
	if err != nil {
		return logAndReturnError("CreateBill", "Ocorreu um erro inesperado, por favor, tente novamente mais tarde.", err)
	}
	switch transactionType {
	case enums.RECEIVABLE:
		patient, err := s.patientService.GetByID(billDTO.PatientID)
		if err != nil {
			return logAndReturnError("CreateBill", "Error fetching patient", err)
		}

		value, err := calculateValue(patient)
		if err != nil {
			return logAndReturnError("CreateBill", "Error retrieving the session price", err)
		}

		if value != 0 {
			cashFlow, err := prepareCashFlow(value, billDTO, psychologistID, transactionType)
			if err != nil {
				return logAndReturnError("CreateBill", "Error preparing cash flow", err)
			}

			if err = s.repo.CreateBill(cashFlow); err != nil {
				return logAndReturnError("CreateBill", "Database error", err)
			}

			logInfo("CreateBill", "Bill to receive saved successfully")
		}
		if billDTO.AppointmentID != 0 {
			if err = s.appointmentService.UpdateStatusAppointment(billDTO.AppointmentID, enums.Atendido); err != nil {
				return logAndReturnError("CreateBill", "Error updating appointment status", err)
			}
			logInfo("CreateBill", "Appointment status updated successfully")
		}
		return nil
	case enums.PAYABLE:
		cashFlow, err := prepareCashFlow(billDTO.Value, billDTO, psychologistID, transactionType)
		if err != nil {
			return logAndReturnError("CreateBill", "Error preparing cash flow", err)
		}

		if err = s.repo.CreateBill(cashFlow); err != nil {
			return logAndReturnError("CreateBill", "Database error", err)
		}
		return nil
	default:
		return errors.New("Operação não disponível, para adicionar um novo registro no fluxo de caixa faça através do \"contas a pagar\" ou \"contas a receber\".")
	}
}

func (s *billToReceiveService) UpdateBill(billDTO *cashflow.BillDTO) error {

	existingBill, err := s.repo.GetByID(uint64(billDTO.ID))
	if err != nil {
		return logAndReturnError("UpdateBill", "Error fetching existing bill", err)
	}
	if existingBill == nil {
		return logAndReturnError("UpdateBill", "Bill not found", err)
	}

	existingBill.Description = billDTO.Description
	existingBill.Value = billDTO.Value
	existingBill.RecordDate, err = time.Parse("2006/01/02", billDTO.RecordDate)

	if err = s.repo.UpdateBill(existingBill); err != nil {
		return logAndReturnError("UpdateBill", "Database error", err)
	}

	logInfo("UpdateBill", "Bill to receive updated successfully")
	return nil
}

func (s *billToReceiveService) StatusUpdateBill(billToPersist *cashflow.CashFlow) error {

	if err := s.repo.UpdateBill(billToPersist); err != nil {
		return logAndReturnError("UpdateStatusBill", "Database error", err)
	}

	logInfo("UpdateBill", "Bill to receive updated successfully")
	return nil
}

func prepareCashFlow(value float64, billDTO *cashflow.BillDTO, psychologistID uint, transactionType enums.TransactionType) (*cashflow.CashFlow, error) {
	var parsedTime time.Time

	if billDTO.RecordDate != "" {
		var err error
		parsedTime, err = time.Parse("2006/01/02", billDTO.RecordDate)
		if err != nil {
			return nil, errors.New("Erro ao analisar a data. Verifique se foi inserida corretamente.")
		}
	} else {
		parsedTime = time.Now()
	}

	return &cashflow.CashFlow{
		PsychologistID:  psychologistID,
		PatientID:       billDTO.PatientID,
		AppointmentID:   billDTO.AppointmentID,
		TransactionType: transactionType.String(),
		Value:           value,
		Description:     billDTO.Description,
		RecordDate:      parsedTime,
		Status:          enums.PENDING.String(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

func calculateValue(patient *person.Patient) (float64, error) {
	if !patient.IsPlan {
		return patient.SessionPrice, nil
	}
	return 0, nil
}

func (s *billToReceiveService) GetByID(id uint64) (*cashflow.CashFlow, error) {
	return s.repo.GetByID(id)
}

func (s *billToReceiveService) Delete(id uint64) error {
	return s.repo.Delete(id)
}

func (s *billToReceiveService) ListCashFlow(psychologistID uint) ([]cashflow.Table, error) {
	return s.repo.ListCashFlow(psychologistID)
}

func (s *billToReceiveService) GetCashFlowByDate(psychologistID uint, filteredDateInitial time.Time, filteredDateFinal time.Time) ([]cashflow.CashFlow, error) {
	return s.repo.GetCashFlowByDate(psychologistID, filteredDateInitial, filteredDateFinal)
}

func (s *billToReceiveService) GetStatusBills(psychologistID uint) ([]cashflow.StatusBill, error) {
	statusBills, err := s.repo.GetStatusBills(psychologistID)
	return statusBills, err
}
