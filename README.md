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

Deposit funds into a wallet.

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

### POST `/withdraw`

Withdraw funds from wallet.

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
