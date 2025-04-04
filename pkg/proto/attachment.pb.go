// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: attachment.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetAttachmentRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            *uint64                `protobuf:"varint,1,opt,name=id,proto3,oneof" json:"id,omitempty"`
	Rid           *string                `protobuf:"bytes,2,opt,name=rid,proto3,oneof" json:"rid,omitempty"`
	UserId        *uint64                `protobuf:"varint,3,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAttachmentRequest) Reset() {
	*x = GetAttachmentRequest{}
	mi := &file_attachment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAttachmentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAttachmentRequest) ProtoMessage() {}

func (x *GetAttachmentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_attachment_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAttachmentRequest.ProtoReflect.Descriptor instead.
func (*GetAttachmentRequest) Descriptor() ([]byte, []int) {
	return file_attachment_proto_rawDescGZIP(), []int{0}
}

func (x *GetAttachmentRequest) GetId() uint64 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return 0
}

func (x *GetAttachmentRequest) GetRid() string {
	if x != nil && x.Rid != nil {
		return *x.Rid
	}
	return ""
}

func (x *GetAttachmentRequest) GetUserId() uint64 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

type GetAttachmentResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Attachment    []byte                 `protobuf:"bytes,1,opt,name=attachment,proto3,oneof" json:"attachment,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAttachmentResponse) Reset() {
	*x = GetAttachmentResponse{}
	mi := &file_attachment_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAttachmentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAttachmentResponse) ProtoMessage() {}

func (x *GetAttachmentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_attachment_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAttachmentResponse.ProtoReflect.Descriptor instead.
func (*GetAttachmentResponse) Descriptor() ([]byte, []int) {
	return file_attachment_proto_rawDescGZIP(), []int{1}
}

func (x *GetAttachmentResponse) GetAttachment() []byte {
	if x != nil {
		return x.Attachment
	}
	return nil
}

type ListAttachmentRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            []uint64               `protobuf:"varint,1,rep,packed,name=id,proto3" json:"id,omitempty"`
	Rid           []string               `protobuf:"bytes,2,rep,name=rid,proto3" json:"rid,omitempty"`
	UserId        *uint64                `protobuf:"varint,3,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListAttachmentRequest) Reset() {
	*x = ListAttachmentRequest{}
	mi := &file_attachment_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListAttachmentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAttachmentRequest) ProtoMessage() {}

func (x *ListAttachmentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_attachment_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAttachmentRequest.ProtoReflect.Descriptor instead.
func (*ListAttachmentRequest) Descriptor() ([]byte, []int) {
	return file_attachment_proto_rawDescGZIP(), []int{2}
}

func (x *ListAttachmentRequest) GetId() []uint64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListAttachmentRequest) GetRid() []string {
	if x != nil {
		return x.Rid
	}
	return nil
}

func (x *ListAttachmentRequest) GetUserId() uint64 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

type ListAttachmentResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Attachments   [][]byte               `protobuf:"bytes,1,rep,name=attachments,proto3" json:"attachments,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListAttachmentResponse) Reset() {
	*x = ListAttachmentResponse{}
	mi := &file_attachment_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListAttachmentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAttachmentResponse) ProtoMessage() {}

