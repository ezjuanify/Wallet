package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/request"
	"github.com/ezjuanify/wallet/internal/model/response"
	"github.com/ezjuanify/wallet/internal/utils"
	"github.com/ezjuanify/wallet/test/integration/testcase"
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

func DoTestBuildExpectedWallet(test_name string, test *testcase.TestCase) error {
	if test.InitialWallet != nil {
		*test.ExpectedWallet = *test.InitialWallet
	} else {
		test.ExpectedWallet.Username = test.SanitizedPayloadUsername()
		test.ExpectedWallet.Balance = 0
	}

	switch test.TxnType {
	case model.TypeDeposit:
		test.ExpectedWallet.Balance += test.Payload.Amount
		test.ExpectedWallet.LastDepositAmount = utils.PtrInt64(test.Payload.Amount)
	case model.TypeWithdraw:
		test.ExpectedWallet.Balance -= test.Payload.Amount
		test.ExpectedWallet.LastWithdrawAmount = utils.PtrInt64(test.Payload.Amount)
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

func DoTestTransactionValidation(test_name string, test testcase.TestCase, transaction *model.Transaction) error {
	if test.ExpectedWallet == nil {
		return fmt.Errorf("%s: wallet is nil", test_name)
	}
	if transaction == nil {
		return fmt.Errorf("%s: transaction is nil", test_name)
	}

	if test.TxnType != model.TransactionType(transaction.Type) {
		return fmt.Errorf("%s: Transaction type - expected %s but got %s", test_name, test.TxnType, transaction.Type)
	}

	if test.ExpectedWallet.Username != transaction.Username {
		return fmt.Errorf("%s: Username - expected %s but got %s", test_name, test.ExpectedWallet.Username, transaction.Username)
	}

	switch test.TxnType {
	case model.TypeDeposit:
		if test.ExpectedWallet.LastDepositAmount == nil {
			return fmt.Errorf("%s: LastDepositAmount - expected a value but got nil", test_name)
		}
		if *test.ExpectedWallet.LastDepositAmount != transaction.Amount {
			return fmt.Errorf("%s: LastDepositAmount - expected %d but got %d", test_name, *test.ExpectedWallet.LastDepositAmount, transaction.Amount)
		}
	case model.TypeWithdraw:
		if test.ExpectedWallet.LastWithdrawAmount == nil {
			return fmt.Errorf("%s: LastWithdrawAmount - expected a value but got nil", test_name)
		}
		if *test.ExpectedWallet.LastWithdrawAmount != transaction.Amount {
			return fmt.Errorf("%s: LastWithdrawAmount - expected %d but got %d", test_name, *test.ExpectedWallet.LastWithdrawAmount, transaction.Amount)
		}
	}

	if test.Payload.Amount != transaction.Amount {
		return fmt.Errorf("%s: Amount - expected %d but got %d", test_name, test.Payload.Amount, transaction.Amount)
	}

	if payloadHash := utils.GenerateTransactionHash(strings.ToUpper(test.SanitizedPayloadUsername()), string(test.TxnType), test.Payload.Amount, nil, transaction.Timestamp.UTC().Format(time.RFC3339)); payloadHash != transaction.Hash {
		return fmt.Errorf("%s: Calculated payload hash %s does not match %s", test_name, payloadHash[:10], transaction.Hash[:10])
	}
	return nil
}
