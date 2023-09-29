// mapper.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0-devel
// 	protoc        v4.23.4
// source: mapper.proto

package mapper

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

type MapperInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageRank      float32 `protobuf:"fixed32,1,opt,name=page_rank,json=pageRank,proto3" json:"page_rank,omitempty"`
	AdjacencyList []int32 `protobuf:"varint,2,rep,packed,name=adjacency_list,json=adjacencyList,proto3" json:"adjacency_list,omitempty"`
}

func (x *MapperInput) Reset() {
	*x = MapperInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mapper_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapperInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapperInput) ProtoMessage() {}

func (x *MapperInput) ProtoReflect() protoreflect.Message {
	mi := &file_mapper_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapperInput.ProtoReflect.Descriptor instead.
func (*MapperInput) Descriptor() ([]byte, []int) {
	return file_mapper_proto_rawDescGZIP(), []int{0}
}

func (x *MapperInput) GetPageRank() float32 {
	if x != nil {
		return x.PageRank
	}
	return 0
}

func (x *MapperInput) GetAdjacencyList() []int32 {
	if x != nil {
		return x.AdjacencyList
	}
	return nil
}

type MapperOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageRankShare float32 `protobuf:"fixed32,1,opt,name=page_rank_share,json=pageRankShare,proto3" json:"page_rank_share,omitempty"`
	AdjacencyList []int32 `protobuf:"varint,2,rep,packed,name=adjacency_list,json=adjacencyList,proto3" json:"adjacency_list,omitempty"`
}

func (x *MapperOutput) Reset() {
	*x = MapperOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mapper_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapperOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapperOutput) ProtoMessage() {}

func (x *MapperOutput) ProtoReflect() protoreflect.Message {
	mi := &file_mapper_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapperOutput.ProtoReflect.Descriptor instead.
func (*MapperOutput) Descriptor() ([]byte, []int) {
	return file_mapper_proto_rawDescGZIP(), []int{1}
}

func (x *MapperOutput) GetPageRankShare() float32 {
	if x != nil {
		return x.PageRankShare
	}
	return 0
}

func (x *MapperOutput) GetAdjacencyList() []int32 {
	if x != nil {
		return x.AdjacencyList
	}
	return nil
}

type CleanUpInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageRank      float32 `protobuf:"fixed32,1,opt,name=page_rank,json=pageRank,proto3" json:"page_rank,omitempty"`
	AdjacencyList []int32 `protobuf:"varint,2,rep,packed,name=adjacency_list,json=adjacencyList,proto3" json:"adjacency_list,omitempty"`
}

func (x *CleanUpInput) Reset() {
	*x = CleanUpInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mapper_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CleanUpInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CleanUpInput) ProtoMessage() {}

func (x *CleanUpInput) ProtoReflect() protoreflect.Message {
	mi := &file_mapper_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CleanUpInput.ProtoReflect.Descriptor instead.
func (*CleanUpInput) Descriptor() ([]byte, []int) {
	return file_mapper_proto_rawDescGZIP(), []int{2}
}

func (x *CleanUpInput) GetPageRank() float32 {
	if x != nil {
		return x.PageRank
	}
	return 0
}

func (x *CleanUpInput) GetAdjacencyList() []int32 {
	if x != nil {
		return x.AdjacencyList
	}
	return nil
}

type CleanUpOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SinkMass float32 `protobuf:"fixed32,1,opt,name=sink_mass,json=sinkMass,proto3" json:"sink_mass,omitempty"`
}

func (x *CleanUpOutput) Reset() {
	*x = CleanUpOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mapper_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CleanUpOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CleanUpOutput) ProtoMessage() {}

func (x *CleanUpOutput) ProtoReflect() protoreflect.Message {
	mi := &file_mapper_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CleanUpOutput.ProtoReflect.Descriptor instead.
func (*CleanUpOutput) Descriptor() ([]byte, []int) {
	return file_mapper_proto_rawDescGZIP(), []int{3}
}

func (x *CleanUpOutput) GetSinkMass() float32 {
	if x != nil {
		return x.SinkMass
	}
	return 0
}

type MapperHeartbeatRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Alive bool `protobuf:"varint,1,opt,name=alive,proto3" json:"alive,omitempty"`
}

func (x *MapperHeartbeatRequest) Reset() {
	*x = MapperHeartbeatRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mapper_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapperHeartbeatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapperHeartbeatRequest) ProtoMessage() {}

func (x *MapperHeartbeatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mapper_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapperHeartbeatRequest.ProtoReflect.Descriptor instead.
func (*MapperHeartbeatRequest) Descriptor() ([]byte, []int) {
	return file_mapper_proto_rawDescGZIP(), []int{4}
}

func (x *MapperHeartbeatRequest) GetAlive() bool {
	if x != nil {
		return x.Alive
	}
	return false
}

type MapperHeartbeatResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Alive bool `protobuf:"varint,1,opt,name=alive,proto3" json:"alive,omitempty"`
}

func (x *MapperHeartbeatResponse) Reset() {
	*x = MapperHeartbeatResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mapper_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapperHeartbeatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapperHeartbeatResponse) ProtoMessage() {}

