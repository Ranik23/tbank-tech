// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.12.4
// source: scrapper/api/proto/scrapper.proto

package gen

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// Запрос для регистрации пользователя
type RegisterUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TgUserId      int64                  `protobuf:"varint,1,opt,name=tg_user_id,json=tgUserId,proto3" json:"tg_user_id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"` // Имя пользователя
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterUserRequest) Reset() {
	*x = RegisterUserRequest{}
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterUserRequest) ProtoMessage() {}

func (x *RegisterUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterUserRequest.ProtoReflect.Descriptor instead.
func (*RegisterUserRequest) Descriptor() ([]byte, []int) {
	return file_scrapper_api_proto_scrapper_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterUserRequest) GetTgUserId() int64 {
	if x != nil {
		return x.TgUserId
	}
	return 0
}

func (x *RegisterUserRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Ответ на регистрацию пользователя
type RegisterUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterUserResponse) Reset() {
	*x = RegisterUserResponse{}
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterUserResponse) ProtoMessage() {}

func (x *RegisterUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterUserResponse.ProtoReflect.Descriptor instead.
func (*RegisterUserResponse) Descriptor() ([]byte, []int) {
	return file_scrapper_api_proto_scrapper_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterUserResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Запрос для удаления пользователя
type DeleteUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TgUserId      int64                  `protobuf:"varint,1,opt,name=tg_user_id,json=tgUserId,proto3" json:"tg_user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteUserRequest) Reset() {
	*x = DeleteUserRequest{}
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserRequest) ProtoMessage() {}

func (x *DeleteUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserRequest) Descriptor() ([]byte, []int) {
	return file_scrapper_api_proto_scrapper_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteUserRequest) GetTgUserId() int64 {
	if x != nil {
		return x.TgUserId
	}
	return 0
}

// Ответ на удаление пользователя
type DeleteUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteUserResponse) Reset() {
	*x = DeleteUserResponse{}
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserResponse) ProtoMessage() {}

func (x *DeleteUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserResponse.ProtoReflect.Descriptor instead.
func (*DeleteUserResponse) Descriptor() ([]byte, []int) {
	return file_scrapper_api_proto_scrapper_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteUserResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Запрос для получения ссылок
type GetLinksRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TgUserId      int64                  `protobuf:"varint,1,opt,name=tg_user_id,json=tgUserId,proto3" json:"tg_user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetLinksRequest) Reset() {
	*x = GetLinksRequest{}
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetLinksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLinksRequest) ProtoMessage() {}

func (x *GetLinksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLinksRequest.ProtoReflect.Descriptor instead.
func (*GetLinksRequest) Descriptor() ([]byte, []int) {
	return file_scrapper_api_proto_scrapper_proto_rawDescGZIP(), []int{4}
}

func (x *GetLinksRequest) GetTgUserId() int64 {
	if x != nil {
		return x.TgUserId
	}
	return 0
}

// Запрос для добавления ссылки
type AddLinkRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TgUserId      int64                  `protobuf:"varint,1,opt,name=tg_user_id,json=tgUserId,proto3" json:"tg_user_id,omitempty"`
	Url           string                 `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddLinkRequest) Reset() {
	*x = AddLinkRequest{}
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddLinkRequest) ProtoMessage() {}

func (x *AddLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddLinkRequest.ProtoReflect.Descriptor instead.
func (*AddLinkRequest) Descriptor() ([]byte, []int) {
	return file_scrapper_api_proto_scrapper_proto_rawDescGZIP(), []int{5}
}

func (x *AddLinkRequest) GetTgUserId() int64 {
	if x != nil {
		return x.TgUserId
	}
	return 0
}

func (x *AddLinkRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

// Запрос для удаления ссылки
type RemoveLinkRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TgUserId      int64                  `protobuf:"varint,1,opt,name=tg_user_id,json=tgUserId,proto3" json:"tg_user_id,omitempty"`
	Url           string                 `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveLinkRequest) Reset() {
	*x = RemoveLinkRequest{}
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveLinkRequest) ProtoMessage() {}

func (x *RemoveLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveLinkRequest.ProtoReflect.Descriptor instead.
func (*RemoveLinkRequest) Descriptor() ([]byte, []int) {
	return file_scrapper_api_proto_scrapper_proto_rawDescGZIP(), []int{6}
}

func (x *RemoveLinkRequest) GetTgUserId() int64 {
	if x != nil {
		return x.TgUserId
	}
	return 0
}

func (x *RemoveLinkRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

// Ответ на ссылку
type LinkResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Url           string                 `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LinkResponse) Reset() {
	*x = LinkResponse{}
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LinkResponse) ProtoMessage() {}

func (x *LinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LinkResponse.ProtoReflect.Descriptor instead.
func (*LinkResponse) Descriptor() ([]byte, []int) {
	return file_scrapper_api_proto_scrapper_proto_rawDescGZIP(), []int{7}
}

func (x *LinkResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LinkResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

// Ответ на получение списка ссылок
type ListLinksResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Links         []*LinkResponse        `protobuf:"bytes,1,rep,name=links,proto3" json:"links,omitempty"`
	Size          int32                  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListLinksResponse) Reset() {
	*x = ListLinksResponse{}
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListLinksResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLinksResponse) ProtoMessage() {}

func (x *ListLinksResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLinksResponse.ProtoReflect.Descriptor instead.
func (*ListLinksResponse) Descriptor() ([]byte, []int) {
	return file_scrapper_api_proto_scrapper_proto_rawDescGZIP(), []int{8}
}

func (x *ListLinksResponse) GetLinks() []*LinkResponse {
	if x != nil {
		return x.Links
	}
	return nil
}

func (x *ListLinksResponse) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

// Ответ на ошибку API
type ApiErrorResponse struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	Description      string                 `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	Code             string                 `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	ExceptionName    string                 `protobuf:"bytes,3,opt,name=exception_name,json=exceptionName,proto3" json:"exception_name,omitempty"`
	ExceptionMessage string                 `protobuf:"bytes,4,opt,name=exception_message,json=exceptionMessage,proto3" json:"exception_message,omitempty"`
	Stacktrace       []string               `protobuf:"bytes,5,rep,name=stacktrace,proto3" json:"stacktrace,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *ApiErrorResponse) Reset() {
	*x = ApiErrorResponse{}
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ApiErrorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiErrorResponse) ProtoMessage() {}

func (x *ApiErrorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_scrapper_api_proto_scrapper_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiErrorResponse.ProtoReflect.Descriptor instead.
func (*ApiErrorResponse) Descriptor() ([]byte, []int) {
	return file_scrapper_api_proto_scrapper_proto_rawDescGZIP(), []int{9}
}

func (x *ApiErrorResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ApiErrorResponse) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *ApiErrorResponse) GetExceptionName() string {
	if x != nil {
		return x.ExceptionName
	}
	return ""
}

func (x *ApiErrorResponse) GetExceptionMessage() string {
	if x != nil {
		return x.ExceptionMessage
	}
	return ""
}

func (x *ApiErrorResponse) GetStacktrace() []string {
	if x != nil {
		return x.Stacktrace
	}
	return nil
}

var File_scrapper_api_proto_scrapper_proto protoreflect.FileDescriptor

var file_scrapper_api_proto_scrapper_proto_rawDesc = string([]byte{
	0x0a, 0x21, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x47, 0x0a, 0x13, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1c, 0x0a, 0x0a, 0x74, 0x67, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x74, 0x67, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x30, 0x0a, 0x14, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x31, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x0a, 0x74,
	0x67, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x08, 0x74, 0x67, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x2e, 0x0a, 0x12, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x2f, 0x0a, 0x0f, 0x47, 0x65, 0x74,
	0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x0a,
	0x74, 0x67, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x74, 0x67, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x40, 0x0a, 0x0e, 0x41, 0x64,
	0x64, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x0a,
	0x74, 0x67, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x74, 0x67, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x43, 0x0a, 0x11,
	0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1c, 0x0a, 0x0a, 0x74, 0x67, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x74, 0x67, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x22, 0x30, 0x0a, 0x0c, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6c, 0x22, 0x55, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x6c, 0x69, 0x6e, 0x6b,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52,
	0x05, 0x6c, 0x69, 0x6e, 0x6b, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x22, 0xbc, 0x01, 0x0a, 0x10, 0x41,
	0x70, 0x69, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65,
	0x78, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2b, 0x0a, 0x11,
	0x65, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x65, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x74, 0x61,
	0x63, 0x6b, 0x74, 0x72, 0x61, 0x63, 0x65, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x73,
	0x74, 0x61, 0x63, 0x6b, 0x74, 0x72, 0x61, 0x63, 0x65, 0x32, 0xd7, 0x03, 0x0a, 0x08, 0x53, 0x63,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x6c, 0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1d, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70, 0x65,
	0x72, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x22, 0x15, 0x2f,
	0x74, 0x67, 0x2d, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x7b, 0x74, 0x67, 0x5f, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x7d, 0x12, 0x66, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x12, 0x1b, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x17, 0x2a, 0x15, 0x2f, 0x74, 0x67, 0x2d, 0x63, 0x68, 0x61, 0x74, 0x2f,
	0x7b, 0x74, 0x67, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0x52, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x12, 0x19, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x0e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x08, 0x12, 0x06, 0x2f, 0x6c, 0x69, 0x6e, 0x6b, 0x73,
	0x12, 0x4e, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x18, 0x2e, 0x73, 0x63,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x41, 0x64, 0x64, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x3a, 0x01, 0x2a, 0x22, 0x06, 0x2f, 0x6c, 0x69, 0x6e, 0x6b, 0x73,
	0x12, 0x51, 0x0a, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x1b,
	0x2e, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x73, 0x63,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x0e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x08, 0x2a, 0x06, 0x2f, 0x6c, 0x69,
	0x6e, 0x6b, 0x73, 0x42, 0x1c, 0x5a, 0x1a, 0x73, 0x63, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x3b, 0x67, 0x65,
	0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_scrapper_api_proto_scrapper_proto_rawDescOnce sync.Once
	file_scrapper_api_proto_scrapper_proto_rawDescData []byte
)

func file_scrapper_api_proto_scrapper_proto_rawDescGZIP() []byte {
	file_scrapper_api_proto_scrapper_proto_rawDescOnce.Do(func() {
		file_scrapper_api_proto_scrapper_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_scrapper_api_proto_scrapper_proto_rawDesc), len(file_scrapper_api_proto_scrapper_proto_rawDesc)))
	})
	return file_scrapper_api_proto_scrapper_proto_rawDescData
}

var file_scrapper_api_proto_scrapper_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_scrapper_api_proto_scrapper_proto_goTypes = []any{
	(*RegisterUserRequest)(nil),  // 0: scrapper.RegisterUserRequest
	(*RegisterUserResponse)(nil), // 1: scrapper.RegisterUserResponse
	(*DeleteUserRequest)(nil),    // 2: scrapper.DeleteUserRequest
	(*DeleteUserResponse)(nil),   // 3: scrapper.DeleteUserResponse
	(*GetLinksRequest)(nil),      // 4: scrapper.GetLinksRequest
	(*AddLinkRequest)(nil),       // 5: scrapper.AddLinkRequest
	(*RemoveLinkRequest)(nil),    // 6: scrapper.RemoveLinkRequest
	(*LinkResponse)(nil),         // 7: scrapper.LinkResponse
	(*ListLinksResponse)(nil),    // 8: scrapper.ListLinksResponse
	(*ApiErrorResponse)(nil),     // 9: scrapper.ApiErrorResponse
}
var file_scrapper_api_proto_scrapper_proto_depIdxs = []int32{
	7, // 0: scrapper.ListLinksResponse.links:type_name -> scrapper.LinkResponse
	0, // 1: scrapper.Scrapper.RegisterUser:input_type -> scrapper.RegisterUserRequest
	2, // 2: scrapper.Scrapper.DeleteUser:input_type -> scrapper.DeleteUserRequest
	4, // 3: scrapper.Scrapper.GetLinks:input_type -> scrapper.GetLinksRequest
	5, // 4: scrapper.Scrapper.AddLink:input_type -> scrapper.AddLinkRequest
	6, // 5: scrapper.Scrapper.RemoveLink:input_type -> scrapper.RemoveLinkRequest
	1, // 6: scrapper.Scrapper.RegisterUser:output_type -> scrapper.RegisterUserResponse
	3, // 7: scrapper.Scrapper.DeleteUser:output_type -> scrapper.DeleteUserResponse
	8, // 8: scrapper.Scrapper.GetLinks:output_type -> scrapper.ListLinksResponse
	7, // 9: scrapper.Scrapper.AddLink:output_type -> scrapper.LinkResponse
	7, // 10: scrapper.Scrapper.RemoveLink:output_type -> scrapper.LinkResponse
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_scrapper_api_proto_scrapper_proto_init() }
func file_scrapper_api_proto_scrapper_proto_init() {
	if File_scrapper_api_proto_scrapper_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_scrapper_api_proto_scrapper_proto_rawDesc), len(file_scrapper_api_proto_scrapper_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_scrapper_api_proto_scrapper_proto_goTypes,
		DependencyIndexes: file_scrapper_api_proto_scrapper_proto_depIdxs,
		MessageInfos:      file_scrapper_api_proto_scrapper_proto_msgTypes,
	}.Build()
	File_scrapper_api_proto_scrapper_proto = out.File
	file_scrapper_api_proto_scrapper_proto_goTypes = nil
	file_scrapper_api_proto_scrapper_proto_depIdxs = nil
}
