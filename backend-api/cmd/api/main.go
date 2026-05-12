package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yourname/fish-game-backend/internal/handlers"
	"github.com/yourname/fish-game-backend/internal/middleware"
	"github.com/yourname/fish-game-backend/internal/repository"
	"github.com/yourname/fish-game-backend/internal/services"
	"github.com/yourname/fish-game-backend/pkg/database"
)

func main() {
	// ─── Load Environment Variables ───────────────────────────────
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ─── Database Connection ───────────────────────────────────────
	db := database.MustConnect()
	defer db.Close()

	// ─── Repositories (Data Access Layer) ─────────────────────────
	authorRepo := repository.NewAuthRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	fishRepo   := repository.NewFishRepository(db)
	roomRepo   := repository.NewRoomRepository(db)

	// ─── Services (Business Logic Layer) ──────────────────────────
	authSvc   := services.NewAuthService(authorRepo)
	playerSvc := services.NewPlayerService(playerRepo)
	fishSvc   := services.NewFishService(fishRepo)
	roomSvc   := services.NewRoomService(roomRepo)

	// ─── Handlers (Presentation Layer) ────────────────────────────
	healthHandler := handlers.NewHealthHandler()
	authHandler   := handlers.NewAuthHandler(authSvc)
	playerHandler := handlers.NewPlayerHandler(playerSvc)
	fishHandler   := handlers.NewFishHandler(fishSvc)
	roomHandler   := handlers.NewRoomHandler(roomSvc)

	// ─── Gin Router ────────────────────────────────────────────────
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// ─── API Routes ────────────────────────────────────────────────
	v1 := router.Group("/api/v1")
	{
		// Health check (public)
		v1.GET("/health", healthHandler.Check)

		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Admin management routes (protected — requires valid JWT)
		admins := v1.Group("/auth/admins")
		admins.Use(middleware.AuthRequired())
		{
			admins.POST("", authHandler.CreateAdmin)
		}

		// Players routes (protected)
		players := v1.Group("/players")
		players.Use(middleware.AuthRequired())
		{
			players.GET("", playerHandler.ListPlayers)
			players.GET("/:id", playerHandler.GetPlayer)
			players.PUT("/:id", playerHandler.UpdatePlayer)
			players.DELETE("/:id", playerHandler.BanPlayer)
		}

		// Fish config routes (protected)
		fish := v1.Group("/fish")
		fish.Use(middleware.AuthRequired())
		{
			fish.GET("", fishHandler.ListFish)
			fish.POST("", fishHandler.CreateFish)
			fish.PUT("/:id", fishHandler.UpdateFish)
			fish.DELETE("/:id", fishHandler.DeleteFish)
		}

		// Rooms routes (protected)
		rooms := v1.Group("/rooms")
		rooms.Use(middleware.AuthRequired())
		{
			rooms.GET("", roomHandler.ListRooms)
			rooms.POST("", roomHandler.CreateRoom)
			rooms.GET("/:id", roomHandler.GetRoom)
			rooms.PUT("/:id", roomHandler.UpdateRoom)
			rooms.DELETE("/:id", roomHandler.CloseRoom)
		}

		// Stats / Dashboard routes (protected)
		stats := v1.Group("/stats")
		stats.Use(middleware.AuthRequired())
		{
			stats.GET("/dashboard", handlers.GetDashboardStats)
		}
	}

	// ─── HTTP Server with Graceful Shutdown ────────────────────────
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🐟 Fish Game API Server listening on port %s (env: %s)", port, appEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
