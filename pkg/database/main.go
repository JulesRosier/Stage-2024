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
	"gorm.io/gorm/logger"
)

type DatabaseClient struct {
	DB *gorm.DB
}

func NewDatabase() *DatabaseClient {

	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbDatabaseOltp := os.Getenv("DB_DATABASE_OLTP")
	DbDatabase := os.Getenv("DB_DATABASE")
	DbHost := os.Getenv("DB_HOST")
	DbPort := os.Getenv("DB_PORT")

	slog.Info("Starting database", "host", DbHost, "database", DbDatabase)

	// connect to db and create database if it does not exist
	createconnstr := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable",
		DbUser, DbPassword, DbHost, DbPort)
	dbcreate, err := gorm.Open(postgres.Open(createconnstr), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		helper.MaybeDieErr(err)
	}

	if err := createDb(dbcreate, DbDatabase); err != nil {
		helper.MaybeDieErr(err)
	}
	if err := createDb(dbcreate, DbDatabaseOltp); err != nil {
		helper.MaybeDieErr(err)
	}

	// connect to database
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		DbUser, DbPassword, DbDatabaseOltp, DbHost, DbPort)

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

// Creates a new database
func createDb(db *gorm.DB, name string) error {
	err := db.Exec("CREATE DATABASE " + name + ";").Error
	if err != nil {
		if err.Error() == fmt.Sprintf("ERROR: database \"%s\" already exists (SQLSTATE 42P04)", name) {
			slog.Info("Database already exists", "database", name)
			return nil
		}
		return err
	}
	slog.Info("Database created", "database", name)
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
