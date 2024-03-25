package main

import "stage2024/pkg/database"

func main() {
	database.Init()

	db := database.GetDb()
	db.Create(&database.Test{Code: "D42", Price: 100})
}
