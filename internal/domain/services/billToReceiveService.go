package services

import (
	"psi-system.be.go.fiber/internal/domain/model/cashflow"
	"psi-system.be.go.fiber/internal/repositories"
)

type BillToReceiveService interface {
}

type billToReceiveService struct {
	repo repositories.BillToReceiveRepository
}

func NewCashFlowService(repo repositories.BillToReceiveRepository) BillToReceiveService {
	return &billToReceiveService{repo: repo}
}

func (s *billToReceiveService) createBillToReceive(billToReceive *cashflow.BillToReceive) error {
	return s.repo.CreateBillToReceive(billToReceive)
}
