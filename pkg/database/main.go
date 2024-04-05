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
func UpdateBike(channel chan helper.Change, records []*Bike) {
	db := GetDb()
	for _, record := range records {
		oldrecord := &Bike{}
		result := db.Limit(1).Find(&oldrecord, "open_data_id = ?", record.OpenDataId)

		// if record does not exist, create it and exit
		if result.RowsAffected == 0 {
			db.Create(&record)
			continue
		}

		// if record exists, change uuid to old uuid and update record
		// also leave fake data the same
		record.Id = oldrecord.Id
		record.IsElectric = oldrecord.IsElectric
		record.IsImmobilized = oldrecord.IsImmobilized
		record.IsAbandoned = oldrecord.IsAbandoned
		record.IsInStorage = oldrecord.IsInStorage
		record.IsDefect = oldrecord.IsDefect

		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&record)

		ColumnChange(oldrecord, record, channel)
	}
}

// Updates records in the database and notifies changes through a channel.
func UpdateStation(channel chan helper.Change, records []*Station) {
	db := GetDb()

	for _, record := range records {
		oldrecord := &Station{}
		result := db.Limit(1).Find(&oldrecord, "open_data_id = ?", record.OpenDataId)

		if result.RowsAffected == 0 {
			db.Create(&record)
			// send change for created record to channel
			channel <- helper.Change{
				Table:      "Station",
				Column:     "Created",
				Id:         record.Id,
				OpenDataId: record.OpenDataId,
			}
			continue
		}

		record.Id = oldrecord.Id
		record.IsActive = oldrecord.IsActive
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&record)

		ColumnChange(oldrecord, record, channel)
	}
}

func GetStationById(id string) (Station, error) {
	if id == "" {
		return Station{}, fmt.Errorf("StationId is empty")
	}
	var station Station
	db.Where("id = ?", id).First(&station)
	return station, nil
}

func GetBikeById(id string) (Bike, error) {
	if id == "" {
		return Bike{}, fmt.Errorf("BikeId is empty")
	}
	var bike Bike
	db.Where("id = ?", id).First(&bike)
	return bike, nil
}

func GetUserById(id string) (User, error) {
	if id == "" {
		return User{}, fmt.Errorf("UserId is empty")
	}
	var user User
	db.Where("id = ?", id).First(&user)
	return user, nil
}

// Updates the bike in the database, does NOT notify changes through a channel.
func UpdateBikeNoNotify(bike *Bike) {
	db = GetDb()
	oldbike := &Bike{}
	result := db.Limit(1).Find(&oldbike, "open_data_id = ?", bike.OpenDataId)

	if result.RowsAffected == 0 {
		helper.MaybeDieErr(fmt.Errorf("bike not found %v", bike))
	}

	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&bike)
}
