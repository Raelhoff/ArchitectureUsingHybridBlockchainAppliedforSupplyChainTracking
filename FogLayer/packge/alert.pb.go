// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/alert.proto

package packge

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type AlertRequest struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	IdEdge               int32    `protobuf:"varint,2,opt,name=idEdge,proto3" json:"idEdge,omitempty"`
	IdNodo               int32    `protobuf:"varint,3,opt,name=idNodo,proto3" json:"idNodo,omitempty"`
	Hash                 string   `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
	Type                 int32    `protobuf:"varint,5,opt,name=type,proto3" json:"type,omitempty"`
	Date                 string   `protobuf:"bytes,6,opt,name=date,proto3" json:"date,omitempty"`
	Temperatura          string   `protobuf:"bytes,7,opt,name=temperatura,proto3" json:"temperatura,omitempty"`
	Umidade              string   `protobuf:"bytes,8,opt,name=umidade,proto3" json:"umidade,omitempty"`
	Rele                 string   `protobuf:"bytes,9,opt,name=rele,proto3" json:"rele,omitempty"`
	Description          string   `protobuf:"bytes,10,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AlertRequest) Reset()         { *m = AlertRequest{} }
func (m *AlertRequest) String() string { return proto.CompactTextString(m) }
func (*AlertRequest) ProtoMessage()    {}
func (*AlertRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_bd2c30df7992c067, []int{0}
}

func (m *AlertRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AlertRequest.Unmarshal(m, b)
}
func (m *AlertRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AlertRequest.Marshal(b, m, deterministic)
}
func (m *AlertRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AlertRequest.Merge(m, src)
}
func (m *AlertRequest) XXX_Size() int {
	return xxx_messageInfo_AlertRequest.Size(m)
}
func (m *AlertRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AlertRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AlertRequest proto.InternalMessageInfo

func (m *AlertRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *AlertRequest) GetIdEdge() int32 {
	if m != nil {
		return m.IdEdge
	}
	return 0
}

func (m *AlertRequest) GetIdNodo() int32 {
	if m != nil {
		return m.IdNodo
	}
	return 0
}

func (m *AlertRequest) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *AlertRequest) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *AlertRequest) GetDate() string {
	if m != nil {
		return m.Date
	}
	return ""
}

func (m *AlertRequest) GetTemperatura() string {
	if m != nil {
		return m.Temperatura
	}
	return ""
}

func (m *AlertRequest) GetUmidade() string {
	if m != nil {
		return m.Umidade
	}
	return ""
}

func (m *AlertRequest) GetRele() string {
	if m != nil {
		return m.Rele
	}
	return ""
}