func (x *MapperHeartbeatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mapper_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapperHeartbeatResponse.ProtoReflect.Descriptor instead.
func (*MapperHeartbeatResponse) Descriptor() ([]byte, []int) {
	return file_mapper_proto_rawDescGZIP(), []int{5}
}

func (x *MapperHeartbeatResponse) GetAlive() bool {
	if x != nil {
		return x.Alive
	}
	return false
}

var File_mapper_proto protoreflect.FileDescriptor

var file_mapper_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x51,
	0x0a, 0x0b, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02,
	0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x64,
	0x6a, 0x61, 0x63, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x05, 0x52, 0x0d, 0x61, 0x64, 0x6a, 0x61, 0x63, 0x65, 0x6e, 0x63, 0x79, 0x4c, 0x69, 0x73,
	0x74, 0x22, 0x5d, 0x0a, 0x0c, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x12, 0x26, 0x0a, 0x0f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x6e, 0x6b, 0x5f, 0x73,
	0x68, 0x61, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0d, 0x70, 0x61, 0x67, 0x65,
	0x52, 0x61, 0x6e, 0x6b, 0x53, 0x68, 0x61, 0x72, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x64, 0x6a,
	0x61, 0x63, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x05, 0x52, 0x0d, 0x61, 0x64, 0x6a, 0x61, 0x63, 0x65, 0x6e, 0x63, 0x79, 0x4c, 0x69, 0x73, 0x74,
	0x22, 0x52, 0x0a, 0x0c, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x55, 0x70, 0x49, 0x6e, 0x70, 0x75, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x6e, 0x6b, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x12, 0x25, 0x0a,
	0x0e, 0x61, 0x64, 0x6a, 0x61, 0x63, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x05, 0x52, 0x0d, 0x61, 0x64, 0x6a, 0x61, 0x63, 0x65, 0x6e, 0x63, 0x79,
	0x4c, 0x69, 0x73, 0x74, 0x22, 0x2c, 0x0a, 0x0d, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x55, 0x70, 0x4f,
	0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x69, 0x6e, 0x6b, 0x5f, 0x6d, 0x61,
	0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x08, 0x73, 0x69, 0x6e, 0x6b, 0x4d, 0x61,
	0x73, 0x73, 0x22, 0x2e, 0x0a, 0x16, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72,
	0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x61, 0x6c, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x61, 0x6c, 0x69,
	0x76, 0x65, 0x22, 0x2f, 0x0a, 0x17, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72,
	0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x61, 0x6c, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x61, 0x6c,
	0x69, 0x76, 0x65, 0x32, 0x56, 0x0a, 0x06, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x22, 0x0a,
	0x03, 0x4d, 0x61, 0x70, 0x12, 0x0c, 0x2e, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x1a, 0x0d, 0x2e, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x12, 0x28, 0x0a, 0x07, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x55, 0x70, 0x12, 0x0d, 0x2e, 0x43,
	0x6c, 0x65, 0x61, 0x6e, 0x55, 0x70, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x0e, 0x2e, 0x43, 0x6c,
	0x65, 0x61, 0x6e, 0x55, 0x70, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x32, 0x4c, 0x0a, 0x0f, 0x4d,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x39,
	0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x17, 0x2e, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x48,
	0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x18, 0x2e, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x6d,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mapper_proto_rawDescOnce sync.Once
	file_mapper_proto_rawDescData = file_mapper_proto_rawDesc
)

func file_mapper_proto_rawDescGZIP() []byte {
	file_mapper_proto_rawDescOnce.Do(func() {
		file_mapper_proto_rawDescData = protoimpl.X.CompressGZIP(file_mapper_proto_rawDescData)
	})
	return file_mapper_proto_rawDescData
}

var file_mapper_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_mapper_proto_goTypes = []interface{}{
	(*MapperInput)(nil),             // 0: MapperInput
	(*MapperOutput)(nil),            // 1: MapperOutput
	(*CleanUpInput)(nil),            // 2: CleanUpInput
	(*CleanUpOutput)(nil),           // 3: CleanUpOutput
	(*MapperHeartbeatRequest)(nil),  // 4: MapperHeartbeatRequest
	(*MapperHeartbeatResponse)(nil), // 5: MapperHeartbeatResponse
}
var file_mapper_proto_depIdxs = []int32{
	0, // 0: Mapper.Map:input_type -> MapperInput
	2, // 1: Mapper.CleanUp:input_type -> CleanUpInput
	4, // 2: MapperHeartbeat.Ping:input_type -> MapperHeartbeatRequest
	1, // 3: Mapper.Map:output_type -> MapperOutput
	3, // 4: Mapper.CleanUp:output_type -> CleanUpOutput
	5, // 5: MapperHeartbeat.Ping:output_type -> MapperHeartbeatResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mapper_proto_init() }
func file_mapper_proto_init() {
	if File_mapper_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mapper_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapperInput); i {
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
		file_mapper_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapperOutput); i {
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
		file_mapper_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CleanUpInput); i {
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
		file_mapper_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CleanUpOutput); i {
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
		file_mapper_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapperHeartbeatRequest); i {
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
		file_mapper_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapperHeartbeatResponse); i {
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
			RawDescriptor: file_mapper_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_mapper_proto_goTypes,
		DependencyIndexes: file_mapper_proto_depIdxs,
		MessageInfos:      file_mapper_proto_msgTypes,
	}.Build()
	File_mapper_proto = out.File
	file_mapper_proto_rawDesc = nil
	file_mapper_proto_goTypes = nil
	file_mapper_proto_depIdxs = nil
}
