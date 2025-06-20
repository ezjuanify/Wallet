package integration

import (
	"testing"

	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/test/integration/testcase"
)

func TestIntegration(t *testing.T) {
	tests := []testcase.TestCase{}
	tests = append(tests, testcase.AddDepositTestCases()...)
	tests = append(tests, testcase.AddWithdrawTestCases()...)
	tests = append(tests, testcase.AddTransferTestCases()...)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var vErrs ValidationErrors

			if err := dbTestHarness.DoTestResetDBState(); err != nil {
				vErrs.Add("Reset DB State", err)
			}

			if test.InitialWallets != nil {
				for _, v := range test.InitialWallets {
					if err := dbTestHarness.DoTestInsertInitialWallet(&v); err != nil {
						vErrs.Add("Insert Initial Wallet", err)
					}
				}
			}

			resp, err := DoTestRequest(test.TxnType, test.Payload, TEST_WALLET_HOST, TEST_WALLET_PORT)
			if err != nil {
				vErrs.Add("Do Request", err)
			}
			defer resp.Body.Close()

			if err := DoTestStatusValidation(test.ExpectErr, resp); err != nil {
				vErrs.Add("Status Validation", err)
			}

			respWallet, err := DoTestRespFormatValidation(resp)
			if !test.ExpectErr && err != nil {
				vErrs.Add("Response Format", err)
			}

			test.BuildExpectedWallet(test.Payload.Username, false)
			if test.Payload.Counterparty != nil {
				test.BuildExpectedWallet(*test.Payload.Counterparty, true)
			}

			if err := DoTestWalletValidation("API Wallet Validation", test.TxnType, test.ExpectedWallet, respWallet, false); !test.ExpectErr && err != nil {
				vErrs.Add("API Wallet Validation", err)
			}

			dbWallet, err := dbTestHarness.DoTestFetchWalletFromDB(test.ExpectedWallet.Username)
			if !test.ExpectErr && err != nil {
				vErrs.Add("Fetch DB Wallet", err)
			}

			if err := DoTestWalletValidation("DB Wallet Validation", test.TxnType, test.ExpectedWallet, dbWallet, false); !test.ExpectErr && err != nil {
				vErrs.Add("DB Wallet Validation", err)
			}

			dbTransaction, err := dbTestHarness.DoTestFetchTransaction(test.ExpectedWallet.Username)
			if !test.ExpectErr && err != nil {
				vErrs.Add("Fetch Transaction", err)
			}

			if err := DoTestTransactionValidation("Transaction Validation", test.TxnType, test.ExpectedWallet, test.ExpectedCounterpartyWallet, dbTransaction, false); !test.ExpectErr && err != nil {
				vErrs.Add("Transaction Validation", err)
			}

			if test.TxnType != model.TypeTransfer {
				return
			}

			counterpartyWallet, err := dbTestHarness.DoTestFetchWalletFromDB(test.ExpectedCounterpartyWallet.Username)
			if !test.ExpectErr && err != nil {
				vErrs.Add("Fetch Counterparty Wallet", err)
			}

			if err := DoTestWalletValidation("Counterparty Wallet Validation", test.TxnType, test.ExpectedCounterpartyWallet, counterpartyWallet, true); !test.ExpectErr && err != nil {
				vErrs.Add("Counterparty Wallet Validation", err)
			}

			counterpartyTransaction, err := dbTestHarness.DoTestFetchTransaction(test.ExpectedCounterpartyWallet.Username)
			if !test.ExpectErr && err != nil {
				vErrs.Add("Fetch Counterparty Transaction", err)
			}

			if err := DoTestTransactionValidation("Counterparty Transaction Validation", test.TxnType, test.ExpectedCounterpartyWallet, test.ExpectedWallet, counterpartyTransaction, true); !test.ExpectErr && err != nil {
				vErrs.Add("Counterparty Transaction Validation", err)
			}

			vErrs.Report(t)
		})
	}
}
