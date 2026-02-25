// Package api 提供 HTTP 接口与 CORS 处理。
package api

import (
	"encoding/json"
	"net/http"
	"fortune-teller-app/internal/fortune"
)

// Health 返回服务健康状态。
// GET /api/health
func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DrawFortune 随机抽取一条运势（可选 query: category=love|career|health|wealth|general）。
// GET /api/fortune
func DrawFortune(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	var item fortune.Item
	if category != "" {
		item = fortune.DrawByCategory(fortune.Category(category))
	} else {
		item = fortune.DrawRandom()
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(item)
}

// ListCategories 返回所有运势类别。
// GET /api/categories
func ListCategories(w http.ResponseWriter, _ *http.Request) {
	all := fortune.All()
	names := make([]struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	}, 0, len(all))
	for id, items := range all {
		if len(items) > 0 {
			names = append(names, struct {
				ID    string `json:"id"`
				Label string `json:"label"`
			}{string(id), id.Label()})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"categories": names})
}
