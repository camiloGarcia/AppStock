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
	var filterDate time.Time
	if dateStr != "" {
		var err error
		filterDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
	} else {
		filterDate = time.Now()
	}

	// 🕒 Cargar zona horaria local
	loc, err := time.LoadLocation("America/Bogota") // o tu zona horaria real
	if err != nil {
		http.Error(w, "Failed to load timezone", http.StatusInternalServerError)
		return
	}

	// 📦 Obtener todos los registros
	stocks, err := repository.GetAllStocks()
	if err != nil {
		http.Error(w, "Failed to load stocks", http.StatusInternalServerError)
		return
	}

	// 🧠 Filtrar por fecha local
	var todaysStocks []model.Stock
	for _, stock := range stocks {
		parsed, err := time.Parse(time.RFC3339, stock.Time)
		if err != nil {
			continue
		}
		localDate := parsed.In(loc).Format("2006-01-02")
		if localDate == filterDate.Format("2006-01-02") {
			todaysStocks = append(todaysStocks, stock)
		}
	}

	// 🧮 Calcular puntuación
	type ScoredStock struct {
		model.Stock
		Score int
	}
	var scored []ScoredStock
	for _, stock := range todaysStocks {
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
