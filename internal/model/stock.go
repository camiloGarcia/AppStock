// internal/model/stock.go
package model

type Stock struct {
	Ticker     string `json:"ticker"`
	Company    string `json:"company"`
	Brokerage  string `json:"brokerage"`
	Action     string `json:"action"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`

	TargetFrom float64 `json:"target_from"`
	TargetTo   float64 `json:"target_to"`

	TargetFromRaw string `json:"-"` // <- este ya no se serializa
	TargetToRaw   string `json:"-"`
	Time          string `json:"time"`
}

type PaginatedStocksResponse struct {
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
	Items []Stock `json:"items"`
}
