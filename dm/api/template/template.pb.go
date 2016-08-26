// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/dm/api/template/template.proto
// DO NOT EDIT!

/*
Package dmTemplate is a generated protocol buffer package.

It is generated from these files:
	github.com/luci/luci-go/dm/api/template/template.proto

It has these top-level messages:
	File
*/
package dmTemplate

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import templateproto "github.com/luci/luci-go/common/data/text/templateproto"
import dm1 "github.com/luci/luci-go/dm/api/service/v1"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// File represents a file full of DM template definitions.
type File struct {
	Template map[string]*File_Template `protobuf:"bytes,1,rep,name=template" json:"template,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *File) Reset()                    { *m = File{} }
func (m *File) String() string            { return proto.CompactTextString(m) }
func (*File) ProtoMessage()               {}
func (*File) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *File) GetTemplate() map[string]*File_Template {
	if m != nil {
		return m.Template
	}
	return nil
}

// Template defines a single template.
type File_Template struct {
	DistributorConfigName string                       `protobuf:"bytes,1,opt,name=distributor_config_name,json=distributorConfigName" json:"distributor_config_name,omitempty"`
	Parameters            *templateproto.File_Template `protobuf:"bytes,2,opt,name=parameters" json:"parameters,omitempty"`
	DistributorParameters *templateproto.File_Template `protobuf:"bytes,3,opt,name=distributor_parameters,json=distributorParameters" json:"distributor_parameters,omitempty"`
	Meta                  *dm1.Quest_Desc_Meta         `protobuf:"bytes,4,opt,name=meta" json:"meta,omitempty"`
}

func (m *File_Template) Reset()                    { *m = File_Template{} }
func (m *File_Template) String() string            { return proto.CompactTextString(m) }
func (*File_Template) ProtoMessage()               {}
func (*File_Template) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *File_Template) GetParameters() *templateproto.File_Template {
	if m != nil {
		return m.Parameters
	}
	return nil
}

func (m *File_Template) GetDistributorParameters() *templateproto.File_Template {
	if m != nil {
		return m.DistributorParameters
	}
	return nil
}

func (m *File_Template) GetMeta() *dm1.Quest_Desc_Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func init() {
	proto.RegisterType((*File)(nil), "dmTemplate.File")
	proto.RegisterType((*File_Template)(nil), "dmTemplate.File.Template")
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/dm/api/template/template.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 327 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x90, 0x51, 0x4b, 0xfb, 0x30,
	0x14, 0xc5, 0xe9, 0xba, 0xff, 0x9f, 0x79, 0x87, 0x20, 0x11, 0x75, 0x16, 0x11, 0xf1, 0xc5, 0xbd,
	0x98, 0xe0, 0x84, 0x21, 0xc3, 0x37, 0x9d, 0x6f, 0x8a, 0x56, 0xf1, 0x75, 0x64, 0x6d, 0xec, 0x82,
	0x4d, 0x53, 0xd2, 0x74, 0xb8, 0xcf, 0xe2, 0x77, 0x15, 0xd3, 0x6c, 0xed, 0x5a, 0x61, 0xe8, 0x4b,
	0x39, 0xcd, 0x39, 0xf9, 0xe5, 0xdc, 0x0b, 0xc3, 0x88, 0xeb, 0x59, 0x3e, 0xc5, 0x81, 0x14, 0x24,
	0xce, 0x03, 0x6e, 0x3f, 0xe7, 0x91, 0x24, 0xa1, 0x20, 0x34, 0xe5, 0x44, 0x33, 0x91, 0xc6, 0x54,
	0xb3, 0x4a, 0xe0, 0x54, 0x49, 0x2d, 0x11, 0x84, 0xe2, 0x65, 0x75, 0xe2, 0x8d, 0x37, 0x31, 0xcc,
	0x81, 0x90, 0x09, 0x09, 0xa9, 0xa6, 0xe6, 0xfe, 0x87, 0xae, 0x20, 0x96, 0xf1, 0x03, 0xe9, 0x8d,
	0x7e, 0xa9, 0x92, 0x31, 0x35, 0xe7, 0x01, 0x23, 0xf3, 0x0b, 0x12, 0x29, 0x9a, 0xce, 0x26, 0x05,
	0x77, 0x79, 0xf7, 0xf4, 0xd3, 0x85, 0xf6, 0x1d, 0x8f, 0x19, 0x1a, 0x41, 0xa7, 0xc4, 0xf6, 0x9c,
	0x13, 0xb7, 0xdf, 0x1d, 0x1c, 0xe3, 0x75, 0x55, 0x5c, 0x64, 0x70, 0xf9, 0x37, 0x4e, 0xb4, 0x5a,
	0xf8, 0x55, 0xde, 0xfb, 0x72, 0xa0, 0x53, 0x7a, 0x68, 0x08, 0x07, 0x21, 0xcf, 0xb4, 0xe2, 0xd3,
	0x5c, 0x4b, 0x35, 0x09, 0x64, 0xf2, 0xc6, 0xa3, 0x49, 0x42, 0x45, 0xc1, 0x75, 0xfa, 0x5b, 0xfe,
	0x5e, 0xcd, 0xbe, 0xb1, 0xee, 0x83, 0x31, 0xd1, 0x35, 0x40, 0x4a, 0x95, 0x51, 0x9a, 0xa9, 0xac,
	0xd7, 0x32, 0xd1, 0xee, 0xe0, 0x08, 0x37, 0x06, 0x6f, 0xb6, 0xf0, 0x6b, 0x79, 0xf4, 0x0c, 0xfb,
	0xf5, 0x57, 0x6b, 0x24, 0xf7, 0x0f, 0xa4, 0x7a, 0xa5, 0xc7, 0x35, 0xf4, 0x0c, 0xda, 0x46, 0xd1,
	0x5e, 0xdb, 0x22, 0x76, 0xcd, 0x3e, 0xf0, 0x53, 0xce, 0x32, 0x8d, 0x6f, 0x59, 0x16, 0xe0, 0x7b,
	0x63, 0xf9, 0x36, 0xe0, 0xbd, 0xc2, 0x76, 0x63, 0x37, 0x68, 0x07, 0xdc, 0x77, 0xb6, 0x58, 0x0d,
	0x5c, 0x48, 0x44, 0xe0, 0xdf, 0x9c, 0xc6, 0x39, 0x5b, 0x4d, 0x76, 0xb8, 0x71, 0xb9, 0xfe, 0x32,
	0x37, 0x6a, 0x5d, 0x39, 0xd3, 0xff, 0xb6, 0xec, 0xe5, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3d,
	0x8e, 0xec, 0x55, 0x6d, 0x02, 0x00, 0x00,
}