func (x *ListAttachmentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_attachment_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAttachmentResponse.ProtoReflect.Descriptor instead.
func (*ListAttachmentResponse) Descriptor() ([]byte, []int) {
	return file_attachment_proto_rawDescGZIP(), []int{3}
}

func (x *ListAttachmentResponse) GetAttachments() [][]byte {
	if x != nil {
		return x.Attachments
	}
	return nil
}

type UpdateVisibilityRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            []uint64               `protobuf:"varint,1,rep,packed,name=id,proto3" json:"id,omitempty"`
	Rid           []string               `protobuf:"bytes,2,rep,name=rid,proto3" json:"rid,omitempty"`
	IsIndexable   bool                   `protobuf:"varint,3,opt,name=is_indexable,json=isIndexable,proto3" json:"is_indexable,omitempty"`
	UserId        *uint64                `protobuf:"varint,4,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateVisibilityRequest) Reset() {
	*x = UpdateVisibilityRequest{}
	mi := &file_attachment_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateVisibilityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateVisibilityRequest) ProtoMessage() {}

func (x *UpdateVisibilityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_attachment_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateVisibilityRequest.ProtoReflect.Descriptor instead.
func (*UpdateVisibilityRequest) Descriptor() ([]byte, []int) {
	return file_attachment_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateVisibilityRequest) GetId() []uint64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *UpdateVisibilityRequest) GetRid() []string {
	if x != nil {
		return x.Rid
	}
	return nil
}

func (x *UpdateVisibilityRequest) GetIsIndexable() bool {
	if x != nil {
		return x.IsIndexable
	}
	return false
}

func (x *UpdateVisibilityRequest) GetUserId() uint64 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

type UpdateVisibilityResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Count         int32                  `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateVisibilityResponse) Reset() {
	*x = UpdateVisibilityResponse{}
	mi := &file_attachment_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateVisibilityResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateVisibilityResponse) ProtoMessage() {}

func (x *UpdateVisibilityResponse) ProtoReflect() protoreflect.Message {
	mi := &file_attachment_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateVisibilityResponse.ProtoReflect.Descriptor instead.
func (*UpdateVisibilityResponse) Descriptor() ([]byte, []int) {
	return file_attachment_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateVisibilityResponse) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type UpdateUsageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            []uint64               `protobuf:"varint,1,rep,packed,name=id,proto3" json:"id,omitempty"`
	Rid           []string               `protobuf:"bytes,2,rep,name=rid,proto3" json:"rid,omitempty"`
	Delta         int64                  `protobuf:"varint,3,opt,name=delta,proto3" json:"delta,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateUsageRequest) Reset() {
	*x = UpdateUsageRequest{}
	mi := &file_attachment_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateUsageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUsageRequest) ProtoMessage() {}

func (x *UpdateUsageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_attachment_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUsageRequest.ProtoReflect.Descriptor instead.
func (*UpdateUsageRequest) Descriptor() ([]byte, []int) {
	return file_attachment_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateUsageRequest) GetId() []uint64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *UpdateUsageRequest) GetRid() []string {
	if x != nil {
		return x.Rid
	}
	return nil
}

func (x *UpdateUsageRequest) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

type UpdateUsageResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Count         int32                  `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateUsageResponse) Reset() {
	*x = UpdateUsageResponse{}
	mi := &file_attachment_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateUsageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUsageResponse) ProtoMessage() {}

func (x *UpdateUsageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_attachment_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUsageResponse.ProtoReflect.Descriptor instead.
func (*UpdateUsageResponse) Descriptor() ([]byte, []int) {
	return file_attachment_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateUsageResponse) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type DeleteAttachmentRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            []uint64               `protobuf:"varint,1,rep,packed,name=id,proto3" json:"id,omitempty"`
	Rid           []string               `protobuf:"bytes,2,rep,name=rid,proto3" json:"rid,omitempty"`
	UserId        *uint64                `protobuf:"varint,3,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteAttachmentRequest) Reset() {
	*x = DeleteAttachmentRequest{}
	mi := &file_attachment_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteAttachmentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAttachmentRequest) ProtoMessage() {}

func (x *DeleteAttachmentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_attachment_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAttachmentRequest.ProtoReflect.Descriptor instead.
func (*DeleteAttachmentRequest) Descriptor() ([]byte, []int) {
	return file_attachment_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteAttachmentRequest) GetId() []uint64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *DeleteAttachmentRequest) GetRid() []string {
	if x != nil {
		return x.Rid
	}
	return nil
}

func (x *DeleteAttachmentRequest) GetUserId() uint64 {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return 0
}

type DeleteAttachmentResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Count         int32                  `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteAttachmentResponse) Reset() {
	*x = DeleteAttachmentResponse{}
	mi := &file_attachment_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteAttachmentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAttachmentResponse) ProtoMessage() {}

