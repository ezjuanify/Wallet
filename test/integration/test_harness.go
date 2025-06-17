package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ezjuanify/wallet/internal/handler/request"
	"github.com/ezjuanify/wallet/internal/handler/response"
	"github.com/ezjuanify/wallet/internal/model"
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

func DoTestSetExpectedWallet(test_name string, txnType model.TransactionType, expectedWallet *model.Wallet, initialWallet *model.Wallet, payload *request.RequestPayload) error {
	if initialWallet != nil {
		*expectedWallet = *initialWallet
	} else {
		expectedWallet.Username = strings.ToUpper(payload.Username)
		expectedWallet.Balance = 0
	}

	switch txnType {
	case model.TypeDeposit:
		expectedWallet.Balance += payload.Amount
		expectedWallet.LastDepositAmount = utils.PtrInt64(payload.Amount)
	case model.TypeWithdraw:
		expectedWallet.Balance -= payload.Amount
		expectedWallet.LastWithdrawAmount = utils.PtrInt64(payload.Amount)
	}

	return nil
}

func DoTestWalletValidation(test_name string, txnType model.TransactionType, expectedWallet *model.Wallet, wallet *model.Wallet) error {
	if wallet == nil {
		return fmt.Errorf("%s: Wallet is nil", test_name)
	}

	if expectedWallet.Username != wallet.Username {
		return fmt.Errorf("%s: Username - expected %s but got %s", test_name, expectedWallet.Username, wallet.Username)
	}

	if expectedWallet.Balance != wallet.Balance {
		return fmt.Errorf("%s: Balance - expected %d but got %d", test_name, expectedWallet.Balance, wallet.Balance)
	}

	// Last Deposit Amount
	if txnType == model.TypeDeposit && wallet.LastDepositAmount == nil {
		return fmt.Errorf("%s: LastDepositAmount - expected %d but got nil", test_name, *expectedWallet.LastDepositAmount)
	} else if txnType == model.TypeDeposit && wallet.LastDepositAmount != nil && *expectedWallet.LastDepositAmount != *wallet.LastDepositAmount {
		return fmt.Errorf("%s: LastDepositAmount - expected %d but got %d", test_name, *expectedWallet.LastDepositAmount, *wallet.LastDepositAmount)
	}

	// Last Deposit Updated
	if txnType == model.TypeDeposit && wallet.LastDepositUpdated == nil {
		return fmt.Errorf("%s: LastDepositUpdated - expected timestamp but got nil", test_name)
	}

	// Last Withdraw Amount
	if txnType == model.TypeWithdraw && wallet.LastWithdrawAmount == nil {
		return fmt.Errorf("%s: LastWithdrawAmount - expected %d but got nil", test_name, *expectedWallet.LastWithdrawAmount)
	} else if txnType == model.TypeWithdraw && wallet.LastWithdrawAmount != nil && *expectedWallet.LastWithdrawAmount != *wallet.LastWithdrawAmount {
		return fmt.Errorf("%s: LastWithdrawAmount - expected %d but got %d", test_name, *expectedWallet.LastWithdrawAmount, *wallet.LastWithdrawAmount)
	}

	// Last Withdraw Updated
	if txnType == model.TypeWithdraw && wallet.LastWithdrawUpdated == nil {
		return fmt.Errorf("%s: LastWithdrawUpdated - expected timestamp but got nil", test_name)
	}
	return nil
}

func DoTestTransactionValidation(test_name string, txnType model.TransactionType, payload *request.RequestPayload, expectedWallet *model.Wallet, transaction *model.Transaction) error {
	if expectedWallet == nil {
		return fmt.Errorf("%s: wallet is nil", test_name)
	}
	if transaction == nil {
		return fmt.Errorf("%s: transaction is nil", test_name)
	}

	if txnType != model.TransactionType(transaction.Type) {
		return fmt.Errorf("%s: Transaction type - expected %s but got %s", test_name, txnType, transaction.Type)
	}

	if expectedWallet.Username != transaction.Username {
		return fmt.Errorf("%s: Username - expected %s but got %s", test_name, expectedWallet.Username, transaction.Username)
	}

	switch txnType {
	case model.TypeDeposit:
		if expectedWallet.LastDepositAmount == nil {
			return fmt.Errorf("%s: LastDepositAmount - expected a value but got nil", test_name)
		}
		if *expectedWallet.LastDepositAmount != transaction.Amount {
			return fmt.Errorf("%s: LastDepositAmount - expected %d but got %d", test_name, *expectedWallet.LastDepositAmount, transaction.Amount)
		}
	case model.TypeWithdraw:
		if expectedWallet.LastWithdrawAmount == nil {
			return fmt.Errorf("%s: LastWithdrawAmount - expected a value but got nil", test_name)
		}
		if *expectedWallet.LastWithdrawAmount != transaction.Amount {
			return fmt.Errorf("%s: LastWithdrawAmount - expected %d but got %d", test_name, *expectedWallet.LastWithdrawAmount, transaction.Amount)
		}
	}

	if payload.Amount != transaction.Amount {
		return fmt.Errorf("%s: Amount - expected %d but got %d", test_name, payload.Amount, transaction.Amount)
	}

	if payloadHash := utils.GenerateTransactionHash(strings.ToUpper(payload.Username), string(txnType), payload.Amount, nil, transaction.Timestamp.UTC().Format(time.RFC3339)); payloadHash != transaction.Hash {
		return fmt.Errorf("%s: Calculated payload hash %s does not match %s", test_name, payloadHash[:10], transaction.Hash[:10])
	}
	return nil
}
