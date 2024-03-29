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
func UpdateBike(channel chan []string, records []*Bike) {
	db := GetDb()
	for _, record := range records {
		oldrecord := &Bike{}
		result := db.Limit(1).Find(&oldrecord, "open_data_id = ?", record.OpenDataId)

		// if record does not exist, create it and exit
		if result.RowsAffected == 0 {
			db.Create(&record)
			continue
		}

		record.Id = oldrecord.Id
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&record)

		changedColumns := helper.ColumnChange(oldrecord, record)
		changedColumns = append([]string{"bikes"}, changedColumns...)

		if len(changedColumns) > 2 {
			channel <- changedColumns
		}
	}
}

// updates records in the database and notifies changes through a channel.
func UpdateStation(channel chan []string, records []*Station) {
	db := GetDb()
	for _, record := range records {
		oldrecord := &Station{}
		result := db.Limit(1).Find(&oldrecord, "open_data_id = ?", record.OpenDataId)

		if result.RowsAffected == 0 {
			db.Create(&record)
			continue
		}

		record.Id = oldrecord.Id
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&record)

		changedColumns := helper.ColumnChange(oldrecord, record)
		changedColumns = append([]string{"stations"}, changedColumns...)

		if len(changedColumns) > 2 {
			channel <- changedColumns
		}
	}
}
