package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-sqlc-starter/internal/auth"
	"github.com/yourusername/go-sqlc-starter/internal/db/sqlc"
)

type AuthHandler struct {
	queries    *sqlc.Queries
	jwtManager *auth.JWTManager
	db         *sql.DB
}

func NewAuthHandler(queries *sqlc.Queries, jwtManager *auth.JWTManager, db *sql.DB) *AuthHandler {
	return &AuthHandler{
		queries:    queries,
		jwtManager: jwtManager,
		db:         db,
	}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents the refresh token request body
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

// UserInfo represents basic user information
type UserInfo struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	IsAdmin  bool   `json:"is_admin"`
}

// Register creates a new user account
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user, err := h.queries.CreateUser(c.Request.Context(), sqlc.CreateUserParams{
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
	})
	if err != nil {
		// Check if email already exists
		if err.Error() == "duplicate key value violates unique constraint" {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Generate tokens
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	// Store refresh token
	expiresAt := time.Now().Add(168 * time.Hour) // 7 days
	_, err = h.queries.CreateRefreshToken(c.Request.Context(), sqlc.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store refresh token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
		User: UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			IsAdmin:  user.IsAdmin,
		},
	})
}

// Login authenticates a user
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.queries.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	}

	// Verify password
	if err := auth.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Generate tokens
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	// Store refresh token
	expiresAt := time.Now().Add(168 * time.Hour)
	_, err = h.queries.CreateRefreshToken(c.Request.Context(), sqlc.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store refresh token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
		User: UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			IsAdmin:  user.IsAdmin,
		},
	})
}

// RefreshToken issues a new access token using a refresh token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	userID, err := h.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// Check if refresh token exists in database
	storedToken, err := h.queries.GetRefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found or expired"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify refresh token"})
		return
	}

	// Get user
	user, err := h.queries.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Generate new access token
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	// Optionally rotate refresh token (recommended for security)
	newRefreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	// Delete old refresh token
	if err := h.queries.DeleteRefreshToken(c.Request.Context(), storedToken.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rotate refresh token"})
		return
	}

	// Store new refresh token
	expiresAt := time.Now().Add(168 * time.Hour)
	_, err = h.queries.CreateRefreshToken(c.Request.Context(), sqlc.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store refresh token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
		User: UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			IsAdmin:  user.IsAdmin,
		},
	})
}

// Logout invalidates a user's refresh token
func (h *AuthHandler) Logout(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete refresh token
	if err := h.queries.DeleteRefreshToken(c.Request.Context(), req.RefreshToken); err != nil {
		// Don't expose if token doesn't exist
		c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
