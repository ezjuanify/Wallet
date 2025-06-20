package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/request"
	"github.com/ezjuanify/wallet/internal/model/response"
	"github.com/ezjuanify/wallet/internal/utils"
)

func waitForHTTPServerReady(maxRetries int, interval time.Duration, host string, port string) error {
	url := fmt.Sprintf("http://%s%s/health", host, port)
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Println("✅ API is up.")
			return nil
		}
		time.Sleep(interval)
	}
	return fmt.Errorf("HTTP server not ready after %d attempts", maxRetries)
}

func DoTestRequest(txnType model.TransactionType, payload *request.RequestPayload, host string, port string) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s%s/%s", host, port, txnType), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	return resp, nil
}

func DoTestStatusValidation(expectErr bool, resp *http.Response) error {
	if expectErr && resp.StatusCode == http.StatusOK {
		return fmt.Errorf("expected error, got success")
	}

	if !expectErr && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	return nil
}

func DoTestRespFormatValidation(resp *http.Response) (*model.Wallet, error) {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %s", err)
	}

	var result response.TransactionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return &result.Wallet, nil
}

func DoTestWalletValidation(test_name string, txnType model.TransactionType, expected *model.Wallet, actual *model.Wallet, isCounterparty bool) error {
	if actual == nil {
		return fmt.Errorf("%s: Wallet is nil", test_name)
	}

	if expected.Username != actual.Username {
		return fmt.Errorf("%s: Username - expected %s but got %s", test_name, expected.Username, actual.Username)
	}

	if expected.Balance != actual.Balance {
		return fmt.Errorf("%s: Balance - expected %d but got %d", test_name, expected.Balance, actual.Balance)
	}

	switch txnType {
	case model.TypeDeposit:
		if actual.LastDepositAmount == nil {
			return fmt.Errorf("%s: LastDepositAmount - expected %d but got nil", test_name, *expected.LastDepositAmount)
		}
		if *actual.LastDepositAmount != *expected.LastDepositAmount {
			return fmt.Errorf("%s: LastDepositAmount - expected %d but got %d", test_name, *expected.LastDepositAmount, *actual.LastDepositAmount)
		}
		if actual.LastDepositUpdated == nil {
			return fmt.Errorf("%s: LastDepositUpdated - expected timestamp but got nil", test_name)
		}
	case model.TypeWithdraw:
		if actual.LastWithdrawAmount == nil {
			return fmt.Errorf("%s: LastWithdrawAmount - expected %d but got nil", test_name, *expected.LastWithdrawAmount)
		}
		if *actual.LastWithdrawAmount != *expected.LastWithdrawAmount {
			return fmt.Errorf("%s: LastWithdrawAmount - expected %d but got %d", test_name, *expected.LastWithdrawAmount, *actual.LastWithdrawAmount)
		}
		if actual.LastWithdrawUpdated == nil {
			return fmt.Errorf("%s: LastWithdrawUpdated - expected timestamp but got nil", test_name)
		}
	case model.TypeTransfer:
		if isCounterparty {
			if actual.LastDepositAmount == nil || actual.LastDepositUpdated == nil {
				return fmt.Errorf("%s: Transfer-In - expected deposit info but got nil", test_name)
			}
			if *actual.LastDepositAmount != *expected.LastDepositAmount {
				return fmt.Errorf("%s: Transfer-In - LastDepositAmount - expected %d but got %d", test_name, *expected.LastDepositAmount, *actual.LastDepositAmount)
			}
		} else {
			if actual.LastWithdrawAmount == nil || actual.LastWithdrawUpdated == nil {
				return fmt.Errorf("%s: Transfer-In - expected withdraw info but got nil", test_name)
			}
			if *actual.LastWithdrawAmount != *expected.LastWithdrawAmount {
				return fmt.Errorf("%s: Transfer-In - LastWithdrawAmount - expected %d but got %d", test_name, *expected.LastWithdrawAmount, *actual.LastWithdrawAmount)
			}
		}
	}
	return nil
}

