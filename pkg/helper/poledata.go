package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const maxRequestCount = 100

type fetchData struct {
	TotalCount int               `json:"total_count"`
	Results    []json.RawMessage `json:"results"`
}

func GetPoleUrls() []string {
	urls := []string{}
	for _, paal := range []string{"coupure-links", "zuidparklaan", "visserij", "isabellakaai", "groendreef", "bataviabrug"} {
		// date needs to be in 2023, code will stop working in 2025
		date := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
		hour := time.Now().Round(5 * time.Minute).Format("15:04:05")
		hour = strings.TrimLeft(hour, "0")

		splitted := strings.Split(hour, ":")
		var hourformatted string
		for _, s := range splitted {
			hourformatted += s + "%3A"
		}
		hourformatted = hourformatted[:len(hourformatted)-3]

		urls = append(urls, fmt.Sprintf("https://data.stad.gent/api/explore/v2.1/catalog/datasets/fietstelpaal-%s-2023-gent/records?where=datum%20%3D%%20date%%27%s%%27%%20and%%20uur5minuten%%20%%3D%%20%%27%s%%27", paal, date, hourformatted))
	}
	return urls
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
	return allItems
}

func request[T any](url string, offset int) (fetchData, error) {
	resp, err := http.Get(fmt.Sprintf("%s&limit=%d&offset=%d", url, maxRequestCount, offset))
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

	return data, nil
}
