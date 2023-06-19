package person

import "time"

type Person struct {
	ID               uint `gorm:"primary_key"`
	Name             string
	Email            string `gorm:"unique;not null"`
	CellPhone        string
	Phone            string
	ZipCode          string
	Address          string
	IsActive         bool `gorm:"default:true"`
	CPF              string
	RG               string
	RegistrationDate time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Psychologist     *Psychologist
	Patient          *Patient
}
