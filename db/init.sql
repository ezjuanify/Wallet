CREATE TABLE IF NOT EXISTS wallets (
    id                    SERIAL  PRIMARY KEY,
    username              TEXT    UNIQUE         NOT NULL,
    balance               BIGINT                 NOT NULL DEFAULT 0,
    last_deposit_amount   BIGINT,
    last_deposit_updated  TIMESTAMP,
    last_withdraw_amount  BIGINT,
    last_withdraw_updated TIMESTAMP
);
ALTER TABLE wallets ADD CONSTRAINT chk_wallet_balance CHECK (balance >= 0 AND balance <= 999999);

CREATE TABLE IF NOT EXISTS transactions (
    id           SERIAL  PRIMARY KEY,
    username     TEXT                  NOT NULL,
    type         TEXT                  NOT NULL CHECK (type IN ('deposit', 'withdraw', 'transfer_in', 'transfer_out')),
    amount       BIGINT                NOT NULL CHECK (amount > 0),
    counterparty TEXT,
    timestamp    TIMESTAMP             NOT NULL DEFAULT now(),
    hash         TEXT                  NOT NULL
);
