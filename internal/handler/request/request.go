package request

type RequestPayload struct {
	Username     string  `json:"username"`
	Amount       int64   `json:"amount"`
	Counterparty *string `json:"counterparty,omitempty"`
}
