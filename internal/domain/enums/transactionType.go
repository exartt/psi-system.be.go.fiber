package enums

type TransactionType uint

const (
	RECEIVABLE TransactionType = iota
	PAYABLE
)

func (s TransactionType) String() string {
	switch s {
	case PAYABLE:
		return "PAYABLE"
	default:
		return "RECEIVABLE"
	}
}
