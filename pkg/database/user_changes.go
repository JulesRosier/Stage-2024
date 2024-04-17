package database

import (
	"log/slog"
	"stage2024/pkg/protogen/users"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func createUserEvent(user *User, eventTime time.Time, db *gorm.DB) error {
	slog.Debug("User created, sending event...", "station", user.Id)

	protostruct := &users.UserRegistered{
		TimeStamp: timestamppb.New(eventTime),
		User:      user.IntoId(),
	}

	if err := createOutboxRecord(timestamppb.New(eventTime), protostruct, db); err != nil {
		return err
	}
	return nil
}
