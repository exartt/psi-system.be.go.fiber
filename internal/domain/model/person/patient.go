package person

import "psi-system.be.go.fiber/internal/domain/model/appointment"

type Patient struct {
	ID             uint `gorm:"primary_key"`
	PersonID       uint `gorm:"not null"`
	PsychologistID uint `gorm:"not null"`
	IsPlan         bool
	SessionPrice   float64
	ConversionType string
	Appointments   []appointment.Appointment `gorm:"foreignKey:PatientID"`
}

type Option struct {
	Value        string `json:"value"`
	Label        string `json:"label"`
	SessionPrice string `json:"session_price"`
}

type PatientsMonthly struct {
	PatientID      uint    `gorm:"column:id;"`
	PsychologistID uint    `gorm:"column:psychologist_id"`
	Value          float64 `gorm:"column:session_price"`
}
