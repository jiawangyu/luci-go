// Code generated by protoc-gen-go. DO NOT EDIT.
// source: go.chromium.org/luci/skylark/skylarkproto/testprotos/another.proto

/*
Package testprotos is a generated protocol buffer package.

It is generated from these files:
	go.chromium.org/luci/skylark/skylarkproto/testprotos/another.proto
	go.chromium.org/luci/skylark/skylarkproto/testprotos/test.proto

It has these top-level messages:
	AnotherMessage
	SimpleFields
	MessageFields
	Simple
	Complex
*/
package testprotos

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

type AnotherMessage struct {
}

func (m *AnotherMessage) Reset()                    { *m = AnotherMessage{} }
func (m *AnotherMessage) String() string            { return proto.CompactTextString(m) }
func (*AnotherMessage) ProtoMessage()               {}
func (*AnotherMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*AnotherMessage)(nil), "testprotos.AnotherMessage")
}

func init() {
	proto.RegisterFile("go.chromium.org/luci/skylark/skylarkproto/testprotos/another.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 102 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x72, 0x4a, 0xcf, 0xd7, 0x4b,
	0xce, 0x28, 0xca, 0xcf, 0xcd, 0x2c, 0xcd, 0xd5, 0xcb, 0x2f, 0x4a, 0xd7, 0xcf, 0x29, 0x4d, 0xce,
	0xd4, 0x2f, 0xce, 0xae, 0xcc, 0x49, 0x2c, 0xca, 0x86, 0xd1, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0xfa,
	0x25, 0xa9, 0xc5, 0x25, 0x60, 0x56, 0xb1, 0x7e, 0x62, 0x5e, 0x7e, 0x49, 0x46, 0x6a, 0x91, 0x1e,
	0x98, 0x2b, 0xc4, 0x85, 0x90, 0x51, 0x12, 0xe0, 0xe2, 0x73, 0x84, 0x48, 0xfa, 0xa6, 0x16, 0x17,
	0x27, 0xa6, 0xa7, 0x26, 0xb1, 0x81, 0x65, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x65, 0x11,
	0x8d, 0x1a, 0x6a, 0x00, 0x00, 0x00,
}
