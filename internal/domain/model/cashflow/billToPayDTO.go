package cashflow

type BillToPayDTO struct {
	ID             uint
	PsychologistID uint
	PatientID      uint
	AppointmentID  uint
	Description    string
}
