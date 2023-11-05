package repositories

import (
	"errors"
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
	"time"
)

type TransactionRepository interface {
	CreateBill(billToReceive *cashflow.CashFlow) error
	ListBillByType(psychologistID uint, typeTransaction enums.TransactionType) ([]cashflow.BillToReceiveTable, error)
	UpdateBill(bill *cashflow.CashFlow) error
	GetByID(id uint64) (*cashflow.CashFlow, error)
	Delete(id uint64) error
	ListCashFlow(psychologistID uint) ([]cashflow.Table, error)
	GetCashFlowByDate(psychologistID uint, filteredDateInitial time.Time, filteredDateFinal time.Time) ([]cashflow.CashFlow, error)
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
