package repositories

import (
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/model/person"
)

type IPsychologistRepository interface {
	CreatePsychologist(psychologist *person.Psychologist) error
}

func NewPsychologistRepository(db *gorm.DB) IPsychologistRepository {
	return &psychologistRepository{db: db}
}

type psychologistRepository struct {
	db *gorm.DB
}

func (r *psychologistRepository) CreatePsychologist(psychologist *person.Psychologist) error {
	result := r.db.Create(psychologist)
	return result.Error
}
