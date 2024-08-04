package http

import (
	"fmt"
	"time"

	"github.com/anthanh17/be-go-api/configs"
	"github.com/anthanh17/be-go-api/internal/dataaccess/cache"
	db "github.com/anthanh17/be-go-api/internal/dataaccess/database/sqlc"
	"github.com/anthanh17/be-go-api/internal/handler/token"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server serves HTTP requests for our service.
type Server struct {
	config     configs.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	sessionCache cache.SessionCache
	logger     *zap.Logger
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config configs.Config, store db.Store, cachier cache.Cachier, logger *zap.Logger) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.Token.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		sessionCache: cache.NewSessionCache(cachier, logger),
		logger:     logger,
	}

	// Router
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

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.GET("/ping", server.Ping)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
