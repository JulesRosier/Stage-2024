package helper

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"stage2024/pkg/protogen/common"

	"github.com/twmb/franz-go/pkg/sr"
)

const protoDir = "./proto"

func ReferenceLocation(rcl *sr.Client) sr.SchemaReference {
	p := common.File_common_location_proto.Path()
	file, err := os.ReadFile(filepath.Join(protoDir, p))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	ssLocation, err := rcl.CreateSchema(context.Background(), p,
		sr.Schema{
			Schema: string(file),
			Type:   sr.TypeProtobuf,
		},
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("created or reusing schema",
		"subject", ssLocation.Subject,
		"version", ssLocation.Version,
		"id", ssLocation.ID,
	)

	return sr.SchemaReference{
		Name:    p,
		Subject: ssLocation.Subject,
		Version: ssLocation.Version,
	}

}
