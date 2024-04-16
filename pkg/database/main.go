package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"reflect"
	"stage2024/pkg/helper"
	"time"

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

// Updates an existing Bike record in the database
func UpdateBike(db *gorm.DB, record *Bike) error {
	slog.Debug("Updating bike")
	err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&record).Error
	if err != nil {
		return err
	}
	return nil
}

// Updates records in the database and sends changed records to function that checks changes
func UpdateStation(records []*Station, db *gorm.DB) {
	slog.Debug("Updating stations")
	for _, record := range records {
		oldrecord := &Station{}
		result := db.Limit(1).Find(&oldrecord, "open_data_id = ?", record.OpenDataId)

		//start transaction
		err := db.Transaction(func(tx *gorm.DB) error {
			if result.RowsAffected == 0 {
				if err := db.Create(&record).Error; err != nil {
					return err
				}

				// send change for created record to OUtbox
				if err := createStationEvent(record, tx); err != nil {
					return err
				}
			} else {
				record.Id = oldrecord.Id

				err := tx.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&record).Error
				if err != nil {
					return err
				}

				if err := ColumnChange(oldrecord, record, tx, helper.Change{}); err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			slog.Warn("Transaction failed", "error", err)
		}
	}
}

func UpdateUser(db *gorm.DB, record *User) error {
	slog.Debug("Updating user")
	err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&record).Error
	if err != nil {
		return err
	}
	return nil
}

func GetStationById(id string, db *gorm.DB) (Station, error) {
	if id == "" {
		return Station{}, fmt.Errorf("StationId is empty")
	}
	var station Station
	err := db.Where("id = ?", id).First(&station).Error
	return station, err
}

func GetBikeById(id string, db *gorm.DB) (Bike, error) {
	if id == "" {
		return Bike{}, fmt.Errorf("BikeId is empty")
	}
	var bike Bike
	err := db.Where("id = ?", id).First(&bike).Error
	return bike, err
}

func GetUserById(id string, db *gorm.DB) (User, error) {
	if id == "" {
		return User{}, fmt.Errorf("UserId is empty")
	}
	var user User
	err := db.Where("id = ?", id).First(&user).Error
	return user, err
}

func createOutboxRecord(now *timestamppb.Timestamp, protostruct proto.Message, db *gorm.DB) error {
	slog.Debug("Creating outbox record")
	payload, err := proto.Marshal(protostruct)
	if err != nil {
		return err
	}
	topic := helper.ToSnakeCase(reflect.TypeOf(protostruct).Elem().Name())
	return db.Create(&Outbox{EventTimestamp: now.AsTime(), Topic: topic, Payload: payload}).Error
}

// adds historical station data
func addHistoricaldata(record *Station, topicname string, db *gorm.DB, amountChanged int32, eventTimeStamp *timestamppb.Timestamp) error {
	slog.Debug("Adding historical data")
	historicaldata := record.ToHistoricalStationData()
	historicaldata.TopicName = topicname
	historicaldata.AmountChanged = amountChanged
	historicaldata.AmountFaked = 0
	historicaldata.EventTimeStamp = eventTimeStamp.AsTime()
	err := db.Create(&historicaldata).Error

	return err
}

// makes unavailable bikes available again
func BikeCleanUp(db *gorm.DB) error {
	slog.Debug("Cleaning up bikes")
	bikes := []Bike{}
	db.Where("(is_defect = true OR is_immobilized = true OR is_abandoned = true OR is_in_storage = true OR is_returned = false) AND in_use_timestamp IS NOT NULL AND extract(epoch from ? - in_use_timestamp)/60 > 0", time.Now().UTC().Format("2006-01-02 15:04:05.999999-07")).Find(&bikes)
	for _, bike := range bikes {

		eventTime1, eventTime2, eventTime3 := getEventTimes(bike)

		change := helper.Change{}
		bike.IsReserved = sql.NullBool{Bool: false, Valid: true}
		bike.PickedUp = sql.NullBool{Bool: false, Valid: true}
		bike.IsDefect = sql.NullBool{Bool: false, Valid: true}
		bike.IsImmobilized = sql.NullBool{Bool: false, Valid: true}
		bike.IsAbandoned = sql.NullBool{Bool: false, Valid: true}
		if !bike.IsInStorage.Bool {
			change.EventTime = eventTime1
			if err := BikeInStorageEvent(bike, change, db); err != nil {
				slog.Warn("Error sending event", "error", err)
			}
		}
		bike.IsInStorage = sql.NullBool{Bool: false, Valid: true}
		if err := BikeRepairedEvent(bike, eventTime2, db); err != nil {
			slog.Warn("Error sending event", "error", err)
		}
		if err := BikeDeployedEvent(bike, eventTime3, db); err != nil {
			return fmt.Errorf("error sending event, error : %v", err)
		}

		bike.IsReturned = sql.NullBool{Bool: true, Valid: true}

		db.Save(&bike)
	}
	return nil
}

// Get event times for bike cleanup function
func getEventTimes(bike Bike) (time.Time, time.Time, time.Time) {
	times := []time.Time{}
	previous := bike.InUseTimestamp.Time
	now := time.Now().UTC()
	for range 3 {
		delta := now.Sub(previous)
		r := rand.Float64()
		if r == 0 {
			r = 0.5
		}
		offset := time.Second * time.Duration(r*delta.Seconds())
		eventTime := previous.Add(offset)
		times = append(times, eventTime)
		previous = eventTime
		slog.Debug("Event time", "time", eventTime)
	}

	return times[0], times[1], times[2]
}
