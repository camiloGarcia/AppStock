// appstock/internal/api/stocks.go
package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"appstock/internal/model"
	"appstock/internal/repository"
)

// Delegable para test
var getStocksWithCount = repository.GetStocksPaginatedWithCount

func GetStocks(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")

	page := 1
	limit := 10
	if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
		limit = l
	}

	// stocks, total, err := repository.GetStocksPaginatedWithCount(page, limit, search, sortBy, sortDir)
	stocks, total, err := getStocksWithCount(page, limit, search, sortBy, sortDir)
	if err != nil {
		http.Error(w, "Error fetching stocks", http.StatusInternalServerError)
		return
	}

	response := model.PaginatedStocksResponse{
		Total: total,
		Page:  page,
		Limit: limit,
		Items: stocks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
