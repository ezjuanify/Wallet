# Wallet API

A RESTful wallet service written in **Go**, supporting **deposit**, **withdrawal**, **transfer**, and **transaction history**, with strict input validation and full integration test coverage.

## Tech Stack

- **Go**
- **PostgreSQL**
- **Docker / Docker Compose**
- `database/sql` with `pgx` driver
- Unit + Integration test suite

## Getting Started

### Prerequisites

- Go `1.xx`
- Docker

### Run with Docker Compose

```bash
git clone https://github.com/yourname/wallet.git
cd wallet
docker-compose up --build
```

App will be available at: `http://localhost:8080`


## API Endpoints

### POST `/deposit`

Deposit funds into a user wallet.

#### Request
```json
{
  "username": "juan",
  "amount": 500
}
```

#### Response
```json
{
  "status": 200,
  "action": "deposit",
  "wallet": {
    "username": "JUAN",
    "balance": 500,
    "lastDepositAmount": 500,
    "lastDepositUpdated": "2025-06-17T09:28:00.376856Z",
    "lastWithdrawAmount": null,
    "lastWithdrawUpdated": null
  }
}
```

---

### POST `/withdraw`

Withdraw funds from user wallet.

#### Request
```json
{
    "username": "juan",
    "amount": 1000
}
```

#### Response
```json
{
    "status": 200,
    "action": "withdraw",
    "wallet": {
        "username": "JUAN",
        "balance": 0,
        "lastDepositAmount": 500,
        "lastDepositUpdated": "2025-06-17T19:13:02.722774Z",
        "lastWithdrawAmount": 1000,
        "lastWithdrawUpdated": "2025-06-17T19:13:04.857005Z"
    }
}
```

---

### POST `/transfer`

Transfer funds from user wallet to counterparty wallet.

#### Request
```json
{
    "username": "juan",
    "amount": 500,
    "counterparty": "mary"
}
```

#### Response
```json
{
    "status": 200,
    "action": "transfer",
    "wallet": {
        "username": "JUAN",
        "balance": 4000,
        "lastDepositAmount": 5000,
        "lastDepositUpdated": "2025-06-19T19:10:11.082386Z",
        "lastWithdrawAmount": 500,
        "lastWithdrawUpdated": "2025-06-19T19:10:14.430453Z"
    },
    "counterparty": "MARY"
}
```

---

### GET `/transactions`

Get transactions based on url parameters. Accepts the following params:

- **username** - Search by username
- **counterparty** - Search by counterparty
- **type** - Search by transaction type (deposit, withdraw, transfer_in, transfer_out)
- **limit** - Number of results to return

#### URL Params
```
localhost:8080/transactions?username=juan
```

#### Response
```json
{
    "status": 200,
    "criteria": {
        "username": "JUAN"
    },
    "transactions": [
        {
            "ID": 6,
            "username": "JUAN",
            "txnType": "transfer_out",
            "amount": 200,
            "counterparty": "MARY",
            "timestamp": "2025-06-20T18:44:24.477541Z",
            "hash": "a7daa5cbb02736bef787cf26b4f8f05a9fc841b36fc77f8b7c3a37a499c19710"
        },
        {
            "ID": 4,
            "username": "JUAN",
            "txnType": "transfer_out",
            "amount": 100,
            "counterparty": "MARY",
            "timestamp": "2025-06-20T18:44:20.031824Z",
            "hash": "8fb9e0158a425d1f710e713adcf0b55f74939126d4db697652215ea3851738c8"
        },
        {
            "ID": 3,
            "username": "JUAN",
            "txnType": "withdraw",
            "amount": 500,
            "counterparty": null,
            "timestamp": "2025-06-20T18:44:18.298866Z",
            "hash": "4c99053a7a8b566f4a1e4eb71efe29f0f977de18b61643f461bc77554556478d"
        },
        {
            "ID": 1,
            "username": "JUAN",
            "txnType": "deposit",
            "amount": 2000,
            "counterparty": null,
            "timestamp": "2025-06-20T18:44:08.593154Z",
            "hash": "394ee8225f8f9a35e2f8b79df17f32533490497a7c98dd0ed26cc42eb8459155"
        }
    ]
}
```

---

### GET `/balance`

Get user wallet. Accepts the following params:

- **username** - Search by username

#### URL Params
```
localhost:8080/balance?username=juan
```

#### Response
```json
{
    "status": 200,
    "wallet": {
        "username": "JUAN",
        "balance": 1300,
        "lastDepositAmount": 2000,
        "lastDepositUpdated": "2025-06-22T12:51:22.490346Z",
        "lastWithdrawAmount": 200,
        "lastWithdrawUpdated": "2025-06-22T12:52:22.519242Z"
    }
}
```

## Testing

### Unit Tests

```bash
go test -v ./internal/...
```

### Integration Tests

```bash
go test -v ./...
```

## License

MIT
