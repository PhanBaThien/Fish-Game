package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/database"
	authHttp "github.com/PhanBaThien/Fish-Game/Fish-Back-End/internal/transport/http"
	"github.com/PhanBaThien/Fish-Game/Fish-Back-End/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Nạp biến môi trường
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Không tìm thấy file .env, sử dụng biến hệ thống")
	}

	// 2. Khởi tạo Database (Trọng tâm)
	rootDSN := os.Getenv("DATABASE_URL")
	seedPath := "internal/scripts/sql/seed.sql"

	db, err := database.InitDBWithAutomation(rootDSN, seedPath)
	if err != nil {
		log.Fatalf("❌ Lỗi khởi tạo Database: %v", err)
	}
	defer db.Close()

	hasher := utils.NewPasswordHasher()
	tokenMaker := utils.NewTokenMaker(os.Getenv("TOKEN_SYMMETRIC_KEY"), jwt.SigningMethodHS256)

	allHandlers, err := InitializeApp(db, hasher, tokenMaker)
	if err != nil {
		log.Fatalf("❌ Lỗi khởi tạo Dependencies (Wire): %v", err)
	}

	// 5. Ném cục Struct Handler đó vào Router
	router := authHttp.SetupRouter(allHandlers)

	// 6. Khởi chạy HTTP Server
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

	go func() {
		log.Printf("🐟 Fish Game Server is running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Lỗi server: %v", err)
		}
	}()

	// 7. Graceful Shutdown
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