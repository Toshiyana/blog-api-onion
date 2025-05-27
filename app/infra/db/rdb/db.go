package rdb

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DB : データベース接続を管理する構造体
type DB struct {
	db *sqlx.DB
}

// NewDB : DBの生成
func NewDB() (*DB, error) {
	// 環境変数から接続情報を取得
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "password")
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	database := getEnv("DB_NAME", "myblog")

	// DSN (Data Source Name) の構築
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo", user, password, host, port, database)

	// データベース接続
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続確認
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// コネクションプールの設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &DB{db: db}, nil
}

// Close : データベース接続のクローズ
func (d *DB) Close() error {
	return d.db.Close()
}

// Begin : トランザクションの開始
func (d *DB) Begin() (*sqlx.Tx, error) {
	return d.db.Beginx()
}

// Read : 読み取り用のデータベース接続を取得
func (d *DB) Read(ctx context.Context) *sqlx.DB {
	return d.db
}

// Write : 書き込み用のデータベース接続を取得
func (d *DB) Write(ctx context.Context) *sqlx.DB {
	return d.db
}

// getEnv : 環境変数を取得（デフォルト値付き）
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// TransactionManager : トランザクション管理インターフェース
type TransactionManager interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// DefaultTransactionManager : デフォルトのトランザクション管理実装
type DefaultTransactionManager struct {
	db *DB
}

// NewDefaultTransactionManager : DefaultTransactionManagerの生成
func NewDefaultTransactionManager(db *DB) TransactionManager {
	return &DefaultTransactionManager{db: db}
}

// Transaction : トランザクションを実行
func (tm *DefaultTransactionManager) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// トランザクションをコンテキストに設定
	txCtx := context.WithValue(ctx, txKey{}, tx)

	// 関数実行
	if err := fn(txCtx); err != nil {
		// エラー発生時はロールバック
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w (original error: %v)", rbErr, err)
		}
		return err
	}

	// コミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// トランザクションをコンテキストに格納するためのキー
type txKey struct{}

// GetTx : コンテキストからトランザクションを取得
func GetTx(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sqlx.Tx)
	return tx, ok
}
