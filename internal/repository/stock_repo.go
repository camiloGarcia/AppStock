// internal/repository/stock_repo.go
package repository

import (
	"fmt"
	"log"
	"strings"
	"time"

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
	r := []model.Stock{}
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

	for rows.Next() {
		var s model.Stock
		err := rows.Scan(&s.Ticker, &s.Company, &s.Brokerage, &s.Action,
			&s.RatingFrom, &s.RatingTo, &s.TargetFrom, &s.TargetTo, &s.Time)
		if err != nil {
			return nil, err
		}
		r = append(r, s)
	}
	return r, nil
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

func GetStocksPaginatedWithCount(page, limit int, search, sortBy, sortDir string) ([]model.Stock, int, error) {
	offset := (page - 1) * limit

	whereClause := ""
	orderClause := "ORDER BY time DESC" // default
	var countArgs []interface{}
	var dataArgs []interface{}

	// Valid sortBy fields
	allowedSortFields := map[string]bool{
		"ticker": true, "company": true, "brokerage": true,
		"target_from": true, "target_to": true,
		"rating_from": true, "rating_to": true, "time": true,
	}

	if search != "" {
		pattern := "%" + strings.ToLower(search) + "%"
		whereClause = `WHERE LOWER(ticker) LIKE $1 OR LOWER(company) LIKE $1 OR LOWER(brokerage) LIKE $1`
		countArgs = append(countArgs, pattern)
		dataArgs = append(dataArgs, pattern)
	}

	// Validate sort field
	if allowedSortFields[sortBy] {
		if sortDir != "desc" {
			sortDir = "asc"
		}
		orderClause = fmt.Sprintf("ORDER BY %s %s", sortBy, sortDir)
	}

	// Append limit and offset
	dataArgs = append(dataArgs, limit, offset)

	// Total count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM stocks %s`, whereClause)
	var total int
	err := db.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Paginated query
	dataQuery := fmt.Sprintf(`
		SELECT ticker, company, brokerage, action,
		       rating_from, rating_to, target_from, target_to, time
		FROM stocks
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderClause, len(dataArgs)-1, len(dataArgs))

	rows, err := db.DB.Query(dataQuery, dataArgs...)
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

func GetStocksFiltered(page, limit int, search string) ([]model.Stock, int, error) {
	offset := (page - 1) * limit
	var rowsQuery string
	var totalQuery string
	var args []interface{}

	pattern := "%" + strings.ToLower(search) + "%"

	if search != "" {
		rowsQuery = `
			SELECT ticker, company, brokerage, action,
			       rating_from, rating_to, target_from, target_to, time
			FROM stocks
			WHERE LOWER(ticker) LIKE $1 OR LOWER(company) LIKE $1 OR LOWER(brokerage) LIKE $1
			ORDER BY time DESC
			LIMIT $2 OFFSET $3
		`
		totalQuery = `
			SELECT COUNT(*) FROM stocks
			WHERE LOWER(ticker) LIKE $1 OR LOWER(company) LIKE $1 OR LOWER(brokerage) LIKE $1
		`
		args = []interface{}{pattern, limit, offset}
	} else {
		rowsQuery = `
			SELECT ticker, company, brokerage, action,
			       rating_from, rating_to, target_from, target_to, time
			FROM stocks
			ORDER BY time DESC
			LIMIT $1 OFFSET $2
		`
		totalQuery = `SELECT COUNT(*) FROM stocks`
		args = []interface{}{limit, offset}
	}

	// Obtener total
	var total int
	var countErr error
	if search != "" {
		countErr = db.DB.QueryRow(totalQuery, pattern).Scan(&total)
	} else {
		countErr = db.DB.QueryRow(totalQuery).Scan(&total)
	}
	if countErr != nil {
		return nil, 0, countErr
	}

	rows, err := db.DB.Query(rowsQuery, args...)
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

func GetStocksByDate(date string) ([]model.Stock, error) {
	query := `
		SELECT ticker, company, brokerage, action,
		       rating_from, rating_to, target_from, target_to, time
		FROM stocks
		WHERE time::date = $1
		ORDER BY time DESC
	`

	rows, err := db.DB.Query(query, date)
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

func GetStocksByLocalDateRange(startUTC, endUTC time.Time) ([]model.Stock, error) {
	query := `
		SELECT ticker, company, brokerage, action,
		       rating_from, rating_to, target_from, target_to, time
		FROM stocks
		WHERE time >= $1 AND time < $2
		ORDER BY time DESC
	`
	rows, err := db.DB.Query(query, startUTC, endUTC)
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
