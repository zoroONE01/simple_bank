# GitHub Copilot Instructions for This Golang "Simple Bank" Project

## 1. Project Overview

This project is a backend API service for a "Simple Bank" application.
Key functionalities include:

- User account creation and management.
- Displaying account balances and details.
- Recording and listing financial transactions.
- Transferring funds between accounts.

The API is built using the Gin web framework (`github.com/gin-gonic/gin`).
**PostgreSQL** is the chosen database, managed via **sqlc** for generating type-safe Go code from SQL queries.
The development environment and services (like Postgres) are containerized using **Docker**.
Testing will heavily leverage **mockery** for generating mocks of interfaces.

The focus is on clear API endpoints, data integrity, robust error handling, and efficient database interactions through `sqlc`-generated code.

## 2. Coding Conventions and Style (Go Specific)

### 2.1. Formatting

- **Adhere to `gofmt`**: All Go code MUST be formatted with `gofmt` (or `goimports`).

### 2.2. Naming Conventions

- **Package Names**: Short, concise, all-lowercase (e.g., `account`, `transaction`, `dbstore`, `api`).
  - The package for `sqlc`-generated code is often named `db` or `dbstore`.
- **Variable Names**: `camelCase` for local, `PascalCase` for exported. Acronyms like `ID`, `URL`, `API` in `PascalCase`.
- **Function Names (API Handlers with Gin)**: `CreateAccountHandler`, `GetAccountHandler`.
- **Interface Names**: Often end with "er" (e.g., `AccountReader`, `TransactionLogger`). `sqlc` generates a `Querier` interface.
- **Struct Names**: `PascalCase`. API request/response structs: `CreateAccountRequest`, `AccountResponse`. `sqlc` generates structs from your SQL schema and queries.

### 2.3. Comments

- **Exported Identifiers**: MUST have a doc comment.
- **SQL Queries (`.sql` files for `sqlc`)**: Comment complex queries or non-obvious logic directly in the SQL files. `sqlc` preserves these comments in the generated Go code.
  - Example SQL comment for `sqlc`:

      ```sql
      -- name: GetAccount :one
      -- Retrieves a specific account by its ID.
      SELECT * FROM accounts
      WHERE id = $1 LIMIT 1;
      ```

- **Business Logic**: Comment complex banking logic.

### 2.4. Imports

- **Group Imports**: Use `goimports`.

## 3. Standard Library Preference

- **Prioritize Standard Library**: Use Go's standard library where appropriate.
- **Context Package**: Always use `context.Context`, passed through Gin handlers down to `sqlc` query methods.
  - `sqlc`-generated methods will typically expect `context.Context` as their first argument.

## 4. Error Handling (API Context)

- **Explicit Error Handling**: Always check errors. `sqlc` methods return an `error`.
- **Error Wrapping**: Use `fmt.Errorf` with `%w`.
- **Database Errors**: `sqlc` often returns standard `database/sql` errors like `sql.ErrNoRows`. Handle these appropriately (e.g., map `sql.ErrNoRows` to HTTP 404).
  - Example with `sqlc` and Gin:

      ```go
      // account, err := server.store.GetAccount(ctx, req.ID) // server.store is the sqlc Querier
      // if err != nil {
      //     if errors.Is(err, sql.ErrNoRows) {
      //         ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
      //         return
      //     }
      //     ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get account"})
      //     return
      // }
      ```

- **Custom Error Types/Codes**: For business logic errors (e.g., insufficient funds), define custom errors.

## 5. Concurrency

- **Goroutines**: Use for background tasks if any. Be mindful with database connections if sharing the `sql.DB` pool.
- **Race Conditions**: Critical. `sqlc` operations are typically atomic at the query level. Ensure business logic spanning multiple queries uses database transactions.

## 6. Testing

- **Table-Driven Tests**: Highly recommended.
- **API Endpoint Testing**: Use `net/http/httptest` to test Gin handlers.
- **Mocking Dependencies**:
  - **`sqlc` `Querier` Interface**: `sqlc` generates a `Querier` interface (often embedded in a `Store` struct). This interface is the primary target for mocking database interactions.
  - **`mockery`**: Use `mockery` (`github.com/vektra/mockery`) to generate mocks for the `Querier` interface and any other service interfaces.
    - Example `mockery` command hint: `//go:generate mockery --name Querier --output ./mocks --outpkg mocks` (place this in the directory with your `Querier` interface).
  - When writing tests, use these mocks to simulate database responses and errors.
  - Example test setup:

      ```go
      // func TestCreateAccountAPI(t *testing.T) {
      //  mockStore := new(mocks.MockQuerier) // Or your specific mock store name
      //  // Setup mock expectations, e.g., mockStore.On("CreateAccount", ...).Return(...)
      //
      //  server := NewServer(mockStore) // Assuming server takes the Querier/Store
      //  router := gin.Default()
      //  // server.SetupRoutes(router)
      //
      //  // ... rest of httptest setup ...
      // }
      ```

