package api

import (
	"appstock/internal/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// 🔹 Test para getActionScore
func TestGetActionScore(t *testing.T) {
	cases := []struct {
		input string
		want  int
	}{
		{"Upgraded by", 3},
		{"Reiterated by", 2},
		{"Initiated by", 2},
		{"Downgraded by", 1},
		{"Unknown", 0},
	}

	for _, c := range cases {
		got := getActionScore(c.input)
		if got != c.want {
			t.Errorf("getActionScore(%q) == %d, want %d", c.input, got, c.want)
		}
	}
}

// 🔹 Test para getRatingScore
func TestGetRatingScore(t *testing.T) {
	cases := []struct {
		from, to string
		want     int
	}{
		{"Sell", "Buy", 11},
		{"Hold", "Neutral", 1},
		{"Buy", "Sell", -11},
		{"Unknown", "Buy", 11}, // Unknown default to 0
	}

	for _, c := range cases {
		got := getRatingScore(c.from, c.to)
		if got != c.want {
			t.Errorf("getRatingScore(%q → %q) == %d, want %d", c.from, c.to, got, c.want)
		}
	}
}

func TestRecommendMultipleStocksHandler(t *testing.T) {
	// Fecha de prueba fija
	now := time.Now().UTC()

	// Simula respuesta del "repositorio"
	getStocksByLocalDateRangeMock = func(start, end time.Time) ([]model.Stock, error) {
		return []model.Stock{
			{
				Ticker:     "AAPL",
				Company:    "Apple Inc.",
				Brokerage:  "JP Morgan",
				Action:     "Upgraded by",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				TargetFrom: 150,
				TargetTo:   180,
				Time:       now.Format(time.RFC3339),
			},
			{
				Ticker:     "TSLA",
				Company:    "Tesla",
				Brokerage:  "Morgan Stanley",
				Action:     "Reiterated by",
				RatingFrom: "Sell",
				RatingTo:   "Hold",
				TargetFrom: 100,
				TargetTo:   110,
				Time:       now.Format(time.RFC3339),
			},
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/recommendation?date="+now.Format("2006-01-02"), nil)
	w := httptest.NewRecorder()

	RecommendMultipleStocks(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var result []model.Stock
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 recommendations, got %d", len(result))
	}

	if !strings.EqualFold(result[0].Ticker, "AAPL") {
		t.Errorf("Expected first recommendation to be AAPL, got %s", result[0].Ticker)
	}
}
