# ep-golang-caching

## REQUIREMENT 🙏

1. `sqlc`

> sqlc generates golang code from SQL.

#### Here's how it works:

- You write queries in SQL.
- You run sqlc to generate code with type-safe interfaces to those queries.
- You write application code that calls the generated code.

```
brew install sqlc # macos
sudo snap install sqlc # ubuntu
```

#### Generates golang code from SQL

```
make sqlc
```

2. `golang-migrate`

> Database migrations using CLI

```
brew install golang-migrate
```

3. `make`

## How to use

```
make databaseup
make migrateup
make server
```
