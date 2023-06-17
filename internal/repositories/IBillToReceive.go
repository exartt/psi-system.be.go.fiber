package repositories

import (
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
)

type BillToReceiveRepository interface {
	CreateBillToReceive(billToReceive *cashflow.BillToReceive) error
}

type billToReceiveRepository struct {
	db *gorm.DB
}

func NewBillToReceiveRepository(db *gorm.DB) BillToReceiveRepository {
	return &billToReceiveRepository{db: db}
}

func (r *billToReceiveRepository) CreateBillToReceive(billToReceive *cashflow.BillToReceive) error {
	result := r.db.Create(billToReceive)
	return result.Error
}
