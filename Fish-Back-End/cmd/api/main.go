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

func mustDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("❌ Giá trị %s không hợp lệ: %v", key, err)
	}
	return d
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Không tìm thấy file .env, sử dụng biến hệ thống")
	} 

	db, err := database.InitDBWithAutomation(os.Getenv("DATABASE_URL"), "internal/scripts/sql/seed.sql")
	if err != nil {
		log.Fatalf("❌ Lỗi khởi tạo Database: %v", err)
	}
	defer db.Close()

	tokenMaker := utils.NewTokenMaker(
		os.Getenv("ACCESS_TOKEN_KEY"),
		os.Getenv("REFRESH_TOKEN_KEY"),
		mustDuration("ACCESS_TOKEN_EXPIRY", 15*time.Minute),
		mustDuration("REFRESH_TOKEN_EXPIRY", 168*time.Hour),
		jwt.SigningMethodHS256,
	)

	allHandlers, err := InitializeApp(db, utils.NewPasswordHasher(), tokenMaker)
	if err != nil {
		log.Fatalf("❌ Lỗi khởi tạo Dependencies (Wire): %v", err)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      authHttp.SetupRouter(allHandlers),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("🐟 Fish Game Server đang chạy tại port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Lỗi server: %v", err)
		}
	}()

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
