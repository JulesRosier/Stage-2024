package gentopendata

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const maxRequestCount = 100

type fetchData[T any] struct {
	TotalCount int `json:"total_count"`
	Results    []T `json:"results"`
}

func Fetch[T any](url string) []T {
	offset := 0
	totalCount := 1

	allItems := make([]T, 0)

	for offset < totalCount {
		data, err := request[T](url, offset)
		if err != nil {
			fmt.Println("Error fetching data:", err)
			break
		}
		totalCount = data.TotalCount
		allItems = append(allItems, data.Results...)
		offset += len(data.Results)
	}

	slog.Info("Total expented items", "count", totalCount)
	slog.Info("Total items fetched", "count", len(allItems))
	return allItems
}

func request[T any](url string, offset int) (fetchData[T], error) {
	slog.Info("Makking request", "offset", offset)
	resp, err := http.Get(fmt.Sprintf("%s?limit=%d&offset=%d", url, maxRequestCount, offset))
	var data fetchData[T]
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}
	slog.Info("Request result", "item_count", len(data.Results))

	return data, nil
}
