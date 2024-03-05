package gentopendata

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const maxRequestCount = 100

type fetchData struct {
	TotalCount int               `json:"total_count"`
	Results    []json.RawMessage `json:"results"`
}

func Fetch[T any](url string, f func([]byte) T) []T {
	offset := 0
	totalCount := 1

	allItems := make([]T, 0)

	for offset < totalCount {
		data, err := request[T](url, offset)
		for _, x := range data.Results {
			allItems = append(allItems, f(x))
		}
		if err != nil {
			slog.Error("Error fetching data:", "error", err)
			break
		}
		totalCount = data.TotalCount
		offset += len(data.Results)
	}

	slog.Info("Total expented items", "count", totalCount)
	slog.Info("Total items fetched", "count", len(allItems))
	return allItems
}

func request[T any](url string, offset int) (fetchData, error) {
	slog.Info("Makking request", "offset", offset)
	resp, err := http.Get(fmt.Sprintf("%s?limit=%d&offset=%d", url, maxRequestCount, offset))
	var data fetchData
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
	// fmt.Println(data)
	slog.Info("Request result", "item_count", len(data.Results))

	return data, nil
}
