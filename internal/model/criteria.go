package model

type Criteria struct {
	Username     string  `json:"username,omitempty"`
	Counterparty string  `json:"counterparty,omitempty"`
	TxnType      TxnType `json:"txnType,omitempty"`
	Limit        int     `json:"limit,omitempty"`
}
