package handler

import (
	"fmt"
	"net/http"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model/response"
	"github.com/ezjuanify/wallet/internal/utils"
	"github.com/ezjuanify/wallet/internal/validation"
	"go.uber.org/zap"
)

func (h *WalletHandler) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	fnName := "WalletHandler.BalanceHandler"

	ctx := r.Context()

	appErrs := validation.NewHandlerErrors()

	defer func() {
		FinalizeTransactionResponse(fnName, nil, w, appErrs)
	}()

	q := r.URL.Query()
	username := q.Get("username")
	logger.Info(fmt.Sprintf("%s - Request received", fnName), zap.String("username", username))

	wallet, appErr := h.walletService.DoFetchWallet(ctx, username)
	if appErr != nil {
		appErr.Status = http.StatusInternalServerError
		appErrs.AddError(*appErr)
		return
	}
	logger.Info(fmt.Sprintf("%s - Wallet fetched successfully", fnName), zap.Any("wallet", wallet))

	resp := &response.WalletResponse{
		Status: http.StatusOK,
		Wallet: wallet,
	}
	logger.Info(fmt.Sprintf("%s - Sending wallet response", fnName), zap.Any("wallet", wallet))
	SendJSONResponse(fnName, w, resp.Status, resp)
}

func (h *WalletHandler) AdminBalanceHandler(w http.ResponseWriter, r *http.Request) {
	fnName := "WalletHandler.AdminBalanceHandler"

	ctx := r.Context()

	appErrs := validation.NewHandlerErrors()

	defer func() {
		FinalizeTransactionResponse(fnName, nil, w, appErrs)
	}()

	wallets, appErr := h.walletService.DoFetchAllWallets(ctx)
	if appErr != nil {
		appErr.Status = http.StatusInternalServerError
		appErrs.AddError(*appErr)
		return
	}
	logger.Info(fmt.Sprintf("%s - All wallets fetched successfully", fnName), zap.Any("wallets", wallets))

	resp := &response.WalletResponse{
		Status:  http.StatusOK,
		Wallets: wallets,
	}
	if len(wallets) == 0 {
		resp.Message = utils.Ptr("No wallets found")
	}
	logger.Info(fmt.Sprintf("%s - Sending wallets response", fnName), zap.Any("wallets", wallets))
	SendJSONResponse(fnName, w, resp.Status, resp)
}
