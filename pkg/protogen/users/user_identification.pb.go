// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: users/user_identification.proto

package users

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UserIdentification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // uuid
	UserName     string `protobuf:"bytes,2,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	EmailAddress string `protobuf:"bytes,3,opt,name=email_address,json=emailAddress,proto3" json:"email_address,omitempty"`
}

func (x *UserIdentification) Reset() {
	*x = UserIdentification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_users_user_identification_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserIdentification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserIdentification) ProtoMessage() {}

func (x *UserIdentification) ProtoReflect() protoreflect.Message {
	mi := &file_users_user_identification_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserIdentification.ProtoReflect.Descriptor instead.
func (*UserIdentification) Descriptor() ([]byte, []int) {
	return file_users_user_identification_proto_rawDescGZIP(), []int{0}
}

func (x *UserIdentification) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UserIdentification) GetUserName() string {
	if x != nil {
		return x.UserName
	}
	return ""
}

func (x *UserIdentification) GetEmailAddress() string {
	if x != nil {
		return x.EmailAddress
	}
	return ""
}

var File_users_user_identification_proto protoreflect.FileDescriptor

var file_users_user_identification_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x67, 0x0a, 0x13, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x42, 0x39, 0x42, 0x17, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x1c, 0x73, 0x74, 0x61, 0x67, 0x65, 0x32, 0x30,
	0x32, 0x34, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_users_user_identification_proto_rawDescOnce sync.Once
	file_users_user_identification_proto_rawDescData = file_users_user_identification_proto_rawDesc
)

func file_users_user_identification_proto_rawDescGZIP() []byte {
	file_users_user_identification_proto_rawDescOnce.Do(func() {
		file_users_user_identification_proto_rawDescData = protoimpl.X.CompressGZIP(file_users_user_identification_proto_rawDescData)
	})
	return file_users_user_identification_proto_rawDescData
}

var file_users_user_identification_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_users_user_identification_proto_goTypes = []interface{}{
	(*UserIdentification)(nil), // 0: user_identification
}
var file_users_user_identification_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_users_user_identification_proto_init() }
func file_users_user_identification_proto_init() {
	if File_users_user_identification_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_users_user_identification_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserIdentification); i {
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
			RawDescriptor: file_users_user_identification_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_users_user_identification_proto_goTypes,
		DependencyIndexes: file_users_user_identification_proto_depIdxs,
		MessageInfos:      file_users_user_identification_proto_msgTypes,
	}.Build()
	File_users_user_identification_proto = out.File
	file_users_user_identification_proto_rawDesc = nil
	file_users_user_identification_proto_goTypes = nil
	file_users_user_identification_proto_depIdxs = nil
}
