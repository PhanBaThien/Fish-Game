package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

// DB is a thin wrapper around *sql.DB to expose the pool for dependency injection.
type DB struct {
	*sql.DB
}

// Config holds all SQL Server connection parameters read from environment variables.
type Config struct {
	Host     string
	Port     int
	User     string // Leave empty to use Windows Authentication (Trusted Connection)
	Password string
	DBName   string
	// Encrypt: "disable" | "false" | "true"
	Encrypt                string
	TrustServerCertificate bool
	// WindowsAuth is true when DB_USER is not set — uses Integrated Security (SSPI)
	WindowsAuth bool
}

// configFromEnv reads DB_* environment variables and returns a Config.
func configFromEnv() Config {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil || port == 0 {
		port = 1433 // SQL Server default port
	}

	trustCert := true
	if v := os.Getenv("DB_TRUST_SERVER_CERT"); v == "false" {
		trustCert = false
	}

	encrypt := os.Getenv("DB_ENCRYPT")
	if encrypt == "" {
		encrypt = "disable"
	}

	user := os.Getenv("DB_USER")
	return Config{
		Host:                   os.Getenv("DB_HOST"),
		Port:                   port,
		User:                   user,
		Password:               os.Getenv("DB_PASSWORD"),
		DBName:                 os.Getenv("DB_NAME"),
		Encrypt:                encrypt,
		TrustServerCertificate: trustCert,
		WindowsAuth:            user == "", // no user → use Windows Integrated Security
	}
}

// dsn builds the SQL Server connection string.
//
// Windows Auth (Trusted Connection):  sqlserver://GIGABYTE?database=MyStoreDB&integrated+security=SSPI&...
// SQL Auth:                           sqlserver://user:pass@host:port?database=...&...
func (c Config) dsn() string {
	trustVal := "false"
	if c.TrustServerCertificate {
		trustVal = "true"
	}

	if c.WindowsAuth {
		// Windows Integrated Security — no username/password needed.
		// go-mssqldb supports "integrated security=SSPI" for Windows Auth.
		return fmt.Sprintf(
			"sqlserver://%s:%d?database=%s&integrated+security=SSPI&encrypt=%s&TrustServerCertificate=%s",
			c.Host, c.Port, c.DBName, c.Encrypt, trustVal,
		)
	}

	// SQL Server Authentication (username + password)
	return fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s&encrypt=%s&TrustServerCertificate=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Encrypt, trustVal,
	)
}

// Connect opens a SQL Server connection pool and verifies the connection with a Ping.
// Returns a *DB wrapper for dependency injection throughout the application.
func Connect() (*DB, error) {
	cfg := configFromEnv()

	// DB_HOST and DB_NAME are always required.
	// DB_USER is optional — when empty, Windows Authentication (SSPI) is used.
	if cfg.Host == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("database: missing required environment variables (DB_HOST, DB_NAME)")
	}
	if !cfg.WindowsAuth && cfg.User == "" {
		return nil, fmt.Errorf("database: DB_USER is required when not using Windows Authentication")
	}

	sqlDB, err := sql.Open("sqlserver", cfg.dsn())
	if err != nil {
		return nil, fmt.Errorf("database: sql.Open failed: %w", err)
	}

	// Connection pool tuning
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	// Verify connectivity with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("database: ping failed (host=%s:%d db=%s): %w", cfg.Host, cfg.Port, cfg.DBName, err)
	}

	authMode := "SQL Auth"
	if cfg.WindowsAuth {
		authMode = "Windows Auth (Trusted Connection)"
	}
	log.Printf("✅ Connected to SQL Server — host: %s:%d | db: %s | auth: %s", cfg.Host, cfg.Port, cfg.DBName, authMode)
	return &DB{sqlDB}, nil
}

// MustConnect calls Connect and fatally exits if connection fails.
// Use only in main() during application startup.
func MustConnect() *DB {
	db, err := Connect()
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	return db
}

// Close gracefully closes all connections in the pool.
func (db *DB) Close() {
	if err := db.DB.Close(); err != nil {
		log.Printf("⚠️  Error closing database connection: %v", err)
		return
	}
	log.Println("🔒 Database connection pool closed.")
}
