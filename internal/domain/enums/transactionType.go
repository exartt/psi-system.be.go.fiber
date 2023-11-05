package enums

type TransactionType uint

const (
	RECEIVABLE TransactionType = iota
	PAYABLE
	CASHFLOW
	CONFIRM_PAYMENT
)

func (s TransactionType) String() string {
	switch s {
	case PAYABLE:
		return "PAYABLE"
	case CASHFLOW:
		return "CASHFLOW"
	case CONFIRM_PAYMENT:
		return "CONFIRM_PAYMENT"
	default:
		return "RECEIVABLE"
	}
}
