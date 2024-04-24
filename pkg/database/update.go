package database

import (
	"log/slog"
	"stage2024/pkg/helper"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Updates an existing Bike record in the database
func UpdateBike(db *gorm.DB, record *Bike) error {
	slog.Debug("Updating bike", "bike", record.Id)
	err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&record).Error
	if err != nil {
		return err
	}
	return nil
}

// Updates an existing Station record in the database and sends events
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

// Updates an existing User record in the database
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
