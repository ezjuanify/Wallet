package validation

import "fmt"

const (
	UPPER_LIMIT = 999999
)

func isAmountTooLow(amount int64) bool {
	return amount <= 0
}

func isAmountTooHigh(amount int64) bool {
	return amount > UPPER_LIMIT
}

func ValidateAmount(amount int64) error {
	switch {
	case isAmountTooLow(amount):
		return fmt.Errorf("amount must be greater than 0")
	case isAmountTooHigh(amount):
		return fmt.Errorf("amount must not exceed %d", UPPER_LIMIT)
	default:
		return nil
	}
}
