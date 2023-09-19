// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/keepalive.proto

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

type KeepAliveRequest struct {
	IdEdge               int32    `protobuf:"varint,1,opt,name=idEdge,proto3" json:"idEdge,omitempty"`
	IdNodo               int32    `protobuf:"varint,2,opt,name=idNodo,proto3" json:"idNodo,omitempty"`
	Type                 int32    `protobuf:"varint,3,opt,name=type,proto3" json:"type,omitempty"`
	Date                 string   `protobuf:"bytes,4,opt,name=date,proto3" json:"date,omitempty"`
	Msg                  string   `protobuf:"bytes,5,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeepAliveRequest) Reset()         { *m = KeepAliveRequest{} }
func (m *KeepAliveRequest) String() string { return proto.CompactTextString(m) }
func (*KeepAliveRequest) ProtoMessage()    {}
func (*KeepAliveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3e815100fbc7a5f2, []int{0}
}

func (m *KeepAliveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeepAliveRequest.Unmarshal(m, b)
}
func (m *KeepAliveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeepAliveRequest.Marshal(b, m, deterministic)
}
func (m *KeepAliveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeepAliveRequest.Merge(m, src)
}
func (m *KeepAliveRequest) XXX_Size() int {
	return xxx_messageInfo_KeepAliveRequest.Size(m)
}
func (m *KeepAliveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_KeepAliveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_KeepAliveRequest proto.InternalMessageInfo

func (m *KeepAliveRequest) GetIdEdge() int32 {
	if m != nil {
		return m.IdEdge
	}
	return 0
}

func (m *KeepAliveRequest) GetIdNodo() int32 {
	if m != nil {
		return m.IdNodo
	}
	return 0
}

func (m *KeepAliveRequest) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *KeepAliveRequest) GetDate() string {
	if m != nil {
		return m.Date
	}
	return ""
}

func (m *KeepAliveRequest) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type KeepAliveResponse struct {
	IdEdge               int32    `protobuf:"varint,1,opt,name=idEdge,proto3" json:"idEdge,omitempty"`
	IdNodo               int32    `protobuf:"varint,2,opt,name=idNodo,proto3" json:"idNodo,omitempty"`
	Type                 int32    `protobuf:"varint,3,opt,name=type,proto3" json:"type,omitempty"`
	Date                 string   `protobuf:"bytes,4,opt,name=date,proto3" json:"date,omitempty"`
	Msg                  string   `protobuf:"bytes,5,opt,name=msg,proto3" json:"msg,omitempty"`
	Confirmation         bool     `protobuf:"varint,6,opt,name=confirmation,proto3" json:"confirmation,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeepAliveResponse) Reset()         { *m = KeepAliveResponse{} }
func (m *KeepAliveResponse) String() string { return proto.CompactTextString(m) }
func (*KeepAliveResponse) ProtoMessage()    {}
func (*KeepAliveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3e815100fbc7a5f2, []int{1}
}

func (m *KeepAliveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeepAliveResponse.Unmarshal(m, b)
}
func (m *KeepAliveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeepAliveResponse.Marshal(b, m, deterministic)
}
func (m *KeepAliveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeepAliveResponse.Merge(m, src)
}
func (m *KeepAliveResponse) XXX_Size() int {
	return xxx_messageInfo_KeepAliveResponse.Size(m)
}
func (m *KeepAliveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_KeepAliveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_KeepAliveResponse proto.InternalMessageInfo

func (m *KeepAliveResponse) GetIdEdge() int32 {
	if m != nil {
		return m.IdEdge
	}
	return 0
}

func (m *KeepAliveResponse) GetIdNodo() int32 {
	if m != nil {
		return m.IdNodo
	}
	return 0
}

func (m *KeepAliveResponse) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *KeepAliveResponse) GetDate() string {
	if m != nil {
		return m.Date
	}
	return ""
}

func (m *KeepAliveResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *KeepAliveResponse) GetConfirmation() bool {
	if m != nil {
		return m.Confirmation
	}
	return false
}

func init() {
	proto.RegisterType((*KeepAliveRequest)(nil), "keepalive.KeepAliveRequest")
	proto.RegisterType((*KeepAliveResponse)(nil), "keepalive.KeepAliveResponse")
}

func init() {
	proto.RegisterFile("proto/keepalive.proto", fileDescriptor_3e815100fbc7a5f2)
}

var fileDescriptor_3e815100fbc7a5f2 = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0xcf, 0x4e, 0x4d, 0x2d, 0x48, 0xcc, 0xc9, 0x2c, 0x4b, 0xd5, 0x03, 0xf3, 0x85, 0x38,
	0xe1, 0x02, 0x4a, 0x35, 0x5c, 0x02, 0xde, 0xa9, 0xa9, 0x05, 0x8e, 0x20, 0x4e, 0x50, 0x6a, 0x61,
	0x69, 0x6a, 0x71, 0x89, 0x90, 0x18, 0x17, 0x5b, 0x66, 0x8a, 0x6b, 0x4a, 0x7a, 0xaa, 0x04, 0xa3,
	0x02, 0xa3, 0x06, 0x6b, 0x10, 0x94, 0x07, 0x11, 0xf7, 0xcb, 0x4f, 0xc9, 0x97, 0x60, 0x82, 0x89,
	0x83, 0x78, 0x42, 0x42, 0x5c, 0x2c, 0x25, 0x95, 0x05, 0xa9, 0x12, 0xcc, 0x60, 0x51, 0x30, 0x1b,
	0x24, 0x96, 0x92, 0x58, 0x92, 0x2a, 0xc1, 0xa2, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0x66, 0x0b, 0x09,
	0x70, 0x31, 0xe7, 0x16, 0xa7, 0x4b, 0xb0, 0x82, 0x85, 0x40, 0x4c, 0xa5, 0x85, 0x8c, 0x5c, 0x82,
	0x48, 0xd6, 0x17, 0x17, 0xe4, 0xe7, 0x15, 0xa7, 0xd2, 0xcf, 0x7e, 0x21, 0x25, 0x2e, 0x9e, 0xe4,
	0xfc, 0xbc, 0xb4, 0xcc, 0xa2, 0xdc, 0xc4, 0x92, 0xcc, 0xfc, 0x3c, 0x09, 0x36, 0x05, 0x46, 0x0d,
	0x8e, 0x20, 0x14, 0x31, 0xa3, 0x48, 0x2e, 0x4e, 0xb8, 0x13, 0x85, 0x7c, 0xb8, 0x78, 0x7d, 0x13,
	0xb3, 0x53, 0x11, 0x02, 0xd2, 0x7a, 0x88, 0xc0, 0x45, 0x0f, 0x48, 0x29, 0x19, 0xec, 0x92, 0x10,
	0x6f, 0x2a, 0x31, 0x38, 0x71, 0x45, 0x71, 0xe8, 0xe9, 0x17, 0x24, 0x26, 0x67, 0xa7, 0xa7, 0x26,
	0xb1, 0x81, 0xa3, 0xc6, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x51, 0x98, 0x08, 0x89, 0xb3, 0x01,
	0x00, 0x00,
}
