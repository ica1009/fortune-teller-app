// Package store 提供 PostgreSQL 用户存储；密码仅存哈希，不暴露。
package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB 封装 PG 连接池。
type DB struct {
	Pool *pgxpool.Pool
}

// User 内部校验用；password_hash 禁止写入响应或日志。
type User struct {
	ID           int64
	Username     string
	PasswordHash string
}

// Open 根据 dataSourceName 创建连接池。
func Open(ctx context.Context, dataSourceName string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(dataSourceName)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &DB{Pool: pool}, nil
}

// Close 关闭连接池。
func (db *DB) Close() {
	db.Pool.Close()
}

// Migrate 创建 users 表。
func (db *DB) Migrate(ctx context.Context) error {
	_, err := db.Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			username VARCHAR(64) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	`)
	return err
}

// CreateUser 插入用户；用户名已存在返回 ErrDuplicateUsername。
func (db *DB) CreateUser(ctx context.Context, username, passwordHash string) (int64, error) {
	var id int64
	err := db.Pool.QueryRow(ctx,
		`INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`,
		username, passwordHash,
	).Scan(&id)
	if err != nil && isDuplicateKey(err) {
		return 0, ErrDuplicateUsername
	}
	return id, err
}

// GetUserByUsername 按用户名查询，仅用于登录校验。
func (db *DB) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	err := db.Pool.QueryRow(ctx,
		`SELECT id, username, password_hash FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// ErrDuplicateUsername 用户名已存在。
var ErrDuplicateUsername = errors.New("username already exists")

func isDuplicateKey(err error) bool {
	var pge *pgconn.PgError
	return errors.As(err, &pge) && pge.Code == "23505"
}