func (x *DeleteAttachmentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_attachment_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAttachmentResponse.ProtoReflect.Descriptor instead.
func (*DeleteAttachmentResponse) Descriptor() ([]byte, []int) {
	return file_attachment_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteAttachmentResponse) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

var File_attachment_proto protoreflect.FileDescriptor

const file_attachment_proto_rawDesc = "" +
	"\n" +
	"\x10attachment.proto\x12\x05proto\"{\n" +
	"\x14GetAttachmentRequest\x12\x13\n" +
	"\x02id\x18\x01 \x01(\x04H\x00R\x02id\x88\x01\x01\x12\x15\n" +
	"\x03rid\x18\x02 \x01(\tH\x01R\x03rid\x88\x01\x01\x12\x1c\n" +
	"\auser_id\x18\x03 \x01(\x04H\x02R\x06userId\x88\x01\x01B\x05\n" +
	"\x03_idB\x06\n" +
	"\x04_ridB\n" +
	"\n" +
	"\b_user_id\"K\n" +
	"\x15GetAttachmentResponse\x12#\n" +
	"\n" +
	"attachment\x18\x01 \x01(\fH\x00R\n" +
	"attachment\x88\x01\x01B\r\n" +
	"\v_attachment\"c\n" +
	"\x15ListAttachmentRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x03(\x04R\x02id\x12\x10\n" +
	"\x03rid\x18\x02 \x03(\tR\x03rid\x12\x1c\n" +
	"\auser_id\x18\x03 \x01(\x04H\x00R\x06userId\x88\x01\x01B\n" +
	"\n" +
	"\b_user_id\":\n" +
	"\x16ListAttachmentResponse\x12 \n" +
	"\vattachments\x18\x01 \x03(\fR\vattachments\"\x88\x01\n" +
	"\x17UpdateVisibilityRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x03(\x04R\x02id\x12\x10\n" +
	"\x03rid\x18\x02 \x03(\tR\x03rid\x12!\n" +
	"\fis_indexable\x18\x03 \x01(\bR\visIndexable\x12\x1c\n" +
	"\auser_id\x18\x04 \x01(\x04H\x00R\x06userId\x88\x01\x01B\n" +
	"\n" +
	"\b_user_id\"0\n" +
	"\x18UpdateVisibilityResponse\x12\x14\n" +
	"\x05count\x18\x01 \x01(\x05R\x05count\"L\n" +
	"\x12UpdateUsageRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x03(\x04R\x02id\x12\x10\n" +
	"\x03rid\x18\x02 \x03(\tR\x03rid\x12\x14\n" +
	"\x05delta\x18\x03 \x01(\x03R\x05delta\"+\n" +
	"\x13UpdateUsageResponse\x12\x14\n" +
	"\x05count\x18\x01 \x01(\x05R\x05count\"e\n" +
	"\x17DeleteAttachmentRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x03(\x04R\x02id\x12\x10\n" +
	"\x03rid\x18\x02 \x03(\tR\x03rid\x12\x1c\n" +
	"\auser_id\x18\x03 \x01(\x04H\x00R\x06userId\x88\x01\x01B\n" +
	"\n" +
	"\b_user_id\"0\n" +
	"\x18DeleteAttachmentResponse\x12\x14\n" +
	"\x05count\x18\x01 \x01(\x05R\x05count2\xa8\x03\n" +
	"\x11AttachmentService\x12L\n" +
	"\rGetAttachment\x12\x1b.proto.GetAttachmentRequest\x1a\x1c.proto.GetAttachmentResponse\"\x00\x12O\n" +
	"\x0eListAttachment\x12\x1c.proto.ListAttachmentRequest\x1a\x1d.proto.ListAttachmentResponse\"\x00\x12U\n" +
	"\x10UpdateVisibility\x12\x1e.proto.UpdateVisibilityRequest\x1a\x1f.proto.UpdateVisibilityResponse\"\x00\x12F\n" +
	"\vUpdateUsage\x12\x19.proto.UpdateUsageRequest\x1a\x1a.proto.UpdateUsageResponse\"\x00\x12U\n" +
	"\x10DeleteAttachment\x12\x1e.proto.DeleteAttachmentRequest\x1a\x1f.proto.DeleteAttachmentResponse\"\x00B\tZ\a.;protob\x06proto3"

var (
	file_attachment_proto_rawDescOnce sync.Once
	file_attachment_proto_rawDescData []byte
)

func file_attachment_proto_rawDescGZIP() []byte {
	file_attachment_proto_rawDescOnce.Do(func() {
		file_attachment_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_attachment_proto_rawDesc), len(file_attachment_proto_rawDesc)))
	})
	return file_attachment_proto_rawDescData
}

