// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: stations/station_full.proto

package stations

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

type StationFull struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TimeStamp   *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=time_stamp,json=timeStamp,proto3" json:"time_stamp,omitempty"`
	Station     *StationIdentification `protobuf:"bytes,2,opt,name=station,proto3" json:"station,omitempty"`
	MaxCapacity int32                  `protobuf:"varint,3,opt,name=max_capacity,json=maxCapacity,proto3" json:"max_capacity,omitempty"`
}

func (x *StationFull) Reset() {
	*x = StationFull{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stations_station_full_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StationFull) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StationFull) ProtoMessage() {}

func (x *StationFull) ProtoReflect() protoreflect.Message {
	mi := &file_stations_station_full_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StationFull.ProtoReflect.Descriptor instead.
func (*StationFull) Descriptor() ([]byte, []int) {
	return file_stations_station_full_proto_rawDescGZIP(), []int{0}
}

func (x *StationFull) GetTimeStamp() *timestamppb.Timestamp {
	if x != nil {
		return x.TimeStamp
	}
	return nil
}

func (x *StationFull) GetStation() *StationIdentification {
	if x != nil {
		return x.Station
	}
	return nil
}

func (x *StationFull) GetMaxCapacity() int32 {
	if x != nil {
		return x.MaxCapacity
	}
	return 0
}

var File_stations_station_full_proto protoreflect.FileDescriptor

var file_stations_station_full_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x66, 0x75, 0x6c, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25,
	0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9f, 0x01, 0x0a, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x66, 0x75, 0x6c, 0x6c, 0x12, 0x39, 0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x31, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x73, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x61, 0x78, 0x5f, 0x63, 0x61, 0x70, 0x61,
	0x63, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x6d, 0x61, 0x78, 0x43,
	0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x42, 0x35, 0x42, 0x10, 0x53, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x46, 0x75, 0x6c, 0x6c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x1f, 0x73,
	0x74, 0x61, 0x67, 0x65, 0x32, 0x30, 0x32, 0x34, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_stations_station_full_proto_rawDescOnce sync.Once
	file_stations_station_full_proto_rawDescData = file_stations_station_full_proto_rawDesc
)

func file_stations_station_full_proto_rawDescGZIP() []byte {
	file_stations_station_full_proto_rawDescOnce.Do(func() {
		file_stations_station_full_proto_rawDescData = protoimpl.X.CompressGZIP(file_stations_station_full_proto_rawDescData)
	})
	return file_stations_station_full_proto_rawDescData
}

var file_stations_station_full_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_stations_station_full_proto_goTypes = []interface{}{
	(*StationFull)(nil),           // 0: station_full
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*StationIdentification)(nil), // 2: station_identification
}
var file_stations_station_full_proto_depIdxs = []int32{
	1, // 0: station_full.time_stamp:type_name -> google.protobuf.Timestamp
	2, // 1: station_full.station:type_name -> station_identification
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_stations_station_full_proto_init() }
func file_stations_station_full_proto_init() {
	if File_stations_station_full_proto != nil {
		return
	}
	file_stations_station_identification_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_stations_station_full_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StationFull); i {
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
			RawDescriptor: file_stations_station_full_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_stations_station_full_proto_goTypes,
		DependencyIndexes: file_stations_station_full_proto_depIdxs,
		MessageInfos:      file_stations_station_full_proto_msgTypes,
	}.Build()
	File_stations_station_full_proto = out.File
	file_stations_station_full_proto_rawDesc = nil
	file_stations_station_full_proto_goTypes = nil
	file_stations_station_full_proto_depIdxs = nil
}
