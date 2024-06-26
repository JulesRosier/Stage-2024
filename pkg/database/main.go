package database

import (
	"fmt"
	"log/slog"
	"reflect"
	"stage2024/pkg/helper"
	"stage2024/pkg/settings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseClient struct {
	DB *gorm.DB
}

func NewDatabase(set settings.Database) *DatabaseClient {
	DbUser := set.User
	DbPassword := set.Paddword
	DbDatabase := set.Database
	DbHost := set.Host
	DbPort := set.Port

	// connect to db and create database if it does not exist
	createConnStr := fmt.Sprintf("user=%s password=%s host=%s port=%d database=postgres sslmode=disable",
		DbUser, DbPassword, DbHost, DbPort)
	dbCreate, err := gorm.Open(postgres.Open(createConnStr), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	helper.MaybeDieErr(err)

	var dbName string
	dbCreate.Raw("SELECT datname FROM pg_catalog.pg_database WHERE datname=?", DbDatabase).Scan(&dbName)
	if dbName != DbDatabase {
		err := dbCreate.Exec("CREATE DATABASE " + DbDatabase + ";").Error
		helper.MaybeDie(err, "Failed to create Database")
	}
	slog.Debug("Created or reusing database", "name", dbName)

	// connect to database
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		DbUser, DbPassword, DbDatabase, DbHost, DbPort)
	slog.Debug("Connecting to database", "connstr", connStr)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	helper.MaybeDieErr(err)

	slog.Debug("Migrating database")
	err = db.AutoMigrate(&User{}, &Bike{}, &Station{}, &Outbox{}, &HistoricalStationData{})
	helper.MaybeDieErr(err)

	slog.Info("Connected to database", "host", DbHost, "database", DbDatabase)
	return &DatabaseClient{
		DB: db,
	}
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
func addHistoricalData(record *Station, topicName string, db *gorm.DB, amountChanged int32, eventTimeStamp *timestamppb.Timestamp) error {
	slog.Debug("Adding historical data")
	historicalData := HistoricalStationData{}
	historicalData.Uuid = record.Id
	historicalData.OpenDataId = record.OpenDataId
	historicalData.TopicName = topicName
	historicalData.AmountChanged = amountChanged
	historicalData.AmountFaked = 0
	historicalData.EventTimeStamp = eventTimeStamp.AsTime()
	err := db.Create(&historicalData).Error

	return err
}
