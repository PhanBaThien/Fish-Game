package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDBWithAutomation sử dụng pgxpool native để kiểm tra và tạo DB tự động
func InitDBWithAutomation(targetDSN, seedFilePath string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	// 1. Phân tích DSN để lấy tên Database mục tiêu
	u, err := url.Parse(targetDSN)
	if err != nil {
		return nil, fmt.Errorf("DSN không hợp lệ: %v", err)
	}

	targetDBName := strings.TrimPrefix(u.Path, "/")
	if targetDBName == "" || targetDBName == "postgres" {
		return connectAndConfig(targetDSN)
	}

	// 2. Kết nối tạm vào database 'postgres' sử dụng pgx.Connect (rất nhẹ)
	tempURL := *u
	tempURL.Path = "/postgres"

	rootConn, err := pgx.Connect(ctx, tempURL.String())
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối vào root db: %v", err)
	}
	defer rootConn.Close(ctx)

	// 3. Kiểm tra xem targetDBName đã tồn tại chưa
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", targetDBName)
	err = rootConn.QueryRow(ctx, checkQuery).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("lỗi kiểm tra tồn tại db: %v", err)
	}

	isNewDB := false
	if !exists {
		log.Printf("🛠 Database '%s' chưa có, đang tạo mới...", targetDBName)
		_, err = rootConn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", targetDBName))
		if err != nil {
			return nil, fmt.Errorf("lỗi tạo database: %v", err)
		}
		isNewDB = true
	}

	// 4. Kết nối chính thức bằng pgxpool
	dbPool, err := connectAndConfig(targetDSN)
	if err != nil {
		return nil, err
	}

	// 5. Nếu là DB mới tạo, thực thi file seed
	if isNewDB {
		log.Printf("📜 Đang khởi tạo dữ liệu từ: %s", seedFilePath)
		if err := runSeed(dbPool, seedFilePath); err != nil {
			dbPool.Close()
			log.Printf("⚠️ Seed lỗi, đang xóa database '%s' để dọn dẹp...", targetDBName)
			_, _ = rootConn.Exec(ctx, fmt.Sprintf("DROP DATABASE %s WITH (FORCE)", targetDBName))
			return nil, fmt.Errorf("seed lỗi: %v", err)
		}
		log.Println("🚀 Khởi tạo database và chạy seed thành công!")
	}

	return dbPool, nil
}

func connectAndConfig(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("lỗi cấu hình DSN: %v", err)
	}

	// Thiết lập các thông số Pool tối ưu
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("lỗi ping db: %v", err)
	}
	return pool, nil
}

func runSeed(pool *pgxpool.Pool, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("không tìm thấy file seed tại %s: %v", path, err)
	}
	_, err = pool.Exec(context.Background(), string(content))
	return err
}
