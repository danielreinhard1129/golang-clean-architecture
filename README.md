# Golang Clean Architecture

A Go project using **Fiber**, **GORM**, and **Atlas** for database migrations.

---

## ⚙️ Setup & Installation

### 1. Install dependencies

```bash
go mod tidy
```

### 2. Install tools

- [Air](https://github.com/air-verse/air) (live reload):

  ```bash
  go install github.com/air-verse/air@latest
  ```

- [Atlas](https://atlasgo.io/cli/getting-started) (database migration):

  ```bash
  curl -sSf https://atlasgo.sh | sh
  ```

### 3. Environment variables

Create a `.env` file:

```env
POSTGRES_HOST=localhost
POSTGRES_PORT=6543
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=appdb
APP_PORT=8000
JWT_SECRET=secret

DATABASE_URL=postgres://postgres:postgres@localhost:6543/appdb?sslmode=disable
DATABASE_DEV_URL=postgres://postgres:postgres@localhost:6543/appdb_dev?sslmode=disable

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=email@mail.com
SMTP_PASS=app_password
SMTP_FROM=email@mail.com

```

### 4. Start database (Recommended: Docker)

It is recommended to use Docker Compose to run the database because the dev database should be temporary, empty, and locally-run. **Atlas** uses it to parse, validate, and analyze SQL definitions, and it is cleaned up after the migration process.

```bash
make docker-up
```

### 5. Run the app

```bash
make watch     # with Air (live reload)
make run       # build & run binary
```

---

## 🗄️ Database Migration

- Create new migration:

  ```bash
  make migrate-diff name=add_table_users
  ```

- Apply migrations:

  ```bash
  make migrate-apply
  ```

- Show migration status:

  ```bash
  make migrate-status
  ```

- Rollback last migration:

  ```bash
  make migrate-down
  ```
