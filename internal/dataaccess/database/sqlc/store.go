package db

import (
	"context"
	"fmt"

	"github.com/anthanh17/be-go-api/configs"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Store defines all functions to execute db queries and transactions
// Repository pattern
type Store interface {
	Querier
	// Expand more transactions in the future
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

func InitializeUpDB(databaseConfig configs.DatabaseConfig, logger *zap.Logger) (Store, func(), error) {
	// postgresql://root:secret@localhost:5432/engineer_pro?sslmode=disable
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		databaseConfig.Username,
		databaseConfig.Password,
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.Database)

	// Connect postgress database
	connPool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		logger.Info("cannot connect to db")
		return nil, nil, err
	}

	// Create database accessor
	store := NewStore(connPool)

	cleanup := func() {
		connPool.Close()
	}

	return store, cleanup, nil
}
