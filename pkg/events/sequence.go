package events

import (
	"log/slog"
	"math"
	"math/rand"
	"stage2024/pkg/database"
	h "stage2024/pkg/helper"
	"stage2024/pkg/kafka"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/stations"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const warpSpeed = true
const chanceAbandond = 0.01
const chanceDefect = 0.01
const chanceImmobilized = 0.5

func RunSequence(kc *kafka.Client) {
	slog.Info("Starting sequence")
	defer slog.Info("Sequence done")
	db := database.GetDb()

	user := database.User{}
	db.Order("random()").Find(&user)
	userId := user.IntoId()

	// TODO: only get bike with is not in use
	bike := database.Bike{}
	db.Order("random()").Find(&bike)
	bikeId := bike.IntoId()

	station := database.Station{}
	db.Order("random()").Find(&station)

	err := kc.Produce(&bikes.BikeReserved{
		TimeStamp: timestamppb.Now(),
		Bike: &bikes.Bike{
			Bike:       bikeId,
			IsElectric: bike.IsElectric.Bool,
		},
		Station: &stations.StationIdentification{
			Id: station.Id,
			Location: &common.Location{
				Latitude:  bike.Lat,
				Longitude: bike.Lon,
			},
			Name: station.Name,
		},
		User: userId,
	})
	h.MaybeDieErr(err)

	randSleep(60*5, 60)

	err = kc.Produce(&bikes.BikePickedUp{
		TimeStamp: timestamppb.Now(),
		Bike:      bikeId,
		Station: &stations.StationIdentification{
			Id: station.Id,
			Location: &common.Location{
				Latitude:  bike.Lat,
				Longitude: bike.Lon,
			},
			Name: station.Name,
		},
		User: userId,
	})
	h.MaybeDieErr(err)

	randSleep(60*12, 60)

	if rand.Float32() <= chanceAbandond {
		err = kc.Produce(&bikes.BikeAbandoned{
			TimeStamp: timestamppb.Now(),
			Bike: &bikes.AbandonedBike{
				Bike: bikeId,
				Location: &common.Location{
					Latitude:  gofakeit.Latitude(),
					Longitude: gofakeit.Longitude(),
				},
			},
			User: userId,
		})
		h.MaybeDieErr(err)
		return
	}

	if rand.Float32() <= chanceDefect {
		immobilized := false
		if rand.Float32() <= chanceImmobilized {
			immobilized = true
		}
		loc := &common.Location{
			Latitude:  gofakeit.Latitude(),
			Longitude: gofakeit.Longitude(),
		}
		err = kc.Produce(&bikes.BikeDefectReported{
			TimeStamp: timestamppb.Now(),
			Bike: &bikes.DefectBike{
				Bike:          bikeId,
				Location:      loc,
				IsElectric:    bike.IsElectric.Bool,
				IsImmobilized: immobilized,
			},
			User:           userId,
			ReportedDefect: gofakeit.LoremIpsumSentence(50),
		})
		h.MaybeDieErr(err)
		if immobilized {
			err = kc.Produce(&bikes.BikeImmobilized{
				TimeStamp: timestamppb.Now(),
				Bike: &bikes.LocationBike{
					Bike:     bikeId,
					Location: loc,
				},
			})
			h.MaybeDieErr(err)
			return
		}
	}

	dropoffStation := database.Station{}
	db.Order("random()").Find(&dropoffStation)

	randSleep(60*2, 60)

	err = kc.Produce(&bikes.BikeReturned{
		TimeStamp: timestamppb.Now(),
		Bike:      bikeId,
		Station: &stations.StationIdentification{
			Id: dropoffStation.Id,
			Location: &common.Location{
				Latitude:  gofakeit.Latitude(),
				Longitude: gofakeit.Latitude(),
			},
			Name: dropoffStation.Name,
		},
		User: userId,
	})
	h.MaybeDieErr(err)

}

func generateNormalDuration(mean, stdDev float64) float64 {
	// Generate two random numbers from a uniform distribution
	u1 := rand.Float64()
	u2 := rand.Float64()

	// Apply Box-Muller transform to get values from a standard normal distribution
	z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	// Scale and shift the values to get values from a normal distribution with given mean and standard deviation
	duration := mean + z0*stdDev

	return duration
}

func randSleep(mean, stdDev int) {
	t := time.Second*time.Duration(generateNormalDuration(float64(mean), float64(stdDev))) + time.Second*30
	if warpSpeed {
		t = t / 60
	}
	slog.Info("Sleeping", "time", t)
	time.Sleep(t)
}
