// Package api 提供注册与登录 HTTP 处理；不返回或记录密码及 password_hash。
package api

import (
	"encoding/json"
	"net/http"

	"fortune-teller-app/internal/auth"
	"fortune-teller-app/internal/store"
)

// AuthHandlers 注册与登录处理器。
type AuthHandlers struct {
	Store  *store.DB
	Secret []byte
}

// RegisterRequest 注册请求体。
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录成功响应。
type LoginResponse struct {
	Token string `json:"token"`
}

// Register 处理 POST /api/auth/register。
func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	if err := auth.Register(r.Context(), h.Store, req.Username, req.Password); err != nil {
		if err == store.ErrDuplicateUsername {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username already exists"})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"message": "registered"})
}

// Login 处理 POST /api/auth/login。
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	token, err := auth.Login(r.Context(), h.Store, req.Username, req.Password, h.Secret)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid username or password"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "login failed"})
		return
	}
	writeJSON(w, http.StatusOK, LoginResponse{Token: token})
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
