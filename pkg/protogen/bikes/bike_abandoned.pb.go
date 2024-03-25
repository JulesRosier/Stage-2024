// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: bikes/bike_abandoned.proto

package bikes

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	common "stage2024/pkg/protogen/common"
	users "stage2024/pkg/protogen/users"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BikeAbandoned struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TimeStamp *timestamppb.Timestamp    `protobuf:"bytes,1,opt,name=time_stamp,json=timeStamp,proto3" json:"time_stamp,omitempty"`
	Bike      *AbandonedBike            `protobuf:"bytes,2,opt,name=bike,proto3" json:"bike,omitempty"`
	User      *users.UserIdentification `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *BikeAbandoned) Reset() {
	*x = BikeAbandoned{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bikes_bike_abandoned_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BikeAbandoned) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BikeAbandoned) ProtoMessage() {}

func (x *BikeAbandoned) ProtoReflect() protoreflect.Message {
	mi := &file_bikes_bike_abandoned_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BikeAbandoned.ProtoReflect.Descriptor instead.
func (*BikeAbandoned) Descriptor() ([]byte, []int) {
	return file_bikes_bike_abandoned_proto_rawDescGZIP(), []int{0}
}

func (x *BikeAbandoned) GetTimeStamp() *timestamppb.Timestamp {
	if x != nil {
		return x.TimeStamp
	}
	return nil
}

func (x *BikeAbandoned) GetBike() *AbandonedBike {
	if x != nil {
		return x.Bike
	}
	return nil
}

func (x *BikeAbandoned) GetUser() *users.UserIdentification {
	if x != nil {
		return x.User
	}
	return nil
}

type AbandonedBike struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bike     *BikeIdentification `protobuf:"bytes,1,opt,name=bike,proto3" json:"bike,omitempty"`
	Location *common.Location    `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *AbandonedBike) Reset() {
	*x = AbandonedBike{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bikes_bike_abandoned_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AbandonedBike) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AbandonedBike) ProtoMessage() {}

func (x *AbandonedBike) ProtoReflect() protoreflect.Message {
	mi := &file_bikes_bike_abandoned_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AbandonedBike.ProtoReflect.Descriptor instead.
func (*AbandonedBike) Descriptor() ([]byte, []int) {
	return file_bikes_bike_abandoned_proto_rawDescGZIP(), []int{1}
}

func (x *AbandonedBike) GetBike() *BikeIdentification {
	if x != nil {
		return x.Bike
	}
	return nil
}

func (x *AbandonedBike) GetLocation() *common.Location {
	if x != nil {
		return x.Location
	}
	return nil
}

var File_bikes_bike_abandoned_proto protoreflect.FileDescriptor

var file_bikes_bike_abandoned_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x62, 0x69, 0x6b, 0x65, 0x73, 0x2f, 0x62, 0x69, 0x6b, 0x65, 0x5f, 0x61, 0x62, 0x61,
	0x6e, 0x64, 0x6f, 0x6e, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x62, 0x69,
	0x6b, 0x65, 0x73, 0x2f, 0x62, 0x69, 0x6b, 0x65, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9a, 0x01, 0x0a, 0x0e, 0x62, 0x69, 0x6b, 0x65, 0x5f,
	0x61, 0x62, 0x61, 0x6e, 0x64, 0x6f, 0x6e, 0x65, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x74, 0x69, 0x6d,
	0x65, 0x5f, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x53,
	0x74, 0x61, 0x6d, 0x70, 0x12, 0x23, 0x0a, 0x04, 0x62, 0x69, 0x6b, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x62, 0x61, 0x6e, 0x64, 0x6f, 0x6e, 0x65, 0x64, 0x5f, 0x62,
	0x69, 0x6b, 0x65, 0x52, 0x04, 0x62, 0x69, 0x6b, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x75,
	0x73, 0x65, 0x72, 0x22, 0x61, 0x0a, 0x0e, 0x61, 0x62, 0x61, 0x6e, 0x64, 0x6f, 0x6e, 0x65, 0x64,
	0x5f, 0x62, 0x69, 0x6b, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x62, 0x69, 0x6b, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x62, 0x69, 0x6b, 0x65, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x62, 0x69, 0x6b, 0x65, 0x12,
	0x25, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x09, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x34, 0x42, 0x12, 0x42, 0x69, 0x6b, 0x65, 0x41, 0x62,
	0x61, 0x6e, 0x64, 0x6f, 0x6e, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x1c,
	0x73, 0x74, 0x61, 0x67, 0x65, 0x32, 0x30, 0x32, 0x34, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x62, 0x69, 0x6b, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bikes_bike_abandoned_proto_rawDescOnce sync.Once
	file_bikes_bike_abandoned_proto_rawDescData = file_bikes_bike_abandoned_proto_rawDesc
)

func file_bikes_bike_abandoned_proto_rawDescGZIP() []byte {
	file_bikes_bike_abandoned_proto_rawDescOnce.Do(func() {
		file_bikes_bike_abandoned_proto_rawDescData = protoimpl.X.CompressGZIP(file_bikes_bike_abandoned_proto_rawDescData)
	})
	return file_bikes_bike_abandoned_proto_rawDescData
}

var file_bikes_bike_abandoned_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_bikes_bike_abandoned_proto_goTypes = []interface{}{
	(*BikeAbandoned)(nil),            // 0: bike_abandoned
	(*AbandonedBike)(nil),            // 1: abandoned_bike
	(*timestamppb.Timestamp)(nil),    // 2: google.protobuf.Timestamp
	(*users.UserIdentification)(nil), // 3: user_identification
	(*BikeIdentification)(nil),       // 4: bike_identification
	(*common.Location)(nil),          // 5: location
}
var file_bikes_bike_abandoned_proto_depIdxs = []int32{
	2, // 0: bike_abandoned.time_stamp:type_name -> google.protobuf.Timestamp
	1, // 1: bike_abandoned.bike:type_name -> abandoned_bike
	3, // 2: bike_abandoned.user:type_name -> user_identification
	4, // 3: abandoned_bike.bike:type_name -> bike_identification
	5, // 4: abandoned_bike.location:type_name -> location
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_bikes_bike_abandoned_proto_init() }
func file_bikes_bike_abandoned_proto_init() {
	if File_bikes_bike_abandoned_proto != nil {
		return
	}
	file_bikes_bike_identification_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_bikes_bike_abandoned_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BikeAbandoned); i {
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
		file_bikes_bike_abandoned_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AbandonedBike); i {
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
			RawDescriptor: file_bikes_bike_abandoned_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_bikes_bike_abandoned_proto_goTypes,
		DependencyIndexes: file_bikes_bike_abandoned_proto_depIdxs,
		MessageInfos:      file_bikes_bike_abandoned_proto_msgTypes,
	}.Build()
	File_bikes_bike_abandoned_proto = out.File
	file_bikes_bike_abandoned_proto_rawDesc = nil
	file_bikes_bike_abandoned_proto_goTypes = nil
	file_bikes_bike_abandoned_proto_depIdxs = nil
}
