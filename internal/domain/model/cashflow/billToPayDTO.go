package cashflow

import "time"

type BillToPayDTO struct {
	ID          uint
	Description string
	Value       float64
	RecordDate  time.Time
}
