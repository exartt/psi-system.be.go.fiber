package repositories

import (
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
	"psi-system.be.go.fiber/internal/domain/model/person"
	"psi-system.be.go.fiber/internal/domain/utils"
	"time"
)

type TransactionRepository interface {
	CreateBill(billToReceive *cashflow.CashFlow) error
	ListBillByType(psychologistID uint, typeTransaction enums.TransactionType) ([]cashflow.BillToReceiveTable, error)
	UpdateBill(bill *cashflow.CashFlow) error
	GetByID(id uint64) (*cashflow.CashFlow, error)
	GetStatusBills(psychologistID uint) ([]cashflow.StatusBill, error)
	Delete(id uint64) error
	ListCashFlow(psychologistID uint) ([]cashflow.Table, error)
	GetCashFlowByDate(psychologistID uint, filteredDateInitial time.Time, filteredDateFinal time.Time) ([]cashflow.CashFlow, error)
	ThrowMonthlyReceives() error
}

func (r *cashFlowRepository) CreateBill(billToReceive *cashflow.CashFlow) error {
	result := r.db.Create(billToReceive)
	return result.Error
}

func (r *cashFlowRepository) ListBillByType(psychologistID uint, typeTransaction enums.TransactionType) ([]cashflow.BillToReceiveTable, error) {
	var billTables []cashflow.BillToReceiveTable
	result := r.db.Table("cash_flows").
		Select("cash_flows.id_fluxo_caixa as ID, cash_flows.flu_valor as Value, cash_flows.descricao as Description,"+
			" cash_flows.status as Status, people.name as \"PatientName\", cash_flows.flu_data_registro as \"RecordDate\"").
		Joins("LEFT JOIN patients ON cash_flows.id_paciente = patients.id").
		Joins("LEFT JOIN people ON patients.person_id = people.id").
		Where("cash_flows.tipo_transacao = ? AND cash_flows.id_psicologo = ?", typeTransaction.String(), psychologistID).
		Scan(&billTables)

	print(result.Error)
	if result != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []cashflow.BillToReceiveTable{}, result.Error
	}
	return billTables, nil
}

func (r *cashFlowRepository) UpdateBill(bill *cashflow.CashFlow) error {
	result := r.db.Save(bill)
	return result.Error
}

func (r *cashFlowRepository) Delete(id uint64) error {
	result := r.db.Delete(&cashflow.CashFlow{}, id)
	return result.Error
}

func (r *cashFlowRepository) ThrowMonthlyReceives() error {
	var patientsMonthly []person.PatientsMonthly

	result := r.db.Table("patients").Select("patients.id as id, patients.psychologist_id as psychologist_id, patients.session_price as session_price").
		Joins("LEFT JOIN people p ON p.id = patients.person_id").
		Where("patients.is_plan = true AND p.is_active = true").
		Where("NOT EXISTS (SELECT 1 FROM cash_flows WHERE cash_flows.id_paciente = patients.id AND cash_flows.flu_data_registro >= ?)", utils.GetFifthDayOfCurrentMonth()).
		Scan(&patientsMonthly)

	if result != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Info("No patients found to insert receives")
	}

	for _, patient := range patientsMonthly {
		insertBill := cashflow.CashFlow{
			PatientID:       patient.PatientID,
			PsychologistID:  patient.PsychologistID,
			Value:           patient.Value,
			Description:     "Lançamento automático mensal",
			Status:          "PENDING",
			TransactionType: "RECEIVE",
			RecordDate:      time.Now(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		err := r.CreateBill(&insertBill)
		if err != nil {
			log.Error(err.Error())
		}
	}

	return result.Error
}
