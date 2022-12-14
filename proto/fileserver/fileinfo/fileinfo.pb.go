// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.4
// source: fileinfo.proto

package fileinfo

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

type BotConstants_BotFileType int32

const (
	BotConstants_BOT_FILE_TYPE_SINGLE    BotConstants_BotFileType = 1 //
	BotConstants_BOT_FILE_TYPE_MULTIPART BotConstants_BotFileType = 2 //
)

// Enum value maps for BotConstants_BotFileType.
var (
	BotConstants_BotFileType_name = map[int32]string{
		1: "BOT_FILE_TYPE_SINGLE",
		2: "BOT_FILE_TYPE_MULTIPART",
	}
	BotConstants_BotFileType_value = map[string]int32{
		"BOT_FILE_TYPE_SINGLE":    1,
		"BOT_FILE_TYPE_MULTIPART": 2,
	}
)

func (x BotConstants_BotFileType) Enum() *BotConstants_BotFileType {
	p := new(BotConstants_BotFileType)
	*p = x
	return p
}

func (x BotConstants_BotFileType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BotConstants_BotFileType) Descriptor() protoreflect.EnumDescriptor {
	return file_fileinfo_proto_enumTypes[0].Descriptor()
}

func (BotConstants_BotFileType) Type() protoreflect.EnumType {
	return &file_fileinfo_proto_enumTypes[0]
}

func (x BotConstants_BotFileType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *BotConstants_BotFileType) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = BotConstants_BotFileType(num)
	return nil
}

// Deprecated: Use BotConstants_BotFileType.Descriptor instead.
func (BotConstants_BotFileType) EnumDescriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{0, 0}
}

type BotConstants struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BotConstants) Reset() {
	*x = BotConstants{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BotConstants) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BotConstants) ProtoMessage() {}

func (x *BotConstants) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BotConstants.ProtoReflect.Descriptor instead.
func (*BotConstants) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{0}
}

type FileUploadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DownKey *string `protobuf:"bytes,1,opt,name=down_key,json=downKey" json:"down_key,omitempty"` //
}

func (x *FileUploadResponse) Reset() {
	*x = FileUploadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileUploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadResponse) ProtoMessage() {}

func (x *FileUploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadResponse.ProtoReflect.Descriptor instead.
func (*FileUploadResponse) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{1}
}

func (x *FileUploadResponse) GetDownKey() string {
	if x != nil && x.DownKey != nil {
		return *x.DownKey
	}
	return ""
}

type FileUploadBeginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileSize *uint64 `protobuf:"varint,1,opt,name=file_size,json=fileSize" json:"file_size,omitempty"` //
}

func (x *FileUploadBeginRequest) Reset() {
	*x = FileUploadBeginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileUploadBeginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadBeginRequest) ProtoMessage() {}

func (x *FileUploadBeginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadBeginRequest.ProtoReflect.Descriptor instead.
func (*FileUploadBeginRequest) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{2}
}

func (x *FileUploadBeginRequest) GetFileSize() uint64 {
	if x != nil && x.FileSize != nil {
		return *x.FileSize
	}
	return 0
}

type FileUploadBeginResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UploadCtx *string `protobuf:"bytes,1,opt,name=upload_ctx,json=uploadCtx" json:"upload_ctx,omitempty"`  //refer: UploadIdCtx
	BlockSize *uint32 `protobuf:"varint,2,opt,name=block_size,json=blockSize" json:"block_size,omitempty"` //
}

func (x *FileUploadBeginResponse) Reset() {
	*x = FileUploadBeginResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileUploadBeginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadBeginResponse) ProtoMessage() {}

func (x *FileUploadBeginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadBeginResponse.ProtoReflect.Descriptor instead.
func (*FileUploadBeginResponse) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{3}
}

func (x *FileUploadBeginResponse) GetUploadCtx() string {
	if x != nil && x.UploadCtx != nil {
		return *x.UploadCtx
	}
	return ""
}

func (x *FileUploadBeginResponse) GetBlockSize() uint32 {
	if x != nil && x.BlockSize != nil {
		return *x.BlockSize
	}
	return 0
}

type FileUploadPartResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FileUploadPartResponse) Reset() {
	*x = FileUploadPartResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileUploadPartResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadPartResponse) ProtoMessage() {}

func (x *FileUploadPartResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadPartResponse.ProtoReflect.Descriptor instead.
func (*FileUploadPartResponse) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{4}
}

type FileUploadEndRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UploadCtx *string `protobuf:"bytes,1,opt,name=upload_ctx,json=uploadCtx" json:"upload_ctx,omitempty"` //
	FileName  *string `protobuf:"bytes,2,opt,name=file_name,json=fileName" json:"file_name,omitempty"`    //
}

func (x *FileUploadEndRequest) Reset() {
	*x = FileUploadEndRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileUploadEndRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadEndRequest) ProtoMessage() {}

func (x *FileUploadEndRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadEndRequest.ProtoReflect.Descriptor instead.
func (*FileUploadEndRequest) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{5}
}

func (x *FileUploadEndRequest) GetUploadCtx() string {
	if x != nil && x.UploadCtx != nil {
		return *x.UploadCtx
	}
	return ""
}

func (x *FileUploadEndRequest) GetFileName() string {
	if x != nil && x.FileName != nil {
		return *x.FileName
	}
	return ""
}

type FileUploadEndResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DownKey *string `protobuf:"bytes,1,opt,name=down_key,json=downKey" json:"down_key,omitempty"` //
}

func (x *FileUploadEndResponse) Reset() {
	*x = FileUploadEndResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileUploadEndResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileUploadEndResponse) ProtoMessage() {}

func (x *FileUploadEndResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileUploadEndResponse.ProtoReflect.Descriptor instead.
func (*FileUploadEndResponse) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{6}
}

func (x *FileUploadEndResponse) GetDownKey() string {
	if x != nil && x.DownKey != nil {
		return *x.DownKey
	}
	return ""
}

type UploadIdCtx struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileSize  *uint64 `protobuf:"varint,1,opt,name=file_size,json=fileSize" json:"file_size,omitempty"`    //
	UploadId  *string `protobuf:"bytes,2,opt,name=upload_id,json=uploadId" json:"upload_id,omitempty"`     //
	BlockSize *uint32 `protobuf:"varint,3,opt,name=block_size,json=blockSize" json:"block_size,omitempty"` //
	FileKey   *string `protobuf:"bytes,4,opt,name=file_key,json=fileKey" json:"file_key,omitempty"`        //
	BotHash   *uint32 `protobuf:"varint,5,opt,name=bot_hash,json=botHash" json:"bot_hash,omitempty"`       //
}

func (x *UploadIdCtx) Reset() {
	*x = UploadIdCtx{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadIdCtx) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadIdCtx) ProtoMessage() {}

func (x *UploadIdCtx) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadIdCtx.ProtoReflect.Descriptor instead.
func (*UploadIdCtx) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{7}
}

func (x *UploadIdCtx) GetFileSize() uint64 {
	if x != nil && x.FileSize != nil {
		return *x.FileSize
	}
	return 0
}

func (x *UploadIdCtx) GetUploadId() string {
	if x != nil && x.UploadId != nil {
		return *x.UploadId
	}
	return ""
}

func (x *UploadIdCtx) GetBlockSize() uint32 {
	if x != nil && x.BlockSize != nil {
		return *x.BlockSize
	}
	return 0
}

func (x *UploadIdCtx) GetFileKey() string {
	if x != nil && x.FileKey != nil {
		return *x.FileKey
	}
	return ""
}

func (x *UploadIdCtx) GetBotHash() uint32 {
	if x != nil && x.BotHash != nil {
		return *x.BotHash
	}
	return 0
}

type GetFileMetaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DownKey []string `protobuf:"bytes,1,rep,name=down_key,json=downKey" json:"down_key,omitempty"` //
}

func (x *GetFileMetaRequest) Reset() {
	*x = GetFileMetaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileMetaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileMetaRequest) ProtoMessage() {}

func (x *GetFileMetaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileMetaRequest.ProtoReflect.Descriptor instead.
func (*GetFileMetaRequest) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{8}
}

func (x *GetFileMetaRequest) GetDownKey() []string {
	if x != nil {
		return x.DownKey
	}
	return nil
}

type FileItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName   *string `protobuf:"bytes,1,opt,name=file_name,json=fileName" json:"file_name,omitempty"`        //
	Hash       *string `protobuf:"bytes,2,opt,name=hash" json:"hash,omitempty"`                                //
	FileSize   *uint64 `protobuf:"varint,3,opt,name=file_size,json=fileSize" json:"file_size,omitempty"`       //
	CreateTime *uint64 `protobuf:"varint,4,opt,name=create_time,json=createTime" json:"create_time,omitempty"` //
	DownKey    *string `protobuf:"bytes,5,opt,name=down_key,json=downKey" json:"down_key,omitempty"`           //
	Exist      *bool   `protobuf:"varint,6,opt,name=exist" json:"exist,omitempty"`                             //
}

