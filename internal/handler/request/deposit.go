package request

type DepositRequest struct {
	Username string `json:"username"`
	Amount   int64  `json:"amount"`
}
