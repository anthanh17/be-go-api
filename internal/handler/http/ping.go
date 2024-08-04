package http

import (
	"net/http"
	"time"

	"github.com/anthanh17/be-go-api/internal/handler/token"
	"github.com/gin-gonic/gin"
)

func (s *Server) Ping(ctx *gin.Context) {
	// Get data by access token
	accessPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Get ping_lock_key
	lockKey := "ping_lock:" + accessPayload.ID.String()

	// Check and set `ping_lock`: using `setnx`
	ok, err := s.sessionCache.SetPingLock(ctx, lockKey, "locked")
	if err != nil {
		s.logger.Info("failed to check and set `ping_lock`")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// If Set ping_lock done (ok == true)
	if ok {
		defer func() {
			// Free `ping_lock` once done
			err := s.sessionCache.Del(ctx, lockKey)
			if err != nil {
				s.logger.Info("Error deleting lock:" + err.Error())
			}
		}()

		// Handle API, include sleep
		time.Sleep(5 * time.Second)
		ctx.JSON(http.StatusOK, gin.H{"message": "Pong"})
	} else {
		// If the lock cannot be set (API is locked)
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "API is currently in use"})
	}
}
