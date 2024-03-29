package kafka

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	h "stage2024/pkg/helper"

	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const protoDir = "./proto"

type Topic struct {
	Name      string
	ProtoFile protoreflect.FileDescriptor
	PType     any
}

func (t Topic) getName(suffix string) string {
	if t.Name == "" {
		base := filepath.Base(t.ProtoFile.Path())
		fileName := base[:len(base)-len(filepath.Ext(base))]
		return fileName + suffix
	} else {
		return t.Name
	}

}

func getSerde(rcl *sr.Client, topics []Topic) *sr.Serde {

	serde := &sr.Serde{}
	serde.SetDefaults()

	for _, topic := range topics {

		file, err := os.ReadFile(filepath.Join("./proto", topic.ProtoFile.Path()))
		h.MaybeDieErr(err)

		subject := topic.getName("-value")

		refs := getReferences(rcl, topic.ProtoFile)

		ss, err := rcl.CreateSchema(context.TODO(), subject, sr.Schema{
			Schema:     string(file),
			Type:       sr.TypeProtobuf,
			References: refs,
		})
		h.MaybeDie(err, "Failed to create schema")
		slog.Debug("Created or reusing schema", "subject", subject)

		serde.Register(
			ss.ID,
			topic.PType,
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
