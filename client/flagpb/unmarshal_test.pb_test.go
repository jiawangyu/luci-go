// Code generated by protoc-gen-go. DO NOT EDIT.
// source: go.chromium.org/luci/client/flagpb/unmarshal_test.proto

/*
Package flagpb is a generated protocol buffer package.

It is generated from these files:
	go.chromium.org/luci/client/flagpb/unmarshal_test.proto

It has these top-level messages:
	M1
	M2
	M3
	MapContainer
*/
package flagpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type E int32

const (
	E_V0 E = 0
	E_V1 E = 1
)

var E_name = map[int32]string{
	0: "V0",
	1: "V1",
}
var E_value = map[string]int32{
	"V0": 0,
	"V1": 1,
}

func (x E) String() string {
	return proto.EnumName(E_name, int32(x))
}
func (E) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type M1 struct {
	S  string  `protobuf:"bytes,1,opt,name=s" json:"s,omitempty"`
	I  int32   `protobuf:"varint,2,opt,name=i" json:"i,omitempty"`
	Ri []int32 `protobuf:"varint,3,rep,packed,name=ri" json:"ri,omitempty"`
	B  bool    `protobuf:"varint,4,opt,name=b" json:"b,omitempty"`
	Rb []bool  `protobuf:"varint,6,rep,packed,name=rb" json:"rb,omitempty"`
	Bb []byte  `protobuf:"bytes,5,opt,name=bb,proto3" json:"bb,omitempty"`
}

func (m *M1) Reset()                    { *m = M1{} }
func (m *M1) String() string            { return proto.CompactTextString(m) }
func (*M1) ProtoMessage()               {}
func (*M1) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *M1) GetS() string {
	if m != nil {
		return m.S
	}
	return ""
}

func (m *M1) GetI() int32 {
	if m != nil {
		return m.I
	}
	return 0
}

func (m *M1) GetRi() []int32 {
	if m != nil {
		return m.Ri
	}
	return nil
}

func (m *M1) GetB() bool {
	if m != nil {
		return m.B
	}
	return false
}

func (m *M1) GetRb() []bool {
	if m != nil {
		return m.Rb
	}
	return nil
}

func (m *M1) GetBb() []byte {
	if m != nil {
		return m.Bb
	}
	return nil
}

type M2 struct {
	M1 *M1 `protobuf:"bytes,1,opt,name=m1" json:"m1,omitempty"`
	E  E   `protobuf:"varint,2,opt,name=e,enum=flagpb.E" json:"e,omitempty"`
}

func (m *M2) Reset()                    { *m = M2{} }
func (m *M2) String() string            { return proto.CompactTextString(m) }
func (*M2) ProtoMessage()               {}
func (*M2) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *M2) GetM1() *M1 {
	if m != nil {
		return m.M1
	}
	return nil
}

func (m *M2) GetE() E {
	if m != nil {
		return m.E
	}
	return E_V0
}

type M3 struct {
	M1 []*M1  `protobuf:"bytes,1,rep,name=m1" json:"m1,omitempty"`
	M2 *M2    `protobuf:"bytes,2,opt,name=m2" json:"m2,omitempty"`
	B  bool   `protobuf:"varint,3,opt,name=b" json:"b,omitempty"`
	S  string `protobuf:"bytes,4,opt,name=s" json:"s,omitempty"`
	Bt []byte `protobuf:"bytes,5,opt,name=bt,proto3" json:"bt,omitempty"`
}

func (m *M3) Reset()                    { *m = M3{} }
func (m *M3) String() string            { return proto.CompactTextString(m) }
func (*M3) ProtoMessage()               {}
func (*M3) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *M3) GetM1() []*M1 {
	if m != nil {
		return m.M1
	}
	return nil
}

func (m *M3) GetM2() *M2 {
	if m != nil {
		return m.M2
	}
	return nil
}

func (m *M3) GetB() bool {
	if m != nil {
		return m.B
	}
	return false
}

func (m *M3) GetS() string {
	if m != nil {
		return m.S
	}
	return ""
}

func (m *M3) GetBt() []byte {
	if m != nil {
		return m.Bt
	}
	return nil
}

type MapContainer struct {
	Ss  map[string]string `protobuf:"bytes,1,rep,name=ss" json:"ss,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Ii  map[int32]int32   `protobuf:"bytes,2,rep,name=ii" json:"ii,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	Sm1 map[string]*M1    `protobuf:"bytes,3,rep,name=sm1" json:"sm1,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *MapContainer) Reset()                    { *m = MapContainer{} }
func (m *MapContainer) String() string            { return proto.CompactTextString(m) }
func (*MapContainer) ProtoMessage()               {}
func (*MapContainer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *MapContainer) GetSs() map[string]string {
	if m != nil {
		return m.Ss
	}
	return nil
}

func (m *MapContainer) GetIi() map[int32]int32 {
	if m != nil {
		return m.Ii
	}
	return nil
}

func (m *MapContainer) GetSm1() map[string]*M1 {
	if m != nil {
		return m.Sm1
	}
	return nil
}

func init() {
	proto.RegisterType((*M1)(nil), "flagpb.M1")
	proto.RegisterType((*M2)(nil), "flagpb.M2")
	proto.RegisterType((*M3)(nil), "flagpb.M3")
	proto.RegisterType((*MapContainer)(nil), "flagpb.MapContainer")
	proto.RegisterEnum("flagpb.E", E_name, E_value)
}

