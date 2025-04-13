package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"appstock/internal/model"
	"appstock/internal/repository"
)

func GetStocks(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
		limit = l
	}

	stocks, total, err := repository.GetStocksPaginatedWithCount(page, limit)
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
