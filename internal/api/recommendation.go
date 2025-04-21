// appstock/internal/api/recommendation.go
package api

import (
	"appstock/internal/model"
	"appstock/internal/repository"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Define una función para obtener la "puntuación" de una acción
func getActionScore(action string) int {
	switch strings.ToLower(action) {
	case "upgraded by":
		return 3
	case "reiterated by", "initiated by":
		return 2
	case "downgraded by":
		return 1
	default:
		return 0
	}
}

// Define una función para obtener el delta de rating (simplificada)
func getRatingScore(from, to string) int {
	rank := map[string]int{
		"sell": 0, "underperform": 1, "underweight": 2,
		"hold": 3, "neutral": 4, "market perform": 5,
		"equal weight": 6, "peer perform": 7, "sector perform": 8,
		"overweight": 9, "outperform": 10, "buy": 11,
		"strong-buy": 12,
	}

	scoreFrom := rank[strings.ToLower(from)]
	scoreTo := rank[strings.ToLower(to)]
	return scoreTo - scoreFrom
}

func RecommendMultipleStocks(w http.ResponseWriter, r *http.Request) {
	// 🔍 Obtener y validar parámetro de fecha
	dateStr := r.URL.Query().Get("date")
	var localDate time.Time
	var err error

	if dateStr != "" {
		localDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
	} else {
		localDate = time.Now()
	}

	// 🕐 Cargar zona horaria local
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		http.Error(w, "Failed to load timezone", http.StatusInternalServerError)
		return
	}

	// 🕓 Convertir a UTC rango [inicio, fin) del día en la zona local
	startOfDay := time.Date(localDate.Year(), localDate.Month(), localDate.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	startUTC := startOfDay.UTC()
	endUTC := endOfDay.UTC()

	// 📦 Obtener registros directamente en el rango UTC
	stocks, err := repository.GetStocksByLocalDateRange(startUTC, endUTC)
	if err != nil {
		http.Error(w, "Failed to load stocks", http.StatusInternalServerError)
		return
	}

	// 🧮 Calcular puntuación
	type ScoredStock struct {
		model.Stock
		Score int
	}
	var scored []ScoredStock
	for _, stock := range stocks {
		actionScore := getActionScore(stock.Action)
		ratingScore := getRatingScore(stock.RatingFrom, stock.RatingTo)
		targetScore := int(stock.TargetTo - stock.TargetFrom)

		totalScore := (actionScore * 2) + ratingScore + targetScore

		scored = append(scored, ScoredStock{
			Stock: stock,
			Score: totalScore,
		})
	}

	// 🔽 Ordenar por puntuación
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	// 🔁 Extraer solo los stocks
	var recommended []model.Stock
	for _, s := range scored {
		recommended = append(recommended, s.Stock)
	}

	// 📨 Enviar respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommended)
}
