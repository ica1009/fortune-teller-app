// Package main 启动算命 API 服务，提供运势抽取等接口，监听 8081。
package main

import (
	"log"
	"net/http"
	"fortune-teller-app/internal/api"
)

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func main() {
	http.HandleFunc("/api/health", cors(api.Health))
	http.HandleFunc("/api/fortune", cors(api.DrawFortune))
	http.HandleFunc("/api/categories", cors(api.ListCategories))
	log.Println("fortune-teller-api listening on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
