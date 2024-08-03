URL = localhost
# URL = ec2-18-141-12-199.ap-southeast-1.compute.amazonaws.com

# ----------------------------- Setup database ---------------------------------
database:
	docker-compose -f ./deployments/docker-compose.yaml up

# ------------------- Read schema sql -> crete or update database --------------
# Migarte database all
migrateup:
	migrate -path internal/dataaccess/database/migration -database "postgresql://root:secret@$(URL):5432/engineer_pro?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/dataaccess/database/migration -database "postgresql://root:secret@$(URL):5432/engineer_pro?sslmode=disable" -verbose down

# Migarte database lastest
migrateup1:
	migrate -path internal/dataaccess/database/migration -database "postgresql://root:secret@$(URL):5432/engineer_pro?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path internal/dataaccess/database/migration -database "postgresql://root:secret@$(URL):5432/engineer_pro?sslmode=disable" -verbose down 1

# ------------------- Read schema and query sqlc -> generate code golang -------
# sqlc gen code golang
sqlc:
	sqlc generate -f ./configs/sqlc.yaml

# Start server http
server:
	go run cmd/main.go

.PHONY: database migrateup migratedown migrateup1 migratedown1 sqlc server
