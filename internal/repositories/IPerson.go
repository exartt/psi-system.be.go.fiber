package repositories

import (
	"gorm.io/gorm"
	"psi-system.be.go.fiber/internal/domain/model/person"
)

type IPersonRepository interface {
	CreatePerson(person *person.Person) error
}

type personRepository struct {
	db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) IPersonRepository {
	return &personRepository{db: db}
}

func (r *personRepository) CreatePerson(person *person.Person) error {
	result := r.db.Create(person)
	return result.Error
}
