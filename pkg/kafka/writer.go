package kafka

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/bikes"
	"stage2024/pkg/protogen/stations"
	"stage2024/pkg/protogen/users"

	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const protoDir = "./proto"

type topic struct {
	name      string
	protoFile protoreflect.FileDescriptor
	pType     any
}

func GetSerde(rcl *sr.Client) *sr.Serde {
	topics := []topic{
		{protoFile: users.File_users_user_registered_proto, pType: &users.UserRegistered{}},

		{protoFile: bikes.File_bikes_bike_abandoned_proto, pType: &bikes.AbandonedBike{}},
		{protoFile: bikes.File_bikes_bike_defetect_reported_proto, pType: &bikes.BikeDefectReported{}},
		{protoFile: bikes.File_bikes_bike_brought_out_proto, pType: &bikes.BikeImmobilized{}},
		{protoFile: bikes.File_bikes_bike_immobilized_proto, pType: &bikes.BikeImmobilized{}},
		{protoFile: bikes.File_bikes_bike_picked_up_proto, pType: &bikes.BikePickedUp{}},
		{protoFile: bikes.File_bikes_bike_reserved_proto, pType: &bikes.BikeReserved{}},
		{protoFile: bikes.File_bikes_bike_returned_proto, pType: &bikes.BikeReturned{}},
		{protoFile: bikes.File_bikes_bike_stored_proto, pType: &bikes.BikeStored{}},

		{protoFile: stations.File_stations_station_capacity_decreased_proto, pType: &stations.StationCapacityDecreased{}},
		{protoFile: stations.File_stations_station_capacity_exhausted_proto, pType: &stations.StationCapacityExhausted{}},
		{protoFile: stations.File_stations_station_capacity_increased_proto, pType: &stations.StationCapacityIncreased{}},
		{protoFile: stations.File_stations_station_created_proto, pType: &stations.StationCreated{}},
		{protoFile: stations.File_stations_station_deprecated_proto, pType: &stations.StationDeprecated{}},
	}

	serde := &sr.Serde{}
	serde.SetDefaults()

	for _, topic := range topics {

		file, err := os.ReadFile(filepath.Join("./proto", topic.protoFile.Path()))
		h.MaybeDieErr(err)

		var subject string
		if topic.name == "" {
			base := filepath.Base(topic.protoFile.Path())
			fileName := base[:len(base)-len(filepath.Ext(base))]
			subject = fileName + "-value"
		} else {
			subject = topic.name
		}

		refs := getReferences(rcl, topic.protoFile)

		ss, err := rcl.CreateSchema(context.TODO(), subject, sr.Schema{
			Schema:     string(file),
			Type:       sr.TypeProtobuf,
			References: refs,
		})
		h.MaybeDie(err, "Failed to create schema")
		slog.Debug("Created or reusing schema", "subject", subject)

		serde.Register(
			ss.ID,
			topic.pType,
			sr.EncodeFn(func(a any) ([]byte, error) {
				return proto.Marshal(a.(proto.Message))
			}),
			sr.Index(0),
			sr.DecodeFn(func(b []byte, a any) error {
				return proto.Unmarshal(b, a.(proto.Message))
			}),
		)
	}
	return serde
}

func getReferences(rcl *sr.Client, protoFile protoreflect.FileDescriptor) []sr.SchemaReference {
	// FIXME: cache referances
	l := protoFile.Imports().Len()
	i := 0
	refs := []sr.SchemaReference{}
	for i < l {
		f := protoFile.Imports().Get(i)
		i++
		sr, err := createReference(rcl, f)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			h.MaybeDieErr(err)
		}
		refs = append(refs, sr)
	}
	return refs
}

func createReference(rcl *sr.Client, protoPath protoreflect.FileDescriptor) (sr.SchemaReference, error) {
	file, err := os.ReadFile(filepath.Join(protoDir, protoPath.Path()))
	if err != nil {
		return sr.SchemaReference{}, err
	}

	ssLocation, err := rcl.CreateSchema(context.Background(), protoPath.Path(),
		sr.Schema{
			Schema:     string(file),
			Type:       sr.TypeProtobuf,
			References: getReferences(rcl, protoPath),
		},
	)
	if err != nil {
		return sr.SchemaReference{}, err
	}
	slog.Debug("Created or reusing schema",
		"subject", ssLocation.Subject,
		"version", ssLocation.Version,
		"id", ssLocation.ID,
	)

	return sr.SchemaReference{
		Name:    protoPath.Path(),
		Subject: ssLocation.Subject,
		Version: ssLocation.Version,
	}, nil

}
