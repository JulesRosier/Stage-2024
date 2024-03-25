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
		Id:           gofakeit.UUID(),
		UserName:     gofakeit.Username(),
		EmailAddress: gofakeit.Email(),
	}
	db.Create(u)
}
