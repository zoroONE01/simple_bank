package utils

const (
	USD = "USD"
	EUR = "EUR"
	VND = "VND"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, VND:
		return true
	default:
		return false
	}
}
