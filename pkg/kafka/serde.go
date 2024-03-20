package kafka

import (
	"log/slog"
	"os"
	"path/filepath"
	writebaqme "stage2024/cmd/write_baqme"
	writebluebike "stage2024/cmd/write_bluebike"
	writebolt "stage2024/cmd/write_bolt"
	writecountpoles "stage2024/cmd/write_countpoles"
	writedonkey "stage2024/cmd/write_donkey"
	writestallinggent "stage2024/cmd/write_stalling_gent"
	writestallingstadskantoor "stage2024/cmd/write_stalling_stadskantoor"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/occupations"
	"stage2024/pkg/protogen/poles"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/twmb/franz-go/pkg/sr"

	"github.com/redpanda-data/console/backend/pkg/config"
	"github.com/redpanda-data/console/backend/pkg/msgpack"
	"github.com/redpanda-data/console/backend/pkg/proto"
	"github.com/redpanda-data/console/backend/pkg/schema"
	"github.com/redpanda-data/console/backend/pkg/serde"
	"go.uber.org/zap"
	protov2 "google.golang.org/protobuf/proto"
)

// Creates  serde.Service and all other services that are need to start it.
func CreateSerde() *serde.Service {
	registry := "http://" + os.Getenv("REGISTRY")

	slog.Info("Creating serde service", "registry", registry)

	logger := zap.L()

	schemaService, err := schema.NewService(
		config.Schema{
			Enabled: true,
			URLs:    []string{registry},
		},
		logger,
	)
	h.MaybeDie(err, "Failed to create schema service")

	protoService, err := proto.NewService(
		config.Proto{
			Enabled: true,
			SchemaRegistry: config.ProtoSchemaRegistry{
				Enabled:         true,
				RefreshInterval: time.Minute * 1,
			},
		},
		logger,
		schemaService,
	)
	h.MaybeDie(err, "Failed to create proto service")

	err = protoService.Start()
	h.MaybeDie(err, "Failed to start proto service")

	seserv := serde.NewService(
		schemaService,
		protoService,
		&msgpack.Service{},
	)

	return seserv
}

// Returns serde with all the schemas registered
func GetSerde(rcl *sr.Client) *sr.Serde {
	topics := []string{writebaqme.Topic, writebluebike.Topic, writebolt.Topic, writecountpoles.Topic, writedonkey.Topic, writestallinggent.Topic, writestallingstadskantoor.Topic}
	filepaths := []string{
		bikes.File_bikes_baqme_proto.Path(),
		occupations.File_occupations_blue_bike_proto.Path(),
		bikes.File_bikes_bolt_proto.Path(),
		poles.File_poles_pole_proto.Path(),
		occupations.File_occupations_donkey_proto.Path(),
		occupations.File_occupations_stallinggent_proto.Path(),
		occupations.File_occupations_stallingstadskantoor_proto.Path(),
	}
	serde := &sr.Serde{}

	for i, topic := range topics {
		protofilepath := filepaths[i]

		file, err := os.ReadFile(filepath.Join("./proto", protofilepath))
		h.MaybeDieErr(err)

		// schema ophalen
		ss := GetSchema(topic, rcl, file)

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
		}

		serde = registerSerde(serde, ss, messageProto)

	}

	return serde
}

// registerSerde registers a serde with a schema
func registerSerde(serde *sr.Serde, ss sr.SubjectSchema, messageProto interface{}) *sr.Serde {
	serde.Register(
		ss.ID,
		messageProto,
		sr.EncodeFn(func(a any) ([]byte, error) {
			return protov2.Marshal(a.(protov2.Message))
		}),
		sr.Index(0),
		sr.DecodeFn(func(b []byte, a any) error {
			return protov2.Unmarshal(b, a.(protov2.Message))
		}),
	)

	return serde
}
