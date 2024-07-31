package main

import (
	"context"
	"ep-golang-caching/configs"
	db "ep-golang-caching/internal/dataaccess/database/sqlc"
	"ep-golang-caching/internal/handler/http"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Connect postgress database
	connPool, err := pgxpool.New(context.Background(), config.Database.Source)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	store := db.NewStore(connPool)

	// Run server HTTP
	runGinServer(config, store)
}

func runGinServer(config configs.Config, store db.Store) {
	server, err := http.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}

	err = server.Start(config.HTTP.Address)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
