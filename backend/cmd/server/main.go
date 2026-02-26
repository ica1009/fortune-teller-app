// Package main 启动算命 API 服务，提供运势抽取与注册/登录，监听 8081。
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"fortune-teller-app/internal/api"
	"fortune-teller-app/internal/store"
)

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	ctx := context.Background()
	db, err := store.Open(ctx, dbURL)
	if err != nil {
		log.Fatal("open db: ", err)
	}
	defer db.Close()
	if err := db.Migrate(ctx); err != nil {
		log.Fatal("migrate: ", err)
	}

	authHandlers := &api.AuthHandlers{Store: db, Secret: []byte(secret)}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", cors(api.Health))
	mux.HandleFunc("/api/fortune", cors(api.DrawFortune))
	mux.HandleFunc("/api/fortune/ai", cors(api.DrawFortuneAI))
	mux.HandleFunc("/api/categories", cors(api.ListCategories))
	mux.HandleFunc("/api/auth/register", authHandlers.Register)
	mux.HandleFunc("/api/auth/login", authHandlers.Login)

	srv := &http.Server{Addr: ":8081", Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Println("fortune-teller-api listening on :8081")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = srv.Shutdown(context.Background())
}
