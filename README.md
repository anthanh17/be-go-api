# github.com/anthanh17/be-go-api
Viết một web server sử dụng gin go + redis với
1. 1 API /login, để tạo session cho mỗi người đăng nhập, dùng redis để lưu session id, user name ấy
2. 1 API /ping chỉ cho phép 1 người được gọi tại một thời điểm ( với sleep ở bên trong api đó trong 5s)
3. đếm số lượng lần 1 người gọi api /ping
4. rate limit mỗi người chỉ được gọi API /ping 2 lần trong 60s
5. 1 API /top/ trả về top 10 người gọi API /ping nhiều nhất
6. Dùng hyperloglog để lưu xấp sỉ số người gọi api /ping , và trả về trong api /count

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
