package helper

import (
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/events"
	"stage2024/pkg/protogen/occupations"
	"stage2024/pkg/protogen/poles"
)

func GetMessageProto(topic string) any {
	var messageProto any
	switch topic {
	case "baqme-locations":
		messageProto = &bikes.BaqmeLocation{}
	case "bluebike-locations":
		messageProto = &occupations.BlueBikeOccupation{}
	case "bolt-locations":
		messageProto = &bikes.BoltLocation{}
	case "countpoles":
		messageProto = &poles.PoleData{}
	case "donkey-locations":
		messageProto = &occupations.DonkeyLocation{}
	case "stalling-gent":
		messageProto = &occupations.StallingGent{}
	case "stalling-stadskantoor":
		messageProto = &occupations.StallingStadskantoor{}
	case "bike.droppedoff":
		messageProto = &events.BikeDropOff{}
	}
	return messageProto
}
