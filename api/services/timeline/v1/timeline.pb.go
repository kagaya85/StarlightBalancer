// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: services/timeline/v1/timeline.proto

package v1

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

type PushTimelineRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Items  []string `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *PushTimelineRequest) Reset() {
	*x = PushTimelineRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_timeline_v1_timeline_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushTimelineRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushTimelineRequest) ProtoMessage() {}

func (x *PushTimelineRequest) ProtoReflect() protoreflect.Message {
	mi := &file_services_timeline_v1_timeline_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushTimelineRequest.ProtoReflect.Descriptor instead.
func (*PushTimelineRequest) Descriptor() ([]byte, []int) {
	return file_services_timeline_v1_timeline_proto_rawDescGZIP(), []int{0}
}

func (x *PushTimelineRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *PushTimelineRequest) GetItems() []string {
	if x != nil {
		return x.Items
	}
	return nil
}

type PushTimelineResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status int32 `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *PushTimelineResponse) Reset() {
	*x = PushTimelineResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_services_timeline_v1_timeline_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushTimelineResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushTimelineResponse) ProtoMessage() {}

func (x *PushTimelineResponse) ProtoReflect() protoreflect.Message {
	mi := &file_services_timeline_v1_timeline_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushTimelineResponse.ProtoReflect.Descriptor instead.
func (*PushTimelineResponse) Descriptor() ([]byte, []int) {
	return file_services_timeline_v1_timeline_proto_rawDescGZIP(), []int{1}
}

func (x *PushTimelineResponse) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

var File_services_timeline_v1_timeline_proto protoreflect.FileDescriptor

var file_services_timeline_v1_timeline_proto_rawDesc = []byte{
	0x0a, 0x23, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x6c,
	0x69, 0x6e, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e,
	0x76, 0x31, 0x22, 0x44, 0x0a, 0x13, 0x50, 0x75, 0x73, 0x68, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69,
	0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x2e, 0x0a, 0x14, 0x50, 0x75, 0x73, 0x68,
	0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0x68, 0x0a, 0x0f, 0x54, 0x69, 0x6d, 0x65,
	0x6c, 0x69, 0x6e, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x55, 0x0a, 0x0c, 0x50,
	0x75, 0x73, 0x68, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x20, 0x2e, 0x74, 0x69,
	0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x54, 0x69,
	0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e,
	0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x73, 0x68,
	0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x1d, 0x5a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_services_timeline_v1_timeline_proto_rawDescOnce sync.Once
	file_services_timeline_v1_timeline_proto_rawDescData = file_services_timeline_v1_timeline_proto_rawDesc
)

func file_services_timeline_v1_timeline_proto_rawDescGZIP() []byte {
	file_services_timeline_v1_timeline_proto_rawDescOnce.Do(func() {
		file_services_timeline_v1_timeline_proto_rawDescData = protoimpl.X.CompressGZIP(file_services_timeline_v1_timeline_proto_rawDescData)
	})
	return file_services_timeline_v1_timeline_proto_rawDescData
}

var file_services_timeline_v1_timeline_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_services_timeline_v1_timeline_proto_goTypes = []interface{}{
	(*PushTimelineRequest)(nil),  // 0: timeline.v1.PushTimelineRequest
	(*PushTimelineResponse)(nil), // 1: timeline.v1.PushTimelineResponse
}
var file_services_timeline_v1_timeline_proto_depIdxs = []int32{
	0, // 0: timeline.v1.TimelineService.PushTimeline:input_type -> timeline.v1.PushTimelineRequest
	1, // 1: timeline.v1.TimelineService.PushTimeline:output_type -> timeline.v1.PushTimelineResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_services_timeline_v1_timeline_proto_init() }
func file_services_timeline_v1_timeline_proto_init() {
	if File_services_timeline_v1_timeline_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_services_timeline_v1_timeline_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushTimelineRequest); i {
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
		file_services_timeline_v1_timeline_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushTimelineResponse); i {
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
			RawDescriptor: file_services_timeline_v1_timeline_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_services_timeline_v1_timeline_proto_goTypes,
		DependencyIndexes: file_services_timeline_v1_timeline_proto_depIdxs,
		MessageInfos:      file_services_timeline_v1_timeline_proto_msgTypes,
	}.Build()
	File_services_timeline_v1_timeline_proto = out.File
	file_services_timeline_v1_timeline_proto_rawDesc = nil
	file_services_timeline_v1_timeline_proto_goTypes = nil
	file_services_timeline_v1_timeline_proto_depIdxs = nil
}