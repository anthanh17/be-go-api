package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/anthanh17/be-go-api/internal/handler/token"
	"github.com/gin-gonic/gin"
)

func (s *Server) ping(ctx *gin.Context) {
	// Get data by access token
	accessPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	/*
	* 1. Count the number of times a person calls the api /ping
	 */
	pingCountKey := "ping_counter:" + accessPayload.Username

	// Get value ping_counter cache
	counter := 0
	count, err := s.sessionCache.Get(ctx, pingCountKey)
	if err == nil {
		fmt.Println("ooo")
		num, err := strconv.Atoi(count)
		if err != nil {
			s.logger.Info("Conversion error:" + err.Error())
		} else {
			fmt.Println("Helele")
			counter = num
		}
	}

	// count 1 unit
	counter++

	// Set value ping_counter cache
	err = s.sessionCache.Set(ctx, pingCountKey, counter)
	if err != nil {
		s.logger.Info("failed - set value ping_counter cache`")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	/*
	* 2. The /ping API only allows 1 caller at a time
	* (with sleep inside that api for 5 seconds).
	 */
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
		time.Sleep(10 * time.Second)
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	} else {
		// If the lock cannot be set (API is locked)
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "API is currently in use"})
	}
}
