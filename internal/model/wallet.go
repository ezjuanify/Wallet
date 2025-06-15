package model

import (
	"time"
)

type Wallet struct {
	Username            string     `json:"username"`
	Balance             int64      `json:"balance"`
	LastDepositAmount   *int64     `json:"lastDepositAmount"`
	LastDepositUpdated  *time.Time `json:"lastDepositUpdated"`
	LastWithdrawAmount  *int64     `json:"lastWithdrawAmount"`
	LastWithdrawUpdated *time.Time `json:"lastWithdrawUpdated"`
}
