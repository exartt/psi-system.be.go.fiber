package repositories

import (
	"errors"
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/enums"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
	"time"
)

type CashFlowRepository interface {
	CreateBill(billToReceive *cashflow.CashFlow) error
	GetByID(id uint64) (*cashflow.CashFlow, error)
	ListBillByType(psychologistID uint, typeTransaction enums.TransactionType) ([]cashflow.BillToReceiveTable, error)
	UpdateBill(billToReceive *cashflow.CashFlow) error
	Delete(id uint64) error
	ListCashFlow(psychologistID uint) ([]cashflow.Table, error)
	GetCashFlowByDate(psychologistID uint, filteredDateInitial time.Time, filteredDateFinal time.Time) ([]cashflow.CashFlow, error)
	GetStatusBills(psychologistID uint) ([]cashflow.StatusBill, error)
	ThrowMonthlyReceives() error
}

type cashFlowRepository struct {
	db *gorm.DB
}

func NewCashFlowRepository(db *gorm.DB) CashFlowRepository {
	return &cashFlowRepository{db: db}
}

func (r *cashFlowRepository) GetByID(id uint64) (*cashflow.CashFlow, error) {
	var bill cashflow.CashFlow
	result := r.db.First(&bill, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	return &bill, nil
}

func (r *cashFlowRepository) ListCashFlow(psychologistID uint) ([]cashflow.Table, error) {
	var cashFlowTable []cashflow.Table
	result := r.db.Table("cash_flows").Where("cash_flows.id_psicologo = ? and cash_flows.status like 'PAID'", psychologistID).Find(&cashFlowTable)

	print(result.Error)
	if result != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []cashflow.Table{}, result.Error
	}

	return cashFlowTable, nil
}

func (r *cashFlowRepository) GetCashFlowByDate(psychologistID uint, filteredDateInitial time.Time, filteredDateFinal time.Time) ([]cashflow.CashFlow, error) {
	var cashFlows []cashflow.CashFlow
	db := r.db
	db = db.Where("id_psicologo = ?", psychologistID)
	if !filteredDateInitial.IsZero() {
		db = db.Where("flu_data_registro >= ?", filteredDateInitial)
	}
	if !filteredDateFinal.IsZero() {
		db = db.Where("flu_data_registro <= ?", filteredDateFinal)
	}
	db = db.Where("status LIKE 'PAID'")
	result := db.Find(&cashFlows)
	if result != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return []cashflow.CashFlow{}, result.Error
	}

	return cashFlows, nil
}

func (r *cashFlowRepository) GetStatusBills(psychologistID uint) ([]cashflow.StatusBill, error) {
	var statusBills []cashflow.StatusBill

	err := r.db.Table("cash_flows").
		Select("CASE WHEN status = 'PENDING' THEN 0 WHEN status = 'PAID' THEN 1 WHEN status = 'OVERDUE' THEN 2 WHEN status = 'CANCELED' THEN 3 WHEN status = 'REFUNDED' THEN 4 ELSE -1 END as status, count(*) as count").
		Where("id_psicologo = ?", psychologistID).
		Group("status").
		Scan(&statusBills)

	if err.Error != nil {
		if !errors.Is(err.Error, gorm.ErrRecordNotFound) {
			return []cashflow.StatusBill{}, nil
		}
		return nil, err.Error
	}

	return statusBills, nil
}
