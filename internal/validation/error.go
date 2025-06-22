package validation

import (
	"fmt"
	"strings"
	"time"

	"github.com/ezjuanify/wallet/internal/logger"
	"go.uber.org/zap"
)

type WalletError struct {
	Name      string
	Status    int
	Code      WalletErrorCode
	Message   string
	Timestamp time.Time
	Err       error
	Context   []zap.Field
}

type WalletErrorCode string

const (
	ERR_TRANSACTION_START_FAILED         WalletErrorCode = "ERR_TRANSACTION_START_FAILED"
	ERR_TRANSACTION_COMMIT_FAILED        WalletErrorCode = "ERR_TRANSACTION_COMMIT_FAILED"
	ERR_INVALID_JSON_BODY                WalletErrorCode = "ERR_INVALID_JSON_BODY"
	ERR_DEPOSIT_FAILED                   WalletErrorCode = "ERR_DEPOSIT_FAILED"
	ERR_WITHDRAW_FAILED                  WalletErrorCode = "ERR_WITHDRAW_FAILED"
	ERR_TRANSFER_OUT_FAILED              WalletErrorCode = "ERR_TRANSFER_OUT_FAILED"
	ERR_TRANSFER_IN_FAILED               WalletErrorCode = "ERR_TRANSFER_IN_FAILED"
	ERR_FETCH_TRANSACTION_FAILED         WalletErrorCode = "ERR_FETCH_TRANSACTION_FAILED"
	ERR_LOG_TRANSACTION_FAILED           WalletErrorCode = "ERR_LOG_TRANSACTION_FAILED"
	ERR_SANITIZE_USERNAME_FAILED         WalletErrorCode = "ERR_SANITIZE_USERNAME_FAILED"
	ERR_AMOUNT_VALIDATION_FAILED         WalletErrorCode = "ERR_AMOUNT_VALIDATION_FAILED"
	ERR_WALLET_BALANCE_VALIDATION_FAILED WalletErrorCode = "ERR_WALLET_BALANCE_VALIDATION_FAILED"
	ERR_INSUFFICIENT_WALLET_BALANCE      WalletErrorCode = "ERR_INSUFFICIENT_WALLET_BALANCE"
	ERR_FETCH_WALLET_FAILED              WalletErrorCode = "ERR_FETCH_WALLET_FAILED"
	ERR_DB_UPSERT_FAILED                 WalletErrorCode = "ERR_DB_UPSERT_FAILED"
	ERR_DB_WITHDRAW_FAILED               WalletErrorCode = "ERR_DB_WITHDRAW_FAILED"
	ERR_WALLET_DOES_NOT_EXIST            WalletErrorCode = "ERR_WALLET_DOES_NOT_EXIST"
	ERR_ZERO_AMOUNT                      WalletErrorCode = "ERR_ZERO_AMOUNT"
	ERR_PANIC_OCCURED                    WalletErrorCode = "ERR_PANIC_OCCURED"
)

type AppErrors struct {
	errs []WalletError
}

func NewHandlerErrors() *AppErrors {
	return &AppErrors{}
}

func (ae *AppErrors) AddError(e WalletError) {
	logger.Warn(fmt.Sprintf("%s - Adding error", e.Name), zap.Any("error", e))
	ae.errs = append(ae.errs, e)
}

func (ae AppErrors) First() *WalletError {
	if len(ae.errs) > 0 {
		return &ae.errs[0]
	}
	return nil
}

func (ae AppErrors) GetErrsCount() int {
	return len(ae.errs)
}

func (ae AppErrors) LogAll() {
	for _, e := range ae.errs {
		zapFields := []zap.Field{
			zap.String("name", e.Name),
			zap.String("code", string(e.Code)),
			zap.String("timestamp_format", e.Timestamp.Format(time.RFC3339)),
		}

		if len(e.Context) > 0 {
			zapFields = append(zapFields, e.Context...)
		}

		if e.Err != nil {
			zapFields = append(zapFields, zap.Error(e.Err))
		}

		logger.Error(e.Message, zapFields...)
	}
}

func (ae AppErrors) LogSummary() {
	logger.Error("Application errors summary", zap.Int("error_count", len(ae.errs)), zap.String("errors", ae.Error()))
}

func (ae AppErrors) Error() string {
	var b strings.Builder
	for _, e := range ae.errs {
		fmt.Fprintf(&b, "%s - %s - %s - %s: %v\n",
			e.Timestamp.Format(time.RFC3339),
			e.Name,
			e.Code,
			e.Message,
			e.Err,
		)
	}
	return b.String()
}
