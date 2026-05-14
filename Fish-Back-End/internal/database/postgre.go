package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// InitDBWithAutomation nhận vào DSN mục tiêu và thực hiện kiểm tra/tạo DB tự động
func InitDBWithAutomation(targetDSN, seedFilePath string) (*sql.DB, error) {
	// 1. Phân tích DSN để lấy tên Database mục tiêu
	u, err := url.Parse(targetDSN)
	if err != nil {
		return nil, fmt.Errorf("DSN không hợp lệ: %v", err)
	}

	// Lấy tên DB từ Path (ví dụ "/fish_game" -> "fish_game")
	targetDBName := strings.TrimPrefix(u.Path, "/")
	if targetDBName == "" || targetDBName == "postgres" {
		// Nếu DSN đang trỏ sẵn vào postgres, chỉ cần kết nối bình thường
		return connectAndConfig(targetDSN)
	}

	// 2. Kết nối tạm vào database 'postgres' để kiểm tra/tạo mới
	// Tạo bản sao URL và đổi Path sang /postgres
	tempURL := *u
	tempURL.Path = "/postgres"

	rootDb, err := sql.Open("pgx", tempURL.String())
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối vào root db: %v", err)
	}
	defer rootDb.Close()

	// 3. Kiểm tra xem targetDBName đã tồn tại chưa
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", targetDBName)
	err = rootDb.QueryRow(checkQuery).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("lỗi kiểm tra tồn tại db: %v", err)
	}

	isNewDB := false
	if !exists {
		log.Printf("🛠 Database '%s' chưa có, đang tạo mới...", targetDBName)
		// Lưu ý: Lệnh CREATE DATABASE không chạy được trong transaction
		_, err = rootDb.Exec(fmt.Sprintf("CREATE DATABASE %s", targetDBName))
		if err != nil {
			return nil, fmt.Errorf("lỗi tạo database: %v", err)
		}
		isNewDB = true
	}

	// 4. Kết nối chính thức bằng DSN ban đầu
	db, err := connectAndConfig(targetDSN)
	if err != nil {
		return nil, err
	}

	// 5. Nếu là DB mới tạo, thực thi file seed
	if isNewDB {
		log.Printf("📜 Đang khởi tạo dữ liệu từ: %s", seedFilePath)
		if err := runSeed(db, seedFilePath); err != nil {
			// Nếu lỗi khi chạy seed cho DB mới, dọn dẹp DB để tránh rác
			db.Close()
			log.Printf("⚠️ Seed lỗi, đang xóa database '%s' để dọn dẹp...", targetDBName)
			_, _ = rootDb.Exec(fmt.Sprintf("DROP DATABASE %s WITH (FORCE)", targetDBName))
			return nil, fmt.Errorf("seed lỗi: %v", err)
		}
		log.Println("🚀 Khởi tạo database và chạy seed thành công!")
	}

	return db, nil
}

// connectAndConfig thiết lập các thông số pool cho connection
func connectAndConfig(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("lỗi kết nối db: %v", err)
	}
	return db, nil
}

// runSeed thực thi file SQL khởi tạo
func runSeed(db *sql.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("không tìm thấy file seed tại %s: %v", path, err)
	}
	_, err = db.Exec(string(content))
	return err
}
