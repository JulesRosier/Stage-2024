package database

import (
	"log/slog"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// Sends BikeImmobilized event. Needs bike and eventTime.
func BikeImmobilizedEvent(bike Bike, eventTime time.Time, db *gorm.DB) error {
	slog.Debug("Bike immobilized, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikeImmobilized{
		TimeStamp: timestamppb.New(eventTime),
		Bike: &bikes.LocationBike{
			Bike: bike.IntoId(),
			Location: &common.Location{
				Latitude:  bike.Lat,
				Longitude: bike.Lon,
			},
		},
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}
	return nil
}

// Sends BikeAbandoned event. Needs user
func BikeAbandonedEvent(bike Bike, eventTime time.Time, user User, db *gorm.DB) error {
	slog.Debug("Bike abandoned, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikeAbandoned{
		TimeStamp: timestamppb.New(eventTime),
		Bike: &bikes.AbandonedBike{
			Bike: bike.IntoId(),
			Location: &common.Location{
				Latitude:  bike.Lat,
				Longitude: bike.Lon,
			},
		},
		User: user.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}
	return nil
}

// sends BikeStored events. Needs bike.
func BikeStoredEvent(bike Bike, eventTime time.Time, db *gorm.DB) error {
	slog.Debug("Bike stored, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikeStored{
		TimeStamp: timestamppb.New(eventTime),
		Bike:      bike.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}
	return nil
}

// sends BikeDeployedEvent event
func BikeDeployedEvent(bike Bike, eventTime time.Time, db *gorm.DB) error {
	slog.Debug("Bike deployed, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikeDeployed{
		TimeStamp: timestamppb.New(eventTime),
		Bike:      bike.IntoId(),
	}

	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}
	return nil
}

// Sends BikeReservedEvent event. Needs station and user.
func BikeReservedEvent(bike Bike, eventTime time.Time, station Station, user User, db *gorm.DB) error {
	slog.Debug("Bike reserved, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikeReserved{
		TimeStamp: timestamppb.New(eventTime),
		Bike: &bikes.Bike{
			Bike:       bike.IntoId(),
			IsElectric: bike.IsElectric.Bool,
		},
		Station: station.IntoId(),
		User:    user.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}
	return nil
}

// Sends BikeDefectReported event. Needs bike, user and defect.
func BikeDefectEvent(bike Bike, eventTime time.Time, user User, defect string, db *gorm.DB) error {
	slog.Debug("Bike defect, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikeDefectReported{
		TimeStamp: timestamppb.New(eventTime),
		Bike: &bikes.DefectBike{
			Bike: bike.IntoId(),
			Location: &common.Location{
				Latitude:  bike.Lat,
				Longitude: bike.Lon,
			},
			IsElectric: bike.IsElectric.Bool,
		},
		User:           user.IntoId(),
		ReportedDefect: defect,
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}
	return nil
}

// Sends BikePickedUp event. Needs bike, station, user.
func BikePickedUpEvent(bike Bike, eventTime time.Time, station Station, user User, db *gorm.DB) error {
	slog.Debug("Bike picked up, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikePickedUp{
		TimeStamp: timestamppb.New(eventTime),
		Bike:      bike.IntoId(),
		Station:   station.IntoId(),
		User:      user.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}

	return nil
}

func BikeReturnedEvent(bike Bike, eventTime time.Time, station Station, user User, db *gorm.DB) error {
	slog.Debug("Bike returned, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikeReturned{
		TimeStamp: timestamppb.New(eventTime),
		Bike:      bike.IntoId(),
		Station:   station.IntoId(),
		User:      user.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}

	return nil
}

func BikeCreatedEvent(bike *Bike, db *gorm.DB, eventTime time.Time) error {
	slog.Debug("Bike created, sending event...", "bike", bike.ID)
	protostruct := &bikes.BikeDeployed{
		TimeStamp: timestamppb.New(eventTime),
		Bike:      bike.IntoId(),
	}
	if err := createOutboxRecord(timestamppb.New(eventTime), protostruct, db); err != nil {
		return err
	}
	return nil
}

// Sends BikeRepaired event. Needs bike and eventTime.
func BikeRepairedEvent(bike Bike, eventTime time.Time, db *gorm.DB) error {
	slog.Debug("Bike repaired, sending event...", "bike", bike.Id)
	protostruct := &bikes.BikeRepaired{
		TimeStamp: timestamppb.New(eventTime),
		Bike:      bike.IntoId(),
	}
	if err := createOutboxRecord(protostruct.TimeStamp, protostruct, db); err != nil {
		return err
	}

	return nil
}
