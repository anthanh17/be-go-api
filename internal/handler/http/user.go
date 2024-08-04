package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/anthanh17/be-go-api/internal/dataaccess/cache"
	db "github.com/anthanh17/be-go-api/internal/dataaccess/database/sqlc"
	"github.com/anthanh17/be-go-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Info("cannot ShouldBindJSON req")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Info("cannot hashedPassword")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		s.logger.Info("cannot CreateUser")
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Info("cannot ShouldBindJSON req")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		s.logger.Info("cannot GetUser")
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		s.logger.Info("cannot CheckPassword")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		s.config.Token.AccessTokenDuration,
	)
	if err != nil {
		s.logger.Info("cannot accessToken")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		user.Role,
		s.config.Token.RefreshTokenDuration,
	)
	if err != nil {
		s.logger.Info("cannot refreshToken")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create session database
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		s.logger.Info("cannot CreateSession")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create session caching
	sessionID, err := uuid.NewRandom()
	if err != nil {
		s.logger.Info("failed gen sessionID uuid")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sessionData := cache.SessionType{
		SessionID: sessionID.String(),
		Username:  user.Username,
	}
	err = s.sessionCache.SetSession(ctx, accessPayload.ID.String(), sessionData)
	if err != nil {
		s.logger.Info("failed to set session key bytes into cache")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
