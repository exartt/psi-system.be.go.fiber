package enums

type TransactionType uint

const (
	RECEIVABLE TransactionType = iota
	PAYABLE
	CASHFLOW
)

func (s TransactionType) String() string {
	switch s {
	case PAYABLE:
		return "PAYABLE"
	case CASHFLOW:
		return "CASHFLOW"
	default:
		return "RECEIVABLE"
	}
}
