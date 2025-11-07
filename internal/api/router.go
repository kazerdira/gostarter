package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/yourusername/go-sqlc-starter/internal/api/handlers"
	"github.com/yourusername/go-sqlc-starter/internal/api/middleware"
	"github.com/yourusername/go-sqlc-starter/internal/auth"
	"github.com/yourusername/go-sqlc-starter/internal/config"
	"github.com/yourusername/go-sqlc-starter/internal/db/sqlc"
)

// NewRouter creates and configures the application router
func NewRouter(cfg *config.Config, db *sql.DB, logger zerolog.Logger) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	router.Use(middleware.RequestID())

	// Initialize dependencies
	queries := sqlc.New(db)
	jwtManager := auth.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTAccessExpiry,
		cfg.JWTRefreshExpiry,
	)

	// Health check endpoints (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"env":    cfg.Env,
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		// Check database connection
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  "database unavailable",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public authentication routes
		authHandler := handlers.NewAuthHandler(queries, jwtManager, db)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected user routes
		userHandler := handlers.NewUserHandler(queries)
		users := v1.Group("/users")
		users.Use(middleware.AuthRequired(jwtManager))
		{
			users.GET("/me", userHandler.GetCurrentUser)
			users.PUT("/me", userHandler.UpdateCurrentUser)
			users.DELETE("/me", userHandler.DeleteCurrentUser)
			
			// Admin only routes
			users.GET("/:id", middleware.AdminRequired(), userHandler.GetUserByID)
			users.GET("", middleware.AdminRequired(), userHandler.ListUsers)
		}
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "route not found",
		})
	})

	return router
}
