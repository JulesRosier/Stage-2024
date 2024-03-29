package database

import (
	"fmt"
	"log/slog"
	"os"
	"stage2024/pkg/helper"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func Init() {

	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbDatabase := os.Getenv("DB_DATABASE_OLTP")
	DbHost := os.Getenv("DB_HOST")
	DbPort := os.Getenv("DB_PORT")

	slog.Info("Starting database", "host", DbHost, "database", DbDatabase)

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DbUser, DbPassword, DbDatabase, DbHost, DbPort)

	database, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		helper.MaybeDieErr(err)
	}
	db = database

	slog.Info("Migrating database")
	err = db.AutoMigrate(&User{}, &Bike{}, &Station{})
	helper.MaybeDieErr(err)
}

func GetDb() *gorm.DB {
	return db
}

// updates records in the database and notifies changes through a channel.
func UpdateRecords(db *gorm.DB, channelCh chan []string, records []*Bike) {
	for _, record := range records {
		oldrecord := &Bike{}
		result := db.Limit(1).Find(&oldrecord, "id = ?", record.Id)

		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&record)

		if result.RowsAffected == 0 {
			continue
		}

		changedColumns := helper.ColumnChange(oldrecord, record)

		if len(changedColumns) > 1 {
			channelCh <- changedColumns
		}
	}
}
