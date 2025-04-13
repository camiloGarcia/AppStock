// internal/model/stock.go
package model

type Stock struct {
	Ticker        string  `json:"ticker"`
	Company       string  `json:"company"`
	Brokerage     string  `json:"brokerage"`
	Action        string  `json:"action"`
	RatingFrom    string  `json:"rating_from"`
	RatingTo      string  `json:"rating_to"`
	TargetFrom    float64 // ya convertido
	TargetTo      float64 // ya convertido
	TargetFromRaw string  `json:"target_from"` // original
	TargetToRaw   string  `json:"target_to"`   // original
	Time          string  `json:"time"`
}

type PaginatedStocksResponse struct {
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
	Items []Stock `json:"items"`
}
