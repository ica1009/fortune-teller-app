// Package auth 提供注册、登录与 JWT；密码仅 bcrypt 哈希存储与校验。
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"fortune-teller-app/internal/store"
)

// ErrInvalidCredentials 登录失败（不区分用户不存在与密码错误）。
var ErrInvalidCredentials = errors.New("invalid username or password")

const (
	bcryptCost     = 10
	usernameMinLen = 3
	usernameMaxLen = 32
	passwordMinLen = 6
)

// Register 注册：校验用户名与唯一性，密码 bcrypt 后入库。
func Register(ctx context.Context, db *store.DB, username, password string) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	_, err = db.CreateUser(ctx, username, string(hash))
	return err
}

// Login 校验用户名与密码，成功签发 JWT。
func Login(ctx context.Context, db *store.DB, username, password string, secret []byte) (string, error) {
	u, err := db.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}
	return issueJWT(u.ID, u.Username, secret)
}

// Claims JWT 载荷。
type Claims struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

func issueJWT(userID int64, username string, secret []byte) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserID:   userID,
		Username: username,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func validateUsername(s string) error {
	if len(s) < usernameMinLen || len(s) > usernameMaxLen {
		return errors.New("username must be 3–32 characters")
	}
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			continue
		}
		return errors.New("username may only contain letters, digits and underscore")
	}
	return nil
}

func validatePassword(s string) error {
	if len(s) < passwordMinLen {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}
