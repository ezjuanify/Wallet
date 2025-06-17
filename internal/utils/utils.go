package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func PtrInt64(v int64) *int64 { return &v }

func GenerateTransactionHash(txUser string, txType string, txAmount int64, txCounterparty *string, timestamp string) string {
	var counterparty string
	if txCounterparty != nil {
		counterparty = *txCounterparty
	}
	raw := fmt.Sprintf("%s|%s|%d|%s|%s", txUser, txType, txAmount, counterparty, timestamp)
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}
