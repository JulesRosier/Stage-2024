// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: poles/pole.proto

package poles

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PoleData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code          string                 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Locatie       string                 `protobuf:"bytes,2,opt,name=locatie,proto3" json:"locatie,omitempty"`
	Datum         string                 `protobuf:"bytes,3,opt,name=datum,proto3" json:"datum,omitempty"`
	Uur5Minuten   string                 `protobuf:"bytes,4,opt,name=uur5minuten,proto3" json:"uur5minuten,omitempty"`
	Ordening      *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=ordening,proto3" json:"ordening,omitempty"`
	Totaal        int32                  `protobuf:"varint,6,opt,name=totaal,proto3" json:"totaal,omitempty"`
	Tegenrichting int32                  `protobuf:"varint,7,opt,name=tegenrichting,proto3" json:"tegenrichting,omitempty"`
	Hoofdrichting int32                  `protobuf:"varint,8,opt,name=hoofdrichting,proto3" json:"hoofdrichting,omitempty"`
}

func (x *PoleData) Reset() {
	*x = PoleData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_poles_pole_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PoleData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PoleData) ProtoMessage() {}

func (x *PoleData) ProtoReflect() protoreflect.Message {
	mi := &file_poles_pole_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PoleData.ProtoReflect.Descriptor instead.
func (*PoleData) Descriptor() ([]byte, []int) {
	return file_poles_pole_proto_rawDescGZIP(), []int{0}
}

func (x *PoleData) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *PoleData) GetLocatie() string {
	if x != nil {
		return x.Locatie
	}
	return ""
}

func (x *PoleData) GetDatum() string {
	if x != nil {
		return x.Datum
	}
	return ""
}

func (x *PoleData) GetUur5Minuten() string {
	if x != nil {
		return x.Uur5Minuten
	}
	return ""
}

func (x *PoleData) GetOrdening() *timestamppb.Timestamp {
	if x != nil {
		return x.Ordening
	}
	return nil
}

func (x *PoleData) GetTotaal() int32 {
	if x != nil {
		return x.Totaal
	}
	return 0
}

func (x *PoleData) GetTegenrichting() int32 {
	if x != nil {
		return x.Tegenrichting
	}
	return 0
}

func (x *PoleData) GetHoofdrichting() int32 {
	if x != nil {
		return x.Hoofdrichting
	}
	return 0
}

var File_poles_pole_proto protoreflect.FileDescriptor

var file_poles_pole_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x6f, 0x6c, 0x65, 0x73, 0x2f, 0x70, 0x6f, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x70, 0x6f, 0x6c, 0x65, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8c, 0x02, 0x0a, 0x08, 0x50,
	0x6f, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6c,
	0x6f, 0x63, 0x61, 0x74, 0x69, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x61, 0x74, 0x75, 0x6d, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x61, 0x74, 0x75, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x75,
	0x75, 0x72, 0x35, 0x6d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x75, 0x75, 0x72, 0x35, 0x6d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x6e, 0x12, 0x36, 0x0a,
	0x08, 0x6f, 0x72, 0x64, 0x65, 0x6e, 0x69, 0x6e, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x6f, 0x72, 0x64,
	0x65, 0x6e, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x6f, 0x74, 0x61, 0x61, 0x6c, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x74, 0x6f, 0x74, 0x61, 0x61, 0x6c, 0x12, 0x24, 0x0a,
	0x0d, 0x74, 0x65, 0x67, 0x65, 0x6e, 0x72, 0x69, 0x63, 0x68, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x74, 0x65, 0x67, 0x65, 0x6e, 0x72, 0x69, 0x63, 0x68, 0x74,
	0x69, 0x6e, 0x67, 0x12, 0x24, 0x0a, 0x0d, 0x68, 0x6f, 0x6f, 0x66, 0x64, 0x72, 0x69, 0x63, 0x68,
	0x74, 0x69, 0x6e, 0x67, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x68, 0x6f, 0x6f, 0x66,
	0x64, 0x72, 0x69, 0x63, 0x68, 0x74, 0x69, 0x6e, 0x67, 0x42, 0x68, 0x0a, 0x09, 0x63, 0x6f, 0x6d,
	0x2e, 0x70, 0x6f, 0x6c, 0x65, 0x73, 0x42, 0x09, 0x50, 0x6f, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x1c, 0x73, 0x74, 0x61, 0x67, 0x65, 0x32, 0x30, 0x32, 0x34, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x6f, 0x6c, 0x65,
	0x73, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x05, 0x50, 0x6f, 0x6c, 0x65, 0x73, 0xca,
	0x02, 0x05, 0x50, 0x6f, 0x6c, 0x65, 0x73, 0xe2, 0x02, 0x11, 0x50, 0x6f, 0x6c, 0x65, 0x73, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x05, 0x50, 0x6f,
	0x6c, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_poles_pole_proto_rawDescOnce sync.Once
	file_poles_pole_proto_rawDescData = file_poles_pole_proto_rawDesc
)

func file_poles_pole_proto_rawDescGZIP() []byte {
	file_poles_pole_proto_rawDescOnce.Do(func() {
		file_poles_pole_proto_rawDescData = protoimpl.X.CompressGZIP(file_poles_pole_proto_rawDescData)
	})
	return file_poles_pole_proto_rawDescData
}

var file_poles_pole_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_poles_pole_proto_goTypes = []interface{}{
	(*PoleData)(nil),              // 0: poles.PoleData
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_poles_pole_proto_depIdxs = []int32{
	1, // 0: poles.PoleData.ordening:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_poles_pole_proto_init() }
func file_poles_pole_proto_init() {
	if File_poles_pole_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_poles_pole_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PoleData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_poles_pole_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_poles_pole_proto_goTypes,
		DependencyIndexes: file_poles_pole_proto_depIdxs,
		MessageInfos:      file_poles_pole_proto_msgTypes,
	}.Build()
	File_poles_pole_proto = out.File
	file_poles_pole_proto_rawDesc = nil
	file_poles_pole_proto_goTypes = nil
	file_poles_pole_proto_depIdxs = nil
}