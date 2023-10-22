package person

import "time"

type PsychologistDto struct {
	ID          uint
	UserID      uint
	PersonID    uint
	Email       string
	Password    string
	IsActive    bool
	ProfileType uint

	Name             string
	CellPhone        string
	Phone            string
	ZipCode          string
	Address          string
	CPF              string
	RG               string
	RegistrationDate time.Time

	Access   int
	TenantID uint

	// Common Fields
	CreatedAt time.Time
	UpdatedAt time.Time
}
