package main

import (
	"stage2024/pkg/database"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
)

const maxUser = 100

func main() {
	database.Init()

	db := database.GetDb()
	db.Create(&database.Test{Code: "D42", Price: 100})

	CreateUsers(db)
}

func CreateUsers(db *gorm.DB) {
	var userCount int64
	db.Model(&database.User{}).Count(&userCount)

	for userCount < maxUser {
		CreateRandomUser(db)
		userCount++
	}
}

func CreateRandomUser(db *gorm.DB) {
	u := &database.User{
		ID:          gofakeit.UUID(),
		UserName:    gofakeit.Username(),
		EmailAdress: gofakeit.Email(),
	}
	db.Create(u)
}
