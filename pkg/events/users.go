package events

import (
	"stage2024/pkg/database"
	"stage2024/pkg/kafka"
	"stage2024/pkg/protogen/users"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func CreateRandomUser(db *gorm.DB, kc *kafka.Client) {
	u := &database.User{
		Id:           gofakeit.UUID(),
		UserName:     gofakeit.Username(),
		EmailAddress: gofakeit.Email(),
	}
	kc.Produce(&users.UserRegistered{
		TimeStamp: timestamppb.Now(),
		User: &users.UserIdentification{
			Id:           u.Id,
			UserName:     u.UserName,
			EmailAddress: u.EmailAddress,
		},
	})
	db.Create(u)
}
