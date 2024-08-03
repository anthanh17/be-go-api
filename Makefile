HOST = localhost
# HOST = ec2-18-141-12-199.ap-southeast-1.compute.amazonaws.com

# ----------------------------- Setup database ---------------------------------
databaseup:
	docker compose -f deployments/docker-compose.yaml up -d

databasedown:
	docker compose -f deployments/docker-compose.yaml down

# ------------------- Read schema sql -> crete or update database --------------
# Migarte database all
migrateup:
	migrate -path internal/dataaccess/database/migration -database "postgresql://root:secret@$(HOST):5432/engineer_pro?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/dataaccess/database/migration -database "postgresql://root:secret@$(HOST):5432/engineer_pro?sslmode=disable" -verbose down

# Migarte database lastest
migrateup1:
	migrate -path internal/dataaccess/database/migration -database "postgresql://root:secret@$(HOST):5432/engineer_pro?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path internal/dataaccess/database/migration -database "postgresql://root:secret@$(HOST):5432/engineer_pro?sslmode=disable" -verbose down 1

# ------------------- Read schema and query sqlc -> generate code golang -------
# sqlc gen code golang
sqlc:
	sqlc generate -f ./configs/sqlc.yaml

# Start server http
server:
	go run cmd/main.go

.PHONY: databaseup databasedown migrateup migratedown migrateup1 migratedown1 sqlc server
