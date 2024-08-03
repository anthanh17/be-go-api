package http

import (
	"context"
	"ep-golang-caching/configs"
	db "ep-golang-caching/internal/dataaccess/database/sqlc"
	"ep-golang-caching/internal/handler/token"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     configs.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	redisdb    *redis.Client
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config configs.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.Token.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		redisdb: redis.NewClient(&redis.Options{
			Addr: config.Cache.Address,
		}),
	}

	err = server.redisdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("cannot connect redis: %w", err)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Use CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Authorization", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
