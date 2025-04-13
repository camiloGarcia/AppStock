// internal/repository/stock_repo.go
package repository

import (
	"log"

	"appstock/internal/model"
	"appstock/pkg/db"
)

func SaveStock(s model.Stock) error {
	query := `
		INSERT INTO stocks (
			ticker, company, brokerage, action,
			rating_from, rating_to, target_from, target_to, time
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`
	_, err := db.DB.Exec(query,
		s.Ticker, s.Company, s.Brokerage, s.Action,
		s.RatingFrom, s.RatingTo, s.TargetFrom, s.TargetTo, s.Time,
	)
	if err != nil {
		log.Printf("Error saving stock %s: %v", s.Ticker, err)
	}
	return err
}

func GetAllStocks() ([]model.Stock, error) {
	rows, err := db.DB.Query(`
		SELECT ticker, company, brokerage, action,
		       rating_from, rating_to, target_from, target_to, time
		FROM stocks
		ORDER BY time DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []model.Stock
	for rows.Next() {
		var s model.Stock
		err := rows.Scan(&s.Ticker, &s.Company, &s.Brokerage, &s.Action,
			&s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &s.Time)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
}

func HasStocks() (bool, error) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM stocks").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetStocksPaginated(page, limit int) ([]model.Stock, error) {
	offset := (page - 1) * limit

	query := `
		SELECT ticker, company, brokerage, action,
		       rating_from, rating_to, target_from, target_to, time
		FROM stocks
		ORDER BY time DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := db.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []model.Stock
	for rows.Next() {
		var s model.Stock
		err := rows.Scan(&s.Ticker, &s.Company, &s.Brokerage, &s.Action,
			&s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &s.Time)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
}

func GetStocksPaginatedWithCount(page, limit int) ([]model.Stock, int, error) {
	offset := (page - 1) * limit

	// Obtener total de registros
	var total int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM stocks`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Obtener datos paginados
	query := `
		SELECT ticker, company, brokerage, action,
		       rating_from, rating_to, target_from, target_to, time
		FROM stocks
		ORDER BY time DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := db.DB.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var stocks []model.Stock
	for rows.Next() {
		var s model.Stock
		err := rows.Scan(&s.Ticker, &s.Company, &s.Brokerage, &s.Action,
			&s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &s.Time)
		if err != nil {
			return nil, 0, err
		}
		stocks = append(stocks, s)
	}
	return stocks, total, nil
}
