// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        (unknown)
// source: bikes/bolt.proto

package bikes

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	common "stage2024/pkg/protogen/common"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BoltLocation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BikeId             string           `protobuf:"bytes,1,opt,name=bike_id,json=bikeId,proto3" json:"bike_id,omitempty"`
	CurrentRangeMeters int32            `protobuf:"varint,2,opt,name=current_range_meters,json=currentRangeMeters,proto3" json:"current_range_meters,omitempty"`
	PricingPlanId      string           `protobuf:"bytes,3,opt,name=pricing_plan_id,json=pricingPlanId,proto3" json:"pricing_plan_id,omitempty"`
	VehicleTypeId      string           `protobuf:"bytes,4,opt,name=vehicle_type_id,json=vehicleTypeId,proto3" json:"vehicle_type_id,omitempty"`
	IsReserved         int32            `protobuf:"varint,5,opt,name=is_reserved,json=isReserved,proto3" json:"is_reserved,omitempty"`
	IsDisabled         int32            `protobuf:"varint,6,opt,name=is_disabled,json=isDisabled,proto3" json:"is_disabled,omitempty"`
	RentalUris         string           `protobuf:"bytes,7,opt,name=rental_uris,json=rentalUris,proto3" json:"rental_uris,omitempty"`
	Location           *common.Location `protobuf:"bytes,8,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *BoltLocation) Reset() {
	*x = BoltLocation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bikes_bolt_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoltLocation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoltLocation) ProtoMessage() {}

func (x *BoltLocation) ProtoReflect() protoreflect.Message {
	mi := &file_bikes_bolt_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoltLocation.ProtoReflect.Descriptor instead.
func (*BoltLocation) Descriptor() ([]byte, []int) {
	return file_bikes_bolt_proto_rawDescGZIP(), []int{0}
}

func (x *BoltLocation) GetBikeId() string {
	if x != nil {
		return x.BikeId
	}
	return ""
}

func (x *BoltLocation) GetCurrentRangeMeters() int32 {
	if x != nil {
		return x.CurrentRangeMeters
	}
	return 0
}

func (x *BoltLocation) GetPricingPlanId() string {
	if x != nil {
		return x.PricingPlanId
	}
	return ""
}

func (x *BoltLocation) GetVehicleTypeId() string {
	if x != nil {
		return x.VehicleTypeId
	}
	return ""
}

func (x *BoltLocation) GetIsReserved() int32 {
	if x != nil {
		return x.IsReserved
	}
	return 0
}

func (x *BoltLocation) GetIsDisabled() int32 {
	if x != nil {
		return x.IsDisabled
	}
	return 0
}

func (x *BoltLocation) GetRentalUris() string {
	if x != nil {
		return x.RentalUris
	}
	return ""
}

func (x *BoltLocation) GetLocation() *common.Location {
	if x != nil {
		return x.Location
	}
	return nil
}

var File_bikes_bolt_proto protoreflect.FileDescriptor

var file_bikes_bolt_proto_rawDesc = []byte{
	0x0a, 0x10, 0x62, 0x69, 0x6b, 0x65, 0x73, 0x2f, 0x62, 0x6f, 0x6c, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x62, 0x69, 0x6b, 0x65, 0x73, 0x1a, 0x15, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xba, 0x02, 0x0a, 0x0c, 0x42, 0x6f, 0x6c, 0x74, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x17, 0x0a, 0x07, 0x62, 0x69, 0x6b, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x62, 0x69, 0x6b, 0x65, 0x49, 0x64, 0x12, 0x30, 0x0a, 0x14, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x12, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12, 0x26, 0x0a, 0x0f,
	0x70, 0x72, 0x69, 0x63, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x70, 0x72, 0x69, 0x63, 0x69, 0x6e, 0x67, 0x50, 0x6c,
	0x61, 0x6e, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x76, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x76,
	0x65, 0x68, 0x69, 0x63, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b,
	0x69, 0x73, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0a, 0x69, 0x73, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x12, 0x1f, 0x0a,
	0x0b, 0x69, 0x73, 0x5f, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x69, 0x73, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x1f,
	0x0a, 0x0b, 0x72, 0x65, 0x6e, 0x74, 0x61, 0x6c, 0x5f, 0x75, 0x72, 0x69, 0x73, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x6e, 0x74, 0x61, 0x6c, 0x55, 0x72, 0x69, 0x73, 0x12,
	0x2c, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x68, 0x0a,
	0x09, 0x63, 0x6f, 0x6d, 0x2e, 0x62, 0x69, 0x6b, 0x65, 0x73, 0x42, 0x09, 0x42, 0x6f, 0x6c, 0x74,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x1c, 0x73, 0x74, 0x61, 0x67, 0x65, 0x32, 0x30,
	0x32, 0x34, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f,
	0x62, 0x69, 0x6b, 0x65, 0x73, 0xa2, 0x02, 0x03, 0x42, 0x58, 0x58, 0xaa, 0x02, 0x05, 0x42, 0x69,
	0x6b, 0x65, 0x73, 0xca, 0x02, 0x05, 0x42, 0x69, 0x6b, 0x65, 0x73, 0xe2, 0x02, 0x11, 0x42, 0x69,
	0x6b, 0x65, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x05, 0x42, 0x69, 0x6b, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bikes_bolt_proto_rawDescOnce sync.Once
	file_bikes_bolt_proto_rawDescData = file_bikes_bolt_proto_rawDesc
)

func file_bikes_bolt_proto_rawDescGZIP() []byte {
	file_bikes_bolt_proto_rawDescOnce.Do(func() {
		file_bikes_bolt_proto_rawDescData = protoimpl.X.CompressGZIP(file_bikes_bolt_proto_rawDescData)
	})
	return file_bikes_bolt_proto_rawDescData
}

var file_bikes_bolt_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_bikes_bolt_proto_goTypes = []interface{}{
	(*BoltLocation)(nil),    // 0: bikes.BoltLocation
	(*common.Location)(nil), // 1: common.Location
}
var file_bikes_bolt_proto_depIdxs = []int32{
	1, // 0: bikes.BoltLocation.location:type_name -> common.Location
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_bikes_bolt_proto_init() }
func file_bikes_bolt_proto_init() {
	if File_bikes_bolt_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bikes_bolt_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoltLocation); i {
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
			RawDescriptor: file_bikes_bolt_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_bikes_bolt_proto_goTypes,
		DependencyIndexes: file_bikes_bolt_proto_depIdxs,
		MessageInfos:      file_bikes_bolt_proto_msgTypes,
	}.Build()
	File_bikes_bolt_proto = out.File
	file_bikes_bolt_proto_rawDesc = nil
	file_bikes_bolt_proto_goTypes = nil
	file_bikes_bolt_proto_depIdxs = nil
}
