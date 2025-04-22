package api

import (
	"appstock/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStocks_Success(t *testing.T) {
	// 🧪 Mock del repositorio
	getStocksWithCount = func(page, limit int, search, sortBy, sortDir string) ([]model.Stock, int, error) {
		return []model.Stock{
			{
				Ticker:    "AAPL",
				Company:   "Apple Inc.",
				Brokerage: "Goldman Sachs",
			},
		}, 1, nil
	}

	// 🧪 Petición simulada
	req := httptest.NewRequest("GET", "/stocks?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	// 🔧 Ejecutar handler
	GetStocks(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}

	var result model.PaginatedStocksResponse
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding response: %v", err)
	}

	if len(result.Items) != 1 || result.Items[0].Ticker != "AAPL" {
		t.Errorf("Unexpected result: %+v", result.Items)
	}
}

func TestGetStocks_Failure(t *testing.T) {
	// 🧪 Simula error del repositorio
	getStocksWithCount = func(page, limit int, search, sortBy, sortDir string) ([]model.Stock, int, error) {
		return nil, 0, errors.New("DB error")
	}

	req := httptest.NewRequest("GET", "/stocks", nil)
	w := httptest.NewRecorder()

	GetStocks(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected 500 Internal Server Error, got %d", resp.StatusCode)
	}
}
