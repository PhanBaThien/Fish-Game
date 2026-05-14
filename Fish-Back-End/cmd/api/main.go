package main

import (
	"context"
	"log"
	"net/http" // Thư viện chuẩn của Go
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/database"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/middleware"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/repository"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/usecase"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"

	// Sử dụng alias authHttp để tránh trùng tên với net/http
	authHttp "github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/transport/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// 1. Nạp biến môi trường từ file .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Không tìm thấy file .env, sử dụng biến môi trường hệ thống")
	}

	// 2. Tự động hóa khởi tạo Database
	rootDSN := os.Getenv("DATABASE_URL")
	seedPath := "internal/scripts/sql/seed.sql"

	// Gọi hàm khởi tạo thông minh (đảm bảo bạn đã cập nhật file postgre.go với hàm này)
	db, err := database.InitDBWithAutomation(rootDSN, seedPath)
	if err != nil {
		log.Fatalf("❌ Lỗi khởi tạo Database: %v", err)
	}
	defer db.Close()

	// 3. Khởi tạo Utilities
	hasher := utils.NewPasswordHasher()
	tokenMaker := utils.NewTokenMaker(os.Getenv("TOKEN_SYMMETRIC_KEY"), jwt.SigningMethodHS256)

	// 4. Dependency Injection (Lắp ráp các tầng)
	adminRepo := repository.NewAdminRepository(db)
	authUsecase := usecase.NewAuthUsecase(adminRepo, hasher, tokenMaker)

	// Sử dụng alias authHttp và đặt tên biến là authHdl để không trùng lặp
	authHdl := authHttp.NewAuthHandler(authUsecase, tokenMaker)

	// 5. Cấu hình Gin Router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// 6. Đăng ký Routes
	v1 := router.Group("/api/v1")
	{
		// Tự động đăng ký các route /login, /register, /logout, /me từ handler
		authHdl.RegisterRoutes(v1)
	}

	// 7. Khởi chạy Server với Graceful Shutdown
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Chạy server trong một goroutine để không chặn luồng chính
	go func() {
		log.Printf("🐟 Fish Game Server is running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Lỗi server: %v", err)
		}
	}()

	// Chờ tín hiệu kết thúc (Ctrl+C hoặc kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Đang đóng server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("❌ Server buộc phải đóng:", err)
	}
	log.Println("✅ Server đã dừng an toàn.")
}
