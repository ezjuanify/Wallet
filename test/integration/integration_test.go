package integration

import (
	"testing"

	"github.com/ezjuanify/wallet/test/integration/testcase"
)

func TestIntegration(t *testing.T) {
	tests := []testcase.TestCase{}
	tests = append(tests, testcase.AddDepositTestCases()...)
	tests = append(tests, testcase.AddWithdrawTestCases()...)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if err := dbTestHarness.DoTestResetDBState(); err != nil {
				t.Errorf("%s", err)
			}

			if test.InitialWallet != nil {
				if err := dbTestHarness.DoTestInsertInitialWallet(test.InitialWallet); err != nil {
					t.Errorf("%s", err)
				}
			}

			resp, err := DoTestRequest(test.TxnType, test.Payload, TEST_WALLET_HOST, TEST_WALLET_PORT)
			if err != nil {
				t.Errorf("%s", err)
			}
			defer resp.Body.Close()

			if err := DoTestStatusValidation(test.ExpectErr, resp); err != nil {
				t.Errorf("%s", err)
			}

			respWallet, err := DoTestRespFormatValidation(resp)
			if !test.ExpectErr && err != nil {
				t.Errorf("%s", err)
			}

			if err := DoTestBuildExpectedWallet("Test Set Expected Wallet", &test); !test.ExpectErr && err != nil {
				t.Errorf("%s", err)
			}

			if err := DoTestWalletValidation("Test API Response Wallet", test.TxnType, test.ExpectedWallet, respWallet); !test.ExpectErr && err != nil {
				t.Errorf("%s", err)
			}

			dbWallet, err := dbTestHarness.DoTestFetchWalletFromDB(test.SanitizedPayloadUsername())
			if !test.ExpectErr && err != nil {
				t.Errorf("%s", err)
			}

			if err := DoTestWalletValidation("Test DB Source Wallet", test.TxnType, test.ExpectedWallet, dbWallet); !test.ExpectErr && err != nil {
				t.Errorf("%s", err)
			}

			dbTransaction, err := dbTestHarness.DoTestFetchTransaction(test.SanitizedPayloadUsername())
			if !test.ExpectErr && err != nil {
				t.Errorf("%s", err)
			}

			if err := DoTestTransactionValidation("Test DB Source Transaction", test, dbTransaction); !test.ExpectErr && err != nil {
				t.Errorf("%s", err)
			}
		})
	}
}
