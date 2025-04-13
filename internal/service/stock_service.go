// internal/service/stock_service.go
package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"appstock/internal/model"
	"appstock/internal/repository"
)

// const apiURL = "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list"
var apiURL = os.Getenv("STOCK_API_URL")
var apiKey = os.Getenv("STOCK_API_KEY")

var bearerToken = fmt.Sprintf("Bearer %s", apiKey)

type ApiResponse struct {
	Items    []model.Stock `json:"items"`
	NextPage string        `json:"next_page"`
}

func FetchAndStoreAllStocks() error {
	nextPage := ""
	for {
		url := apiURL
		if nextPage != "" {
			url += "?next_page=" + nextPage
		}

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", bearerToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var data ApiResponse
		if err := json.Unmarshal(body, &data); err != nil {
			return err
		}

		for _, stock := range data.Items {
			// Limpia $ en los valores target
			stock.TargetFrom = parseDollar(stock.TargetFromRaw)
			stock.TargetTo = parseDollar(stock.TargetToRaw)
			repository.SaveStock(stock)
		}

		if data.NextPage == "" {
			break
		}
		nextPage = data.NextPage
	}
	return nil
}

func parseDollar(s string) float64 {
	s = strings.ReplaceAll(s, "$", "")
	var val float64
	fmt.Sscanf(s, "%f", &val)
	return val
}
