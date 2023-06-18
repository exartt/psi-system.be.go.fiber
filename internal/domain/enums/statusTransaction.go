package enums

type StatusTransaction uint

const (
	PENDING StatusTransaction = iota
	PAID
	OVERDUE
	CANCELED
	REFUNDED
)

func (s StatusTransaction) String() string {
	switch s {
	case PAID:
		return "PAID"
	case OVERDUE:
		return "OVERDUE"
	case CANCELED:
		return "CANCELED"
	case REFUNDED:
		return "REFUNDED"
	default:
		return "PENDING"
	}
}
