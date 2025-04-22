package service

import (
	"appstock/internal/model"
	"appstock/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Simula valores guardados
var savedStocks []model.Stock

// Mock de SaveStock
var mockSaveStock = func(s model.Stock) error {
	savedStocks = append(savedStocks, s)
	return nil
}

func TestFetchAndStoreAllStocks_Success(t *testing.T) {
	// Respuesta simulada de la API
	mockResponse := ApiResponse{
		Items: []model.Stock{
			{
				Ticker:        "AAPL",
				Company:       "Apple Inc.",
				Brokerage:     "Goldman Sachs",
				Action:        "upgraded by",
				RatingFrom:    "neutral",
				RatingTo:      "buy",
				TargetFromRaw: "150",
				TargetToRaw:   "180",
			},
		},
		NextPage: "",
	}

	// Simula servidor HTTP de prueba
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonBytes, _ := json.Marshal(mockResponse)
		w.Write(jsonBytes)
	}))
	defer server.Close()

	// Redefine variables globales del paquete
	apiURL = server.URL
	apiKey = "mock-key"
	bearerToken = fmt.Sprintf("Bearer %s", apiKey)

	// Redirige SaveStock al mock
	originalSaveStock := repository.SaveStock
	repository.SaveStock = mockSaveStock
	defer func() { repository.SaveStock = originalSaveStock }()

	// Limpia estado previo
	savedStocks = []model.Stock{}

	// Ejecuta función
	err := FetchAndStoreAllStocks()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verifica
	if len(savedStocks) != 1 {
		t.Fatalf("Expected 1 stock saved, got %d", len(savedStocks))
	}

	s := savedStocks[0]
	if s.Ticker != "AAPL" {
		t.Errorf("Expected Ticker 'AAPL', got '%s'", s.Ticker)
	}

}
