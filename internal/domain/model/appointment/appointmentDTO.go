package appointment

import (
	"psi-system.be.go.fiber/internal/domain/enums"
	"time"
)

type DTO struct {
	ID             uint
	PsychologistID uint
	PatientID      uint
	EventID        string
	Start          time.Time
	End            time.Time
	Summary        string
	Description    string
	Status         enums.StatusAgendamento
}