func (x *FileItem) Reset() {
	*x = FileItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileItem) ProtoMessage() {}

func (x *FileItem) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileItem.ProtoReflect.Descriptor instead.
func (*FileItem) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{9}
}

func (x *FileItem) GetFileName() string {
	if x != nil && x.FileName != nil {
		return *x.FileName
	}
	return ""
}

func (x *FileItem) GetHash() string {
	if x != nil && x.Hash != nil {
		return *x.Hash
	}
	return ""
}

func (x *FileItem) GetFileSize() uint64 {
	if x != nil && x.FileSize != nil {
		return *x.FileSize
	}
	return 0
}

func (x *FileItem) GetCreateTime() uint64 {
	if x != nil && x.CreateTime != nil {
		return *x.CreateTime
	}
	return 0
}

func (x *FileItem) GetDownKey() string {
	if x != nil && x.DownKey != nil {
		return *x.DownKey
	}
	return ""
}

func (x *FileItem) GetExist() bool {
	if x != nil && x.Exist != nil {
		return *x.Exist
	}
	return false
}

type GetFileMetaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []*FileItem `protobuf:"bytes,1,rep,name=list" json:"list,omitempty"` //
}

func (x *GetFileMetaResponse) Reset() {
	*x = GetFileMetaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFileMetaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileMetaResponse) ProtoMessage() {}

func (x *GetFileMetaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileMetaResponse.ProtoReflect.Descriptor instead.
func (*GetFileMetaResponse) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{10}
}

func (x *GetFileMetaResponse) GetList() []*FileItem {
	if x != nil {
		return x.List
	}
	return nil
}

type BotFileExtra struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BotHash   *uint32 `protobuf:"varint,1,opt,name=bot_hash,json=botHash" json:"bot_hash,omitempty"`       //
	FileType  *int32  `protobuf:"varint,2,opt,name=file_type,json=fileType" json:"file_type,omitempty"`    //
	BlockSize *int64  `protobuf:"varint,3,opt,name=block_size,json=blockSize" json:"block_size,omitempty"` //
	FileSize  *int64  `protobuf:"varint,4,opt,name=file_size,json=fileSize" json:"file_size,omitempty"`    //
}

func (x *BotFileExtra) Reset() {
	*x = BotFileExtra{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BotFileExtra) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BotFileExtra) ProtoMessage() {}

func (x *BotFileExtra) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BotFileExtra.ProtoReflect.Descriptor instead.
func (*BotFileExtra) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{11}
}

func (x *BotFileExtra) GetBotHash() uint32 {
	if x != nil && x.BotHash != nil {
		return *x.BotHash
	}
	return 0
}

func (x *BotFileExtra) GetFileType() int32 {
	if x != nil && x.FileType != nil {
		return *x.FileType
	}
	return 0
}

func (x *BotFileExtra) GetBlockSize() int64 {
	if x != nil && x.BlockSize != nil {
		return *x.BlockSize
	}
	return 0
}

func (x *BotFileExtra) GetFileSize() int64 {
	if x != nil && x.FileSize != nil {
		return *x.FileSize
	}
	return 0
}

type BotUploadContext struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileSize  *int64   `protobuf:"varint,1,opt,name=file_size,json=fileSize" json:"file_size,omitempty"`    //
	BlockSize *int64   `protobuf:"varint,2,opt,name=block_size,json=blockSize" json:"block_size,omitempty"` //
	Blocks    []string `protobuf:"bytes,3,rep,name=blocks" json:"blocks,omitempty"`                         //
}

func (x *BotUploadContext) Reset() {
	*x = BotUploadContext{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BotUploadContext) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BotUploadContext) ProtoMessage() {}

func (x *BotUploadContext) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BotUploadContext.ProtoReflect.Descriptor instead.
func (*BotUploadContext) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{12}
}

func (x *BotUploadContext) GetFileSize() int64 {
	if x != nil && x.FileSize != nil {
		return *x.FileSize
	}
	return 0
}

func (x *BotUploadContext) GetBlockSize() int64 {
	if x != nil && x.BlockSize != nil {
		return *x.BlockSize
	}
	return 0
}

