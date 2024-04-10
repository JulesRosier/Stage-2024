package database

import (
	"log/slog"
	"stage2024/pkg/helper"
	"stage2024/pkg/protogen/users"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func createdUser(user User, change helper.Change, db *gorm.DB) error {
	slog.Debug("User created, sending event...", "station", user.Id)
	now := timestamppb.Now()
	protostruct := &users.UserRegistered{
		TimeStamp: now,
		User:      user.IntoId(),
	}

	if err := createOutboxRecord(now, protostruct, db); err != nil {
		return err
	}
	return nil
}
