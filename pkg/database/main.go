package database

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"stage2024/pkg/helper"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DatabaseClient struct {
	DB *gorm.DB
}

func NewDatabase() *DatabaseClient {

	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbDatabase := os.Getenv("DB_DATABASE_OLTP")
	DbHost := os.Getenv("DB_HOST")
	DbPort := os.Getenv("DB_PORT")

	slog.Info("Starting database", "host", DbHost, "database", DbDatabase)

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DbUser, DbPassword, DbDatabase, DbHost, DbPort)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		helper.MaybeDieErr(err)
	}

	slog.Info("Migrating database")
	err = db.AutoMigrate(&User{}, &Bike{}, &Station{}, &Outbox{}, &HistoricalStationData{})
	helper.MaybeDieErr(err)

	return &DatabaseClient{
		DB: db,
	}
}

// updates records in the database and notifies changes through a channel.
func UpdateBike(records []*Bike, db *gorm.DB) {
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

		ColumnChange(oldrecord, record, db)
	}
}

// Updates records in the database and sends changed records to function that checks changes
func UpdateStation(records []*Station, db *gorm.DB) {

	for _, record := range records {
		oldrecord := &Station{}
		result := db.Limit(1).Find(&oldrecord, "open_data_id = ?", record.OpenDataId)

		//start transaction
		err := db.Transaction(func(tx *gorm.DB) error {
			if result.RowsAffected == 0 {

				if err := db.Create(&record).Error; err != nil {
					return err
				}

				// send change for created record to function
				created := helper.Change{
					Table:      "Station",
					Column:     "Created",
					Id:         record.Id,
					OpenDataId: record.OpenDataId,
				}

				err := ChangeDetected(created, tx)
				if err != nil {
					return err
				}

				return nil

			} else {
				record.Id = oldrecord.Id

				err := tx.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&record).Error
				if err != nil {
					return err
				}

				if err := ColumnChange(oldrecord, record, tx); err != nil {
					return err
				}

				return nil
			}
		})

		if err != nil {
			slog.Warn("Transaction failed", "error", err)
		}
	}
}

func GetStationById(id string, db *gorm.DB) (Station, error) {
	if id == "" {
		return Station{}, fmt.Errorf("StationId is empty")
	}
	var station Station
	db.Where("id = ?", id).First(&station)
	return station, nil
}

func GetBikeById(id string, db *gorm.DB) (Bike, error) {
	if id == "" {
		return Bike{}, fmt.Errorf("BikeId is empty")
	}
	var bike Bike
	db.Where("id = ?", id).First(&bike)
	return bike, nil
}

func GetUserById(id string, db *gorm.DB) (User, error) {
	if id == "" {
		return User{}, fmt.Errorf("UserId is empty")
	}
	var user User
	db.Where("id = ?", id).First(&user)
	return user, nil
}

func createOutboxRecord(now *timestamppb.Timestamp, protostruct proto.Message, db *gorm.DB) error {
	payload, err := proto.Marshal(protostruct)
	if err != nil {
		return err
	}
	topic := helper.ToSnakeCase(reflect.TypeOf(protostruct).Elem().Name())
	return db.Create(&Outbox{EventTimestamp: now.AsTime(), Topic: topic, Payload: payload}).Error
}

// adds historical station data
func addHistoricaldata(record *Station, topicname string, db *gorm.DB) error {
	historicaldata := record.ToHistoricalStationData()
	historicaldata.TopicName = topicname
	if err := db.Create(&historicaldata).Error; err != nil {
		return err
	}
	return nil
}