func (x *BotUploadContext) GetBlocks() []string {
	if x != nil {
		return x.Blocks
	}
	return nil
}

type PartPair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PartId   *int32  `protobuf:"varint,1,opt,name=part_id,json=partId" json:"part_id,omitempty"`       //
	PartKey  *string `protobuf:"bytes,2,opt,name=part_key,json=partKey" json:"part_key,omitempty"`     //
	PartSize *int64  `protobuf:"varint,3,opt,name=part_size,json=partSize" json:"part_size,omitempty"` //
	Md5Value *string `protobuf:"bytes,4,opt,name=md5_value,json=md5Value" json:"md5_value,omitempty"`  //
}

func (x *PartPair) Reset() {
	*x = PartPair{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fileinfo_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartPair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartPair) ProtoMessage() {}

func (x *PartPair) ProtoReflect() protoreflect.Message {
	mi := &file_fileinfo_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartPair.ProtoReflect.Descriptor instead.
func (*PartPair) Descriptor() ([]byte, []int) {
	return file_fileinfo_proto_rawDescGZIP(), []int{13}
}

func (x *PartPair) GetPartId() int32 {
	if x != nil && x.PartId != nil {
		return *x.PartId
	}
	return 0
}

func (x *PartPair) GetPartKey() string {
	if x != nil && x.PartKey != nil {
		return *x.PartKey
	}
	return ""
}

func (x *PartPair) GetPartSize() int64 {
	if x != nil && x.PartSize != nil {
		return *x.PartSize
	}
	return 0
}

func (x *PartPair) GetMd5Value() string {
	if x != nil && x.Md5Value != nil {
		return *x.Md5Value
	}
	return ""
}

var File_fileinfo_proto protoreflect.FileDescriptor

var file_fileinfo_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x66, 0x69, 0x6c, 0x65, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x67, 0x61, 0x6d, 0x65, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x54, 0x0a, 0x0c, 0x42, 0x6f,
	0x74, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x73, 0x22, 0x44, 0x0a, 0x0b, 0x42, 0x6f,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x14, 0x42, 0x4f, 0x54,
	0x5f, 0x46, 0x49, 0x4c, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x49, 0x4e, 0x47, 0x4c,
	0x45, 0x10, 0x01, 0x12, 0x1b, 0x0a, 0x17, 0x42, 0x4f, 0x54, 0x5f, 0x46, 0x49, 0x4c, 0x45, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x55, 0x4c, 0x54, 0x49, 0x50, 0x41, 0x52, 0x54, 0x10, 0x02,
	0x22, 0x2f, 0x0a, 0x12, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x6f, 0x77, 0x6e, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x6f, 0x77, 0x6e, 0x4b, 0x65,
	0x79, 0x22, 0x35, 0x0a, 0x16, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x42,
	0x65, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08,
	0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x57, 0x0a, 0x17, 0x46, 0x69, 0x6c, 0x65,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x65, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x63, 0x74,
	0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x43,
	0x74, 0x78, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x69, 0x7a,
	0x65, 0x22, 0x18, 0x0a, 0x16, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50,
	0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x52, 0x0a, 0x14, 0x46,
	0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x45, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x63, 0x74,
	0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x43,
	0x74, 0x78, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x32, 0x0a, 0x15, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x45, 0x6e, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x6f, 0x77, 0x6e,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x6f, 0x77, 0x6e,
	0x4b, 0x65, 0x79, 0x22, 0x9c, 0x01, 0x0a, 0x0b, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x64,
	0x43, 0x74, 0x78, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x19, 0x0a, 0x08,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x66, 0x69, 0x6c, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x62, 0x6f, 0x74, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x62, 0x6f, 0x74, 0x48, 0x61,
	0x73, 0x68, 0x22, 0x2f, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x4d, 0x65, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x6f, 0x77, 0x6e,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x64, 0x6f, 0x77, 0x6e,
	0x4b, 0x65, 0x79, 0x22, 0xaa, 0x01, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x74, 0x65, 0x6d,
	0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73,
	0x68, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1f,
	0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x64, 0x6f, 0x77, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x64, 0x6f, 0x77, 0x6e, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x78,
	0x69, 0x73, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x65, 0x78, 0x69, 0x73, 0x74,
	0x22, 0x3d, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x69, 0x6e, 0x66, 0x6f,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x22,
	0x82, 0x01, 0x0a, 0x0c, 0x42, 0x6f, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x45, 0x78, 0x74, 0x72, 0x61,
	0x12, 0x19, 0x0a, 0x08, 0x62, 0x6f, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x62, 0x6f, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1b, 0x0a, 0x09, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x66, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65,
	0x53, 0x69, 0x7a, 0x65, 0x22, 0x66, 0x0a, 0x10, 0x42, 0x6f, 0x74, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x53, 0x69, 0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x22, 0x78, 0x0a, 0x08,
	0x50, 0x61, 0x72, 0x74, 0x50, 0x61, 0x69, 0x72, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x61, 0x72, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x70, 0x61, 0x72, 0x74, 0x49,
	0x64, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x61, 0x72, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x72, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x1b, 0x0a, 0x09,
	0x70, 0x61, 0x72, 0x74, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x08, 0x70, 0x61, 0x72, 0x74, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x64, 0x35,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x64,
	0x35, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x15, 0x5a, 0x13, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x69, 0x6e, 0x66, 0x6f,
}

var (
	file_fileinfo_proto_rawDescOnce sync.Once
	file_fileinfo_proto_rawDescData = file_fileinfo_proto_rawDesc
)

func file_fileinfo_proto_rawDescGZIP() []byte {
	file_fileinfo_proto_rawDescOnce.Do(func() {
		file_fileinfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_fileinfo_proto_rawDescData)
	})
	return file_fileinfo_proto_rawDescData
}

