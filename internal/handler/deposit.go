package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ezjuanify/wallet/internal/handler/request"
	"github.com/ezjuanify/wallet/internal/handler/response"
	"github.com/ezjuanify/wallet/internal/model"
)

func (wh *WalletHandler) DepositResponse(w http.ResponseWriter, r *http.Request) {
	var req request.DepositRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	wallet, err := wh.walletService.DoDeposit(ctx, req.Username, req.Amount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Deposit Error: %v", err), http.StatusInternalServerError)
		return
	}

	if err := wh.transactionService.LogTransaction(ctx, req.Username, model.TypeDeposit, req.Amount, nil); err != nil {
		fmt.Printf("Error logging transaction: %v\n", err)
	}

	resp := response.TransactionResponse{
		Status:          http.StatusOK,
		TransactionType: model.TypeDeposit,
		Wallet:          *wallet,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
	}
}
