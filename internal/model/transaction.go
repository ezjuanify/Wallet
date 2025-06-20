package model

import (
	"time"
)

type Transaction struct {
	ID           int64     `json:"ID"`
	Username     string    `json:"username"`
	TxnType      TxnType   `json:"txnType"`
	Amount       int64     `json:"amount"`
	Counterparty *string   `json:"counterparty"`
	Timestamp    time.Time `json:"timestamp"`
	Hash         string    `json:"hash"`
}

type TxnType string

const (
	TypeDeposit     TxnType = "deposit"
	TypeWithdraw    TxnType = "withdraw"
	TypeTransfer    TxnType = "transfer"
	TypeTransferIn  TxnType = "transfer_in"
	TypeTransferOut TxnType = "transfer_out"
)

var txnTypes = map[TxnType]struct{}{
	TypeDeposit:     {},
	TypeWithdraw:    {},
	TypeTransfer:    {},
	TypeTransferIn:  {},
	TypeTransferOut: {},
}

func IsTxnTypeValid(txnType string) bool {
	_, ok := txnTypes[TxnType(txnType)]
	return ok
}
