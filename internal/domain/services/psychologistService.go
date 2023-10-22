package services

import (
	"psi-system.be.go.fiber/internal/domain/model/person"
	"psi-system.be.go.fiber/internal/repositories"
)

type PsychologistService interface {
	SavePsychologist(data *person.PsychologistDto) error
}

type psychologistService struct {
	repo    repositories.IPersonRepository
	psyRepo repositories.IPsychologistRepository
}

func NewPsychologistService(repo repositories.IPersonRepository, psyRepo repositories.IPsychologistRepository) PsychologistService {
	return &psychologistService{
		repo:    repo,
		psyRepo: psyRepo,
	}
}

func (s *psychologistService) SavePsychologist(data *person.PsychologistDto) error {
	personDto := &person.Person{
		Name:             data.Name,
		Email:            data.Email,
		CellPhone:        data.CellPhone,
		Phone:            data.Phone,
		ZipCode:          data.ZipCode,
		Address:          data.Address,
		IsActive:         data.IsActive,
		CPF:              data.CPF,
		RG:               data.RG,
		RegistrationDate: data.RegistrationDate,
	}

	if err := s.repo.CreatePerson(personDto); err != nil {
		return err
	}

	data.PersonID = personDto.ID

	psychologist := &person.Psychologist{
		PersonID: data.PersonID,
		Access:   data.Access,
		TenantID: data.TenantID,
		UserID:   data.UserID,
	}

	if err := s.psyRepo.CreatePsychologist(psychologist); err != nil {
		return err
	}

	return nil
}
