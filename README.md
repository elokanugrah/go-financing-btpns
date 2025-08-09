# go-financing-btpns
project test btpns

## Tech Stack

  * **Language:** Go
  * **Framework:** Gin Gonic
  * **Database:** PostgreSQL
  * **Containerization:** Docker
  * **Testing:** Testify, Mockery

## Getting Started

### Prerequisites

  * Go (v1.23+)
  * Docker & Docker Compose
  * `migrate-cli`
  * `mockery`

### Tooling Installation

  * `migrate-cli` is a command-line tool for managing database migrations.

    ```bash
    go install -tags 'postgres' [github.com/golang-migrate/migrate/v4/cmd/migrate@latest](https://github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
    ```
  * `mockery` is a tool for automatically generating mocks for your Go interfaces.
  
    ```bash
    go install [github.com/vektra/mockery/v2@latest](https://github.com/vektra/mockery/v2@latest)
    ```

### Running with Docker

1.  **Clone the repository**

    ```bash
    git clone https://github.com/elokanugrah/go-financing-btpns.git
    cd go-order-system
    ```

2.  **Create `.env` file**
    Create a `.env` file in the root directory.

    ```ini
    SERVER_PORT=9000
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=user
    DB_PASSWORD=password
    DB_NAME=financing_db
    ```

3.  **Run Services**

    ```bash
    docker-compose up --build
    ```

4.  **Run Migrations & Seeder** (in a new terminal)

    ```bash
    # Create database schema
    migrate -database "postgres://user:password@localhost:5432/financing_db?sslmode=disable" -path migration up

    # Seed product data
    go run ./cmd/seed
    ```

The API is now running at `http://localhost:9000`.

## API Endpoints


| Method | Endpoint              | Description              |
| :----- | :-------------------- | :----------------------- |
| `POST` | `/calculate-installments`      | Get calculation.    |
| `POST`  | `/submit-financing`      | Submit financing.       |


**Example: Get Installment Calculations**

```bash
curl --location '127.0.0.1:9000/calculate-installments' \
--header 'Content-Type: application/json' \
--data '{
  "amount": 10000000
}'
```

**Example: Get Submit Financing**

```bash
curl --location '127.0.0.1:9000/submit-financing' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": 1,
    "facility_limit_id": 1,
    "amount": 5000000,
    "tenor": 12,
    "start_date": "2025-03-01"
}'
```

## Running Tests

To run all unit and integration tests, ensure the database is running and execute:

```bash
go test -v ./...
```