func init() {
	proto.RegisterFile("go.chromium.org/luci/client/flagpb/unmarshal_test.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 385 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x5f, 0xcb, 0xd3, 0x30,
	0x14, 0xc6, 0xcd, 0xe9, 0xdb, 0xba, 0x9d, 0x77, 0xbc, 0x8c, 0x28, 0x58, 0x86, 0x42, 0xe9, 0x55,
	0x11, 0x69, 0x6d, 0x87, 0xf8, 0xe7, 0x52, 0xd9, 0x85, 0x17, 0xbd, 0x89, 0xe0, 0xa5, 0xd2, 0x8c,
	0xba, 0x05, 0x9b, 0x76, 0x24, 0xa9, 0xb0, 0x0f, 0xe9, 0x77, 0x92, 0xa4, 0xed, 0xac, 0x50, 0xaf,
	0xda, 0xc3, 0xf3, 0x3c, 0xe7, 0x77, 0x72, 0x12, 0x7c, 0x7b, 0xea, 0xd2, 0xe3, 0x59, 0x75, 0x52,
	0xf4, 0x32, 0xed, 0xd4, 0x29, 0x6b, 0xfa, 0xa3, 0xc8, 0x8e, 0x8d, 0xa8, 0x5b, 0x93, 0xfd, 0x68,
	0xaa, 0xd3, 0x85, 0x67, 0x7d, 0x2b, 0x2b, 0xa5, 0xcf, 0x55, 0xf3, 0xdd, 0xd4, 0xda, 0xa4, 0x17,
	0xd5, 0x99, 0x8e, 0x06, 0x83, 0x18, 0x7f, 0x43, 0x28, 0x73, 0xba, 0x41, 0xa2, 0x43, 0x12, 0x91,
	0x64, 0xcd, 0x88, 0xb6, 0x95, 0x08, 0x21, 0x22, 0x89, 0xcf, 0x88, 0xa0, 0x0f, 0x08, 0x4a, 0x84,
	0x5e, 0xe4, 0x25, 0x3e, 0x03, 0x25, 0xac, 0xca, 0xc3, 0xbb, 0x88, 0x24, 0x2b, 0x46, 0xb8, 0x53,
	0x79, 0x18, 0x44, 0x5e, 0xb2, 0x62, 0xa0, 0x5c, 0xcd, 0x79, 0xe8, 0x47, 0x24, 0xd9, 0x30, 0xe0,
	0x3c, 0x7e, 0x8f, 0x50, 0x16, 0x74, 0x87, 0x20, 0x73, 0x07, 0xb8, 0x2f, 0x30, 0x1d, 0xd0, 0x69,
	0x99, 0x33, 0x90, 0x39, 0x7d, 0x86, 0xa4, 0x76, 0xb4, 0x87, 0x62, 0x3d, 0x49, 0x07, 0x46, 0xea,
	0xf8, 0x8c, 0x50, 0xee, 0x6f, 0x51, 0x6f, 0x21, 0x6a, 0xb5, 0xc2, 0x65, 0xe7, 0x5a, 0xc1, 0x40,
	0x16, 0xc3, 0x98, 0xde, 0x34, 0xa6, 0x3b, 0xe0, 0xdd, 0x74, 0x40, 0x3b, 0xa4, 0xb9, 0x0d, 0x69,
	0xe2, 0xdf, 0x80, 0x9b, 0xb2, 0xba, 0x7c, 0xea, 0x5a, 0x53, 0x89, 0xb6, 0x56, 0xf4, 0x15, 0x82,
	0xd6, 0x23, 0xf4, 0xf9, 0xad, 0xf1, 0xcc, 0x91, 0x7e, 0xd1, 0x87, 0xd6, 0xa8, 0x2b, 0x03, 0xad,
	0xad, 0x5b, 0xd8, 0x85, 0xfd, 0xdf, 0xfd, 0x59, 0x8c, 0x6e, 0x21, 0x68, 0x86, 0x9e, 0x96, 0xb9,
	0x5b, 0xe8, 0x7d, 0xf1, 0x62, 0xb9, 0xb9, 0xcc, 0x07, 0xbf, 0x75, 0xee, 0xde, 0xe0, 0xe3, 0x91,
	0x46, 0xb7, 0xe8, 0xfd, 0xac, 0xaf, 0xe3, 0x4d, 0xd9, 0x5f, 0xfa, 0x14, 0xfd, 0x5f, 0x55, 0xd3,
	0x0f, 0x1b, 0x5c, 0xb3, 0xa1, 0xf8, 0x00, 0xef, 0x88, 0x8d, 0x8d, 0xd8, 0x79, 0xcc, 0x5f, 0x88,
	0xf9, 0xf3, 0xd8, 0x47, 0x5c, 0x4d, 0xf8, 0x05, 0x5c, 0x34, 0xcf, 0xfd, 0x7b, 0x21, 0x7f, 0x7b,
	0xbc, 0x7c, 0x82, 0xe4, 0x40, 0x03, 0x84, 0xaf, 0xaf, 0xb7, 0x8f, 0xdc, 0x37, 0xdf, 0x12, 0x1e,
	0xb8, 0x87, 0xb7, 0xff, 0x13, 0x00, 0x00, 0xff, 0xff, 0x7f, 0x7d, 0x43, 0xa5, 0xb3, 0x02, 0x00,
	0x00,
}