## 7. Dependencies and Modules

- **Go Modules**: Project uses Go Modules.
- **Key Dependencies**: Gin, `sqlc`, `mockery`, `lib/pq` (or `jackc/pgx`).

## 8. Patterns to Prefer

- **Clear API Design**: RESTful endpoints.
- **Request Validation**: Gin binding/validation.
- **Service Layer**: Encapsulate business logic. `sqlc` forms the core of the data access layer.
- **Database Transactions with `sqlc`**: For operations requiring atomicity (e.g., transfers), `sqlc` does not manage transactions directly. You need to manage `sql.Tx` objects and pass them to `sqlc`'s `WithTx` method or have `sqlc` generate methods that accept a `*sql.Tx`.
  - Often, a custom `Store` struct will wrap the `sqlc.Queries` and `*sql.DB` to provide transaction helper methods.
  - Example transaction helper in a `Store`:

      ```go
      // func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
      //  tx, err := store.db.BeginTx(ctx, nil)
      //  if err != nil {
      //      return err
      //  }
      //  q := New(tx) // New is the sqlc generated constructor
      //  err = fn(q)
      //  if err != nil {
      //      if rbErr := tx.Rollback(); rbErr != nil {
      //          return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
      //      }
      //      return err
      //  }
      //  return tx.Commit()
      // }
      ```

  - Copilot should be guided to use these transaction helpers when implementing multi-step database operations.

## 9. Patterns to Avoid

- **Fat Gin Handlers**: Delegate to service/store methods.
- **Direct `sql.DB` usage in handlers**: Abstract via a store/repository interface (which `sqlc` helps define).
- **Ignoring `sqlc` generated types**: Use the structs and methods generated by `sqlc` for type safety.

## 10. Logging

- **Structured Logging**: `zerolog`, `zap`, or `slog`.
- **Gin Logging Middleware**: For request details.
- **Log SQL queries (with caution)**: During development, logging executed SQL queries (and their arguments) can be helpful. Be extremely careful about logging sensitive data in production. `sqlc` itself doesn't log queries; this would be part of your Go code or a database driver feature.

## 11. Project-Specific Libraries/Frameworks & Tools

- **Docker & Docker Compose**:
  - Used for setting up and managing the development environment, especially the **PostgreSQL** database.
  - Assume a `docker-compose.yml` file defines services.
  - Copilot might be asked to help with Dockerfile or docker-compose configurations.

- **PostgreSQL**:
  - The primary database.
  - SQL queries for `sqlc` should be written in PostgreSQL syntax.
  - Database schema migrations should be managed (e.g., using `golang-migrate/migrate`, `pressly/goose`, or plain SQL scripts). *(Specify your migration tool if decided)*.

- **`sqlc` (`github.com/sqlc-dev/sqlc`)**:
  - **Configuration**: `sqlc.yaml` defines how `sqlc` discovers SQL queries and generates Go code.
  - **SQL Source Files**: Raw SQL queries are typically stored in `.sql` files (e.g., in a `db/query/` directory).
    - Each query should have a `sqlc` comment defining its name and type (e.g., `:one`, `:many`, `:exec`).
    - Example: `-- name: CreateAccount :one`
  - **Generated Code**: `sqlc` generates Go structs for your schema, input parameters for queries, and methods to execute these queries. It also generates a `Querier` interface.
  - **Usage**: Interact with the database by calling methods on the `sqlc`-generated `*Queries` struct (or the `Querier` interface).
  - Copilot should help write correct SQL for `sqlc` and utilize the generated Go methods.

- **`mockery` (`github.com/vektra/mockery`)**:
  - Used to generate mocks for Go interfaces.
  - Primarily for mocking the `sqlc`-generated `Querier` interface or any custom service interfaces for unit testing.
  - Mocks are typically stored in a `mocks/` directory.
  - Command to generate mocks (often put in a `go:generate` directive): `mockery --name InterfaceName --output ./mocks --outpkg mocks`

- **Gin Web Framework (`github.com/gin-gonic/gin`)**:
  - Routing, request binding/validation, middleware, JSON responses as previously detailed.

- **Database Driver for PostgreSQL**:
  - `github.com/lib/pq` (common) or `github.com/jackc/pgx/v5/stdlib` (modern, often preferred for new projects due to performance and features). *(Specify which one you're using in `sqlc.yaml` and your `go.mod`)*.

## 12. Specific Goals for Copilot

- "Help write a SQL query for `sqlc` in `account.sql` to update an account balance, ensuring it returns the updated account record."
- "Generate the Gin handler for `POST /transfers`, ensuring it uses the `sqlc`-generated `Store`'s transaction helper method for atomicity."
- "Show how to set up a unit test for the `TransferMoney` service method, using a `mockery`-generated mock of the `Querier`."
- "Assist in creating a `docker-compose.yml` file to run the Postgres database and the bank application."
- "What's the `sqlc.yaml` configuration for using `pgx` as the Go database driver?"

---
