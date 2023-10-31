package cashflow

import "time"

type BillDTO struct {
	ID              uint
	PatientID       uint
	AppointmentID   uint
	Value           float64
	Description     string
	RecordDate      string
	TransactionType string
}

type BillToReceiveTable struct {
	ID          uint
	Value       float64
	PatientName string
	Description string
	Status      string
	RecordDate  time.Time
}
