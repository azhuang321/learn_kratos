// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.0--rc1
// source: group/service/v1/group.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type GetGroupInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetGroupInfoRequest) Reset() {
	*x = GetGroupInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_service_v1_group_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetGroupInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGroupInfoRequest) ProtoMessage() {}

func (x *GetGroupInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_group_service_v1_group_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGroupInfoRequest.ProtoReflect.Descriptor instead.
func (*GetGroupInfoRequest) Descriptor() ([]byte, []int) {
	return file_group_service_v1_group_proto_rawDescGZIP(), []int{0}
}

func (x *GetGroupInfoRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetGroupInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	GroupName string `protobuf:"bytes,2,opt,name=group_name,json=groupName,proto3" json:"group_name,omitempty"`
}

func (x *GetGroupInfoResponse) Reset() {
	*x = GetGroupInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_service_v1_group_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetGroupInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGroupInfoResponse) ProtoMessage() {}

func (x *GetGroupInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_group_service_v1_group_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGroupInfoResponse.ProtoReflect.Descriptor instead.
func (*GetGroupInfoResponse) Descriptor() ([]byte, []int) {
	return file_group_service_v1_group_proto_rawDescGZIP(), []int{1}
}

func (x *GetGroupInfoResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GetGroupInfoResponse) GetGroupName() string {
	if x != nil {
		return x.GroupName
	}
	return ""
}

var File_group_service_v1_group_proto protoreflect.FileDescriptor

var file_group_service_v1_group_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f,
	0x76, 0x31, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14,
	0x61, 0x70, 0x69, 0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x25, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x45, 0x0a, 0x14, 0x47, 0x65, 0x74,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x61, 0x6d, 0x65,
	0x32, 0x88, 0x01, 0x0a, 0x05, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x7f, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x29, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x22, 0x0d, 0x2f, 0x76, 0x31, 0x2f, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x3a, 0x01, 0x2a, 0x42, 0x19, 0x5a, 0x17, 0x61,
	0x70, 0x69, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_group_service_v1_group_proto_rawDescOnce sync.Once
	file_group_service_v1_group_proto_rawDescData = file_group_service_v1_group_proto_rawDesc
)

func file_group_service_v1_group_proto_rawDescGZIP() []byte {
	file_group_service_v1_group_proto_rawDescOnce.Do(func() {
		file_group_service_v1_group_proto_rawDescData = protoimpl.X.CompressGZIP(file_group_service_v1_group_proto_rawDescData)
	})
	return file_group_service_v1_group_proto_rawDescData
}

var file_group_service_v1_group_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_group_service_v1_group_proto_goTypes = []interface{}{
	(*GetGroupInfoRequest)(nil),  // 0: api.group.service.v1.GetGroupInfoRequest
	(*GetGroupInfoResponse)(nil), // 1: api.group.service.v1.GetGroupInfoResponse
}
var file_group_service_v1_group_proto_depIdxs = []int32{
	0, // 0: api.group.service.v1.Group.GetGroupInfo:input_type -> api.group.service.v1.GetGroupInfoRequest
	1, // 1: api.group.service.v1.Group.GetGroupInfo:output_type -> api.group.service.v1.GetGroupInfoResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_group_service_v1_group_proto_init() }
func file_group_service_v1_group_proto_init() {
	if File_group_service_v1_group_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_group_service_v1_group_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetGroupInfoRequest); i {
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
		file_group_service_v1_group_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetGroupInfoResponse); i {
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
			RawDescriptor: file_group_service_v1_group_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_group_service_v1_group_proto_goTypes,
		DependencyIndexes: file_group_service_v1_group_proto_depIdxs,
		MessageInfos:      file_group_service_v1_group_proto_msgTypes,
	}.Build()
	File_group_service_v1_group_proto = out.File
	file_group_service_v1_group_proto_rawDesc = nil
	file_group_service_v1_group_proto_goTypes = nil
	file_group_service_v1_group_proto_depIdxs = nil
}