func (m *AlertRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type AlertResponse struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	Confirmation         bool     `protobuf:"varint,2,opt,name=confirmation,proto3" json:"confirmation,omitempty"`
	Type                 string   `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AlertResponse) Reset()         { *m = AlertResponse{} }
func (m *AlertResponse) String() string { return proto.CompactTextString(m) }
func (*AlertResponse) ProtoMessage()    {}
func (*AlertResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_bd2c30df7992c067, []int{1}
}

func (m *AlertResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AlertResponse.Unmarshal(m, b)
}
func (m *AlertResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AlertResponse.Marshal(b, m, deterministic)
}
func (m *AlertResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AlertResponse.Merge(m, src)
}
func (m *AlertResponse) XXX_Size() int {
	return xxx_messageInfo_AlertResponse.Size(m)
}
func (m *AlertResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AlertResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AlertResponse proto.InternalMessageInfo

func (m *AlertResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *AlertResponse) GetConfirmation() bool {
	if m != nil {
		return m.Confirmation
	}
	return false
}

func (m *AlertResponse) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func init() {
	proto.RegisterType((*AlertRequest)(nil), "alert.AlertRequest")
	proto.RegisterType((*AlertResponse)(nil), "alert.AlertResponse")
}

func init() {
	proto.RegisterFile("proto/alert.proto", fileDescriptor_bd2c30df7992c067)
}

var fileDescriptor_bd2c30df7992c067 = []byte{
	// 286 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0x31, 0x4f, 0xf3, 0x30,
	0x10, 0x86, 0xbf, 0xb4, 0x4d, 0x9a, 0xdc, 0x57, 0x10, 0x18, 0x84, 0x2c, 0xa6, 0x28, 0x53, 0xa7,
	0x54, 0x82, 0x85, 0xb5, 0x48, 0xac, 0x0c, 0x61, 0x82, 0xcd, 0xc4, 0x47, 0x6a, 0xd1, 0xc4, 0xc6,
	0x76, 0x06, 0xfe, 0x3b, 0x03, 0xf2, 0xb9, 0x45, 0xe9, 0xf6, 0xbc, 0x8f, 0xad, 0x37, 0x39, 0x1f,
	0x5c, 0x1a, 0xab, 0xbd, 0xde, 0x88, 0x3d, 0x5a, 0x5f, 0x13, 0xb3, 0x94, 0x42, 0xf5, 0x93, 0xc0,
	0x6a, 0x1b, 0xa8, 0xc1, 0xaf, 0x11, 0x9d, 0x67, 0xe7, 0x30, 0x53, 0x92, 0x27, 0x65, 0xb2, 0x4e,
	0x9b, 0x99, 0x92, 0xec, 0x06, 0x32, 0x25, 0x9f, 0x64, 0x87, 0x7c, 0x46, 0xee, 0x90, 0xa2, 0x7f,
	0xd6, 0x52, 0xf3, 0xf9, 0xd1, 0x87, 0xc4, 0x18, 0x2c, 0x76, 0xc2, 0xed, 0xf8, 0xa2, 0x4c, 0xd6,
	0x45, 0x43, 0x1c, 0x9c, 0xff, 0x36, 0xc8, 0x53, 0xba, 0x49, 0x1c, 0x9c, 0x14, 0x1e, 0x79, 0x16,
	0xef, 0x05, 0x66, 0x25, 0xfc, 0xf7, 0xd8, 0x1b, 0xb4, 0xc2, 0x8f, 0x56, 0xf0, 0x25, 0x1d, 0x4d,
	0x15, 0xe3, 0xb0, 0x1c, 0x7b, 0x25, 0x85, 0x44, 0x9e, 0xd3, 0xe9, 0x31, 0x86, 0x3e, 0x8b, 0x7b,
	0xe4, 0x45, 0xec, 0x0b, 0x1c, 0xfa, 0x24, 0xba, 0xd6, 0x2a, 0xe3, 0x95, 0x1e, 0x38, 0xc4, 0xbe,
	0x89, 0xaa, 0x5e, 0xe1, 0xec, 0x30, 0xbd, 0x33, 0x7a, 0x70, 0xc8, 0x2e, 0x60, 0xde, 0xbb, 0x8e,
	0xe6, 0x2f, 0x9a, 0x80, 0xac, 0x82, 0x55, 0xab, 0x87, 0x0f, 0x65, 0x7b, 0x41, 0x2d, 0xe1, 0x19,
	0xf2, 0xe6, 0xc4, 0xfd, 0x0d, 0x38, 0x8f, 0x1f, 0x0f, 0x7c, 0xb7, 0x85, 0x94, 0xaa, 0xd9, 0x03,
	0x14, 0x2f, 0x38, 0xc8, 0x18, 0xae, 0xea, 0xb8, 0x84, 0xe9, 0x9b, 0xdf, 0x5e, 0x9f, 0xca, 0xf8,
	0x2b, 0xd5, 0xbf, 0x47, 0x78, 0xcb, 0xeb, 0x8d, 0x11, 0xed, 0x67, 0x87, 0xef, 0x19, 0xad, 0xed,
	0xfe, 0x37, 0x00, 0x00, 0xff, 0xff, 0x4d, 0x6b, 0x88, 0xde, 0xcb, 0x01, 0x00, 0x00,
}
