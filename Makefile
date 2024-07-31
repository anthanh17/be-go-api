# ----------------------------------------- Setup create database ------------------------------------------------------
# docker run -itd --name my-redis -p 6379:6379 redis:latest
# Setup postgres database docker
createdbcontainer:
	docker run --name monday-auth-database -p 5432:5432 -v /data/postgres:/var/lib/postgresql/data -e POSTGRES_USER=root -e POSTGRES_PASSWORD=abc123 -d postgres:12-alpine

createdb:
	docker exec -it monday-auth-database createdb --username=root --owner=root monday_auth

dropdb:
	docker exec -it monday-auth-database dropdb monday_auth
# --------------------------------------------------------------------------------------------------------------------------
# -------------------------------------- Read file schema sql crete or update database --------------------------------------
# Migarte database all
migrateup:
	migrate -path db/migration -database "postgresql://root:abc123@ec2-18-141-12-199.ap-southeast-1.compute.amazonaws.com:5432/monday_auth?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:abc123@ec2-18-141-12-199.ap-southeast-1.compute.amazonaws.com:5432/monday_auth?sslmode=disable" -verbose down

# Migarte database lastest
migrateup1:
	migrate -path db/migration -database "postgresql://root:abc123@ec2-18-141-12-199.ap-southeast-1.compute.amazonaws.com:5432/monday_auth?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:abc123@ec2-18-141-12-199.ap-southeast-1.compute.amazonaws.com:5432/monday_auth?sslmode=disable" -verbose down 1
# --------------------------------------------------------------------------------------------------------------------------
# ---------------------------------- Define schema dabase and define query sqlc generate code golang -----------------------
# create file config sqlc.yaml
sqlcinit:
	sqlc init

# sqlc gen code golang
sqlc:
	sqlc generate -f ./configs/sqlc.yaml
# --------------------------------------------------------------------------------------------------------------------------
# Unit test
test:
	go test -v -cover ./...

# Start server http
server:
	go run main.go

.PHONY: createdbcontainer createdb dropdb migrateup migratedown sqlcinit sqlc test server migratedown1 migrateup1