var file_attachment_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_attachment_proto_goTypes = []any{
	(*GetAttachmentRequest)(nil),     // 0: proto.GetAttachmentRequest
	(*GetAttachmentResponse)(nil),    // 1: proto.GetAttachmentResponse
	(*ListAttachmentRequest)(nil),    // 2: proto.ListAttachmentRequest
	(*ListAttachmentResponse)(nil),   // 3: proto.ListAttachmentResponse
	(*UpdateVisibilityRequest)(nil),  // 4: proto.UpdateVisibilityRequest
	(*UpdateVisibilityResponse)(nil), // 5: proto.UpdateVisibilityResponse
	(*UpdateUsageRequest)(nil),       // 6: proto.UpdateUsageRequest
	(*UpdateUsageResponse)(nil),      // 7: proto.UpdateUsageResponse
	(*DeleteAttachmentRequest)(nil),  // 8: proto.DeleteAttachmentRequest
	(*DeleteAttachmentResponse)(nil), // 9: proto.DeleteAttachmentResponse
}
var file_attachment_proto_depIdxs = []int32{
	0, // 0: proto.AttachmentService.GetAttachment:input_type -> proto.GetAttachmentRequest
	2, // 1: proto.AttachmentService.ListAttachment:input_type -> proto.ListAttachmentRequest
	4, // 2: proto.AttachmentService.UpdateVisibility:input_type -> proto.UpdateVisibilityRequest
	6, // 3: proto.AttachmentService.UpdateUsage:input_type -> proto.UpdateUsageRequest
	8, // 4: proto.AttachmentService.DeleteAttachment:input_type -> proto.DeleteAttachmentRequest
	1, // 5: proto.AttachmentService.GetAttachment:output_type -> proto.GetAttachmentResponse
	3, // 6: proto.AttachmentService.ListAttachment:output_type -> proto.ListAttachmentResponse
	5, // 7: proto.AttachmentService.UpdateVisibility:output_type -> proto.UpdateVisibilityResponse
	7, // 8: proto.AttachmentService.UpdateUsage:output_type -> proto.UpdateUsageResponse
	9, // 9: proto.AttachmentService.DeleteAttachment:output_type -> proto.DeleteAttachmentResponse
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_attachment_proto_init() }
func file_attachment_proto_init() {
	if File_attachment_proto != nil {
		return
	}
	file_attachment_proto_msgTypes[0].OneofWrappers = []any{}
	file_attachment_proto_msgTypes[1].OneofWrappers = []any{}
	file_attachment_proto_msgTypes[2].OneofWrappers = []any{}
	file_attachment_proto_msgTypes[4].OneofWrappers = []any{}
	file_attachment_proto_msgTypes[8].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_attachment_proto_rawDesc), len(file_attachment_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_attachment_proto_goTypes,
		DependencyIndexes: file_attachment_proto_depIdxs,
		MessageInfos:      file_attachment_proto_msgTypes,
	}.Build()
	File_attachment_proto = out.File
	file_attachment_proto_goTypes = nil
	file_attachment_proto_depIdxs = nil
}
