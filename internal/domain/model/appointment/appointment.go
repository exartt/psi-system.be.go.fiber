package appointment

import "time"

type Appointment struct {
	ID             uint `gorm:"primary_key"`
	PsychologistID uint // Add this
	PatientID      uint // Add this
	TenantID       uint `gorm:"not null"`
	CalendarID     string
	Start          time.Time
	End            time.Time
	Summary        string
	Description    string
	Location       string
	Status         string
	Notify         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