func DoTestTransactionValidation(test_name string, txnType model.TransactionType, expected *model.Wallet, counterparty *model.Wallet, transaction *model.Transaction, isCounterparty bool) error {
	if transaction == nil {
		return fmt.Errorf("%s: transaction is nil", test_name)
	}

	if expected == nil {
		return fmt.Errorf("%s: wallet is nil", test_name)
	}

	if expected.Username != transaction.Username {
		return fmt.Errorf("%s: Username - expected %s but got %s", test_name, expected.Username, transaction.Username)
	}

	var amount int64
	var counterpartyUsername *string
	if counterparty != nil {
		counterpartyUsername = &counterparty.Username
	}

	switch txnType {
	case model.TypeDeposit:
		if model.TransactionType(transaction.Type) != model.TypeDeposit {
			return fmt.Errorf("%s: Transaction type - expected %s but got %s", test_name, model.TypeDeposit, transaction.Type)
		}
		if expected.LastDepositAmount == nil {
			return fmt.Errorf("%s: LastDepositAmount - expected a value but got nil", test_name)
		}
		if *expected.LastDepositAmount != transaction.Amount {
			return fmt.Errorf("%s: LastDepositAmount - expected %d but got %d", test_name, *expected.LastDepositAmount, transaction.Amount)
		}
		amount = *expected.LastDepositAmount
	case model.TypeWithdraw:
		if model.TransactionType(transaction.Type) != model.TypeWithdraw {
			return fmt.Errorf("%s: Transaction type - expected %s but got %s", test_name, model.TypeWithdraw, transaction.Type)
		}
		if expected.LastWithdrawAmount == nil {
			return fmt.Errorf("%s: LastWithdrawAmount - expected a value but got nil", test_name)
		}
		if *expected.LastWithdrawAmount != transaction.Amount {
			return fmt.Errorf("%s: LastWithdrawAmount - expected %d but got %d", test_name, *expected.LastWithdrawAmount, transaction.Amount)
		}
		amount = *expected.LastWithdrawAmount
	case model.TypeTransfer:
		if isCounterparty {
			if model.TransactionType(transaction.Type) != model.TypeTransferIn {
				return fmt.Errorf("%s: Transaction type - expected %s but got %s", test_name, model.TypeTransferIn, transaction.Type)
			}
			if expected.LastDepositAmount == nil {
				return fmt.Errorf("%s: LastDepositAmount - expected a value but got nil", test_name)
			}
			if *expected.LastDepositAmount != transaction.Amount {
				return fmt.Errorf("%s: LastDepositAmount - expected %d but got %d", test_name, *expected.LastDepositAmount, transaction.Amount)
			}
			amount = *expected.LastDepositAmount
		} else {
			if model.TransactionType(transaction.Type) != model.TypeTransferOut {
				return fmt.Errorf("%s: Transaction type - expected %s but got %s", test_name, model.TypeTransferOut, transaction.Type)
			}
			if expected.LastWithdrawAmount == nil {
				return fmt.Errorf("%s: LastWithdrawAmount - expected a value but got nil", test_name)
			}
			if *expected.LastWithdrawAmount != transaction.Amount {
				return fmt.Errorf("%s: LastWithdrawAmount - expected %d but got %d", test_name, *expected.LastWithdrawAmount, transaction.Amount)
			}
			amount = *expected.LastWithdrawAmount
		}

		if counterpartyUsername == nil {
			return fmt.Errorf("%s: Counterparty - expected counterparty username but got nil", test_name)
		}
		if transaction.Counterparty == nil {
			return fmt.Errorf("%s: Counterparty - expected %s but got nil", test_name, *counterpartyUsername)
		}
		if *counterpartyUsername != *transaction.Counterparty {
			return fmt.Errorf("%s: Counterparty - expected %s but got %s", test_name, *counterpartyUsername, *transaction.Counterparty)
		}
	}

	if payloadHash := utils.GenerateTransactionHash(
		expected.Username,
		transaction.Type,
		amount,
		counterpartyUsername,
		transaction.Timestamp.UTC().Format(time.RFC3339),
	); payloadHash != transaction.Hash {
		return fmt.Errorf("%s: Calculated payload hash %s does not match %s", test_name, payloadHash[:10], transaction.Hash[:10])
	}
	return nil
}
