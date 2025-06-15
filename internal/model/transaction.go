package model

import (
	"time"
)

type Transaction struct {
	ID           int64
	Username     string
	Type         string
	Amount       int64
	Counterparty *string
	Timestamp    time.Time
	Hash         string
}

type TransactionType string

const (
	TypeDeposit     TransactionType = "deposit"
	TypeWithdraw    TransactionType = "withdraw"
	TypeTransferIn  TransactionType = "transfer_in"
	TypeTransferOut TransactionType = "transfer_out"
)