var file_fileinfo_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_fileinfo_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_fileinfo_proto_goTypes = []interface{}{
	(BotConstants_BotFileType)(0),   // 0: gameinfo.BotConstants.BotFileType
	(*BotConstants)(nil),            // 1: gameinfo.BotConstants
	(*FileUploadResponse)(nil),      // 2: gameinfo.FileUploadResponse
	(*FileUploadBeginRequest)(nil),  // 3: gameinfo.FileUploadBeginRequest
	(*FileUploadBeginResponse)(nil), // 4: gameinfo.FileUploadBeginResponse
	(*FileUploadPartResponse)(nil),  // 5: gameinfo.FileUploadPartResponse
	(*FileUploadEndRequest)(nil),    // 6: gameinfo.FileUploadEndRequest
	(*FileUploadEndResponse)(nil),   // 7: gameinfo.FileUploadEndResponse
	(*UploadIdCtx)(nil),             // 8: gameinfo.UploadIdCtx
	(*GetFileMetaRequest)(nil),      // 9: gameinfo.GetFileMetaRequest
	(*FileItem)(nil),                // 10: gameinfo.FileItem
	(*GetFileMetaResponse)(nil),     // 11: gameinfo.GetFileMetaResponse
	(*BotFileExtra)(nil),            // 12: gameinfo.BotFileExtra
	(*BotUploadContext)(nil),        // 13: gameinfo.BotUploadContext
	(*PartPair)(nil),                // 14: gameinfo.PartPair
}
var file_fileinfo_proto_depIdxs = []int32{
	10, // 0: gameinfo.GetFileMetaResponse.list:type_name -> gameinfo.FileItem
	1,  // [1:1] is the sub-list for method output_type
	1,  // [1:1] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_fileinfo_proto_init() }
func file_fileinfo_proto_init() {
	if File_fileinfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fileinfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BotConstants); i {
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
		file_fileinfo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileUploadResponse); i {
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
		file_fileinfo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileUploadBeginRequest); i {
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
		file_fileinfo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileUploadBeginResponse); i {
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
		file_fileinfo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileUploadPartResponse); i {
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
		file_fileinfo_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileUploadEndRequest); i {
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
		file_fileinfo_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileUploadEndResponse); i {
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
		file_fileinfo_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadIdCtx); i {
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
		file_fileinfo_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileMetaRequest); i {
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
		file_fileinfo_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileItem); i {
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
		file_fileinfo_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFileMetaResponse); i {
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
		file_fileinfo_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BotFileExtra); i {
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
		file_fileinfo_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BotUploadContext); i {
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
		file_fileinfo_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PartPair); i {
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
			RawDescriptor: file_fileinfo_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fileinfo_proto_goTypes,
		DependencyIndexes: file_fileinfo_proto_depIdxs,
		EnumInfos:         file_fileinfo_proto_enumTypes,
		MessageInfos:      file_fileinfo_proto_msgTypes,
	}.Build()
	File_fileinfo_proto = out.File
	file_fileinfo_proto_rawDesc = nil
	file_fileinfo_proto_goTypes = nil
	file_fileinfo_proto_depIdxs = nil
}
