// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/common/tsmon/ts_mon_proto_v2/acquisition_network_device.proto
// DO NOT EDIT!

/*
Package ts_mon_proto_v2 is a generated protocol buffer package.

It is generated from these files:
	github.com/luci/luci-go/common/tsmon/ts_mon_proto_v2/acquisition_network_device.proto
	github.com/luci/luci-go/common/tsmon/ts_mon_proto_v2/acquisition_task.proto
	github.com/luci/luci-go/common/tsmon/ts_mon_proto_v2/any.proto
	github.com/luci/luci-go/common/tsmon/ts_mon_proto_v2/endpoint.proto
	github.com/luci/luci-go/common/tsmon/ts_mon_proto_v2/metrics.proto
	github.com/luci/luci-go/common/tsmon/ts_mon_proto_v2/timestamp.proto

It has these top-level messages:
	NetworkDevice
	Task
	Any
	Request
	MetricsPayload
	MetricsCollection
	MetricsDataSet
	MetricsData
	Annotations
	Timestamp
*/
package ts_mon_proto_v2

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

type NetworkDevice_TypeId int32

const (
	NetworkDevice_MESSAGE_TYPE_ID NetworkDevice_TypeId = 34049749
)

var NetworkDevice_TypeId_name = map[int32]string{
	34049749: "MESSAGE_TYPE_ID",
}
var NetworkDevice_TypeId_value = map[string]int32{
	"MESSAGE_TYPE_ID": 34049749,
}

func (x NetworkDevice_TypeId) Enum() *NetworkDevice_TypeId {
	p := new(NetworkDevice_TypeId)
	*p = x
	return p
}
func (x NetworkDevice_TypeId) String() string {
	return proto.EnumName(NetworkDevice_TypeId_name, int32(x))
}
func (x *NetworkDevice_TypeId) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(NetworkDevice_TypeId_value, data, "NetworkDevice_TypeId")
	if err != nil {
		return err
	}
	*x = NetworkDevice_TypeId(value)
	return nil
}
func (NetworkDevice_TypeId) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type NetworkDevice struct {
	ProxyEnvironment *string `protobuf:"bytes,5,opt,name=proxy_environment,json=proxyEnvironment" json:"proxy_environment,omitempty"`
	AcquisitionName  *string `protobuf:"bytes,10,opt,name=acquisition_name,json=acquisitionName" json:"acquisition_name,omitempty"`
	Pop              *string `protobuf:"bytes,30,opt,name=pop" json:"pop,omitempty"`
	Alertable        *bool   `protobuf:"varint,101,opt,name=alertable" json:"alertable,omitempty"`
	Realm            *string `protobuf:"bytes,102,opt,name=realm" json:"realm,omitempty"`
	Asn              *int64  `protobuf:"varint,103,opt,name=asn" json:"asn,omitempty"`
	Metro            *string `protobuf:"bytes,104,opt,name=metro" json:"metro,omitempty"`
	Role             *string `protobuf:"bytes,105,opt,name=role" json:"role,omitempty"`
	Hostname         *string `protobuf:"bytes,106,opt,name=hostname" json:"hostname,omitempty"`
	Vendor           *string `protobuf:"bytes,70,opt,name=vendor" json:"vendor,omitempty"`
	Hostgroup        *string `protobuf:"bytes,108,opt,name=hostgroup" json:"hostgroup,omitempty"`
	ProxyZone        *string `protobuf:"bytes,100,opt,name=proxy_zone,json=proxyZone" json:"proxy_zone,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *NetworkDevice) Reset()                    { *m = NetworkDevice{} }
func (m *NetworkDevice) String() string            { return proto.CompactTextString(m) }
func (*NetworkDevice) ProtoMessage()               {}
func (*NetworkDevice) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *NetworkDevice) GetProxyEnvironment() string {
	if m != nil && m.ProxyEnvironment != nil {
		return *m.ProxyEnvironment
	}
	return ""
}

func (m *NetworkDevice) GetAcquisitionName() string {
	if m != nil && m.AcquisitionName != nil {
		return *m.AcquisitionName
	}
	return ""
}

func (m *NetworkDevice) GetPop() string {
	if m != nil && m.Pop != nil {
		return *m.Pop
	}
	return ""
}

func (m *NetworkDevice) GetAlertable() bool {
	if m != nil && m.Alertable != nil {
		return *m.Alertable
	}
	return false
}

func (m *NetworkDevice) GetRealm() string {
	if m != nil && m.Realm != nil {
		return *m.Realm
	}
	return ""
}

func (m *NetworkDevice) GetAsn() int64 {
	if m != nil && m.Asn != nil {
		return *m.Asn
	}
	return 0
}

func (m *NetworkDevice) GetMetro() string {
	if m != nil && m.Metro != nil {
		return *m.Metro
	}
	return ""
}

func (m *NetworkDevice) GetRole() string {
	if m != nil && m.Role != nil {
		return *m.Role
	}
	return ""
}

func (m *NetworkDevice) GetHostname() string {
	if m != nil && m.Hostname != nil {
		return *m.Hostname
	}
	return ""
}

func (m *NetworkDevice) GetVendor() string {
	if m != nil && m.Vendor != nil {
		return *m.Vendor
	}
	return ""
}

func (m *NetworkDevice) GetHostgroup() string {
	if m != nil && m.Hostgroup != nil {
		return *m.Hostgroup
	}
	return ""
}

func (m *NetworkDevice) GetProxyZone() string {
	if m != nil && m.ProxyZone != nil {
		return *m.ProxyZone
	}
	return ""
}

func init() {
	proto.RegisterType((*NetworkDevice)(nil), "ts_mon.proto.v2.NetworkDevice")
	proto.RegisterEnum("ts_mon.proto.v2.NetworkDevice_TypeId", NetworkDevice_TypeId_name, NetworkDevice_TypeId_value)
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/common/tsmon/ts_mon_proto_v2/acquisition_network_device.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 335 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x54, 0x90, 0xd1, 0x4a, 0xc3, 0x30,
	0x14, 0x86, 0x99, 0x73, 0x63, 0x0b, 0xc8, 0x6a, 0x90, 0x11, 0x44, 0xa5, 0xec, 0x6a, 0x22, 0xae,
	0xb0, 0x37, 0x10, 0x56, 0x65, 0x17, 0x0e, 0xd9, 0xe6, 0x85, 0xde, 0x84, 0xae, 0x3d, 0x76, 0xd1,
	0x26, 0xa7, 0xa6, 0x69, 0x75, 0xbe, 0x8b, 0xbe, 0x95, 0x6f, 0xe3, 0x85, 0x34, 0x11, 0xa7, 0x37,
	0xe1, 0xfc, 0xdf, 0xff, 0xa5, 0x3d, 0x84, 0xdc, 0xa6, 0xc2, 0xac, 0xcb, 0xd5, 0x28, 0x46, 0x19,
	0x64, 0x65, 0x2c, 0xec, 0x71, 0x9e, 0x62, 0x10, 0xa3, 0x94, 0xa8, 0x02, 0x53, 0xb8, 0x93, 0x4b,
	0x54, 0x3c, 0xd7, 0x68, 0x90, 0x57, 0xe3, 0x20, 0x8a, 0x9f, 0x4b, 0x51, 0x08, 0x23, 0x50, 0x71,
	0x05, 0xe6, 0x05, 0xf5, 0x13, 0x4f, 0xa0, 0x12, 0x31, 0x8c, 0xac, 0x43, 0x7b, 0xee, 0x86, 0x4b,
	0xa3, 0x6a, 0x3c, 0xf8, 0xda, 0x21, 0x7b, 0x33, 0x67, 0x4e, 0xac, 0x48, 0xcf, 0xc8, 0x7e, 0xae,
	0xf1, 0x75, 0xc3, 0x41, 0x55, 0x42, 0xa3, 0x92, 0xa0, 0x0c, 0x6b, 0xf9, 0x8d, 0x61, 0x77, 0xee,
	0xd9, 0x22, 0xdc, 0x72, 0x7a, 0x4a, 0xbc, 0x7f, 0xff, 0x8c, 0x24, 0x30, 0x62, 0xdd, 0xde, 0x1f,
	0x3e, 0x8b, 0x24, 0x50, 0x8f, 0x34, 0x73, 0xcc, 0xd9, 0x89, 0x6d, 0xeb, 0x91, 0x1e, 0x91, 0x6e,
	0x94, 0x81, 0x36, 0xd1, 0x2a, 0x03, 0x06, 0x7e, 0x63, 0xd8, 0x99, 0x6f, 0x01, 0x3d, 0x20, 0x2d,
	0x0d, 0x51, 0x26, 0xd9, 0x83, 0xbd, 0xe1, 0x42, 0xfd, 0x95, 0xa8, 0x50, 0x2c, 0xf5, 0x1b, 0xc3,
	0xe6, 0xbc, 0x1e, 0x6b, 0x4f, 0x82, 0xd1, 0xc8, 0xd6, 0xce, 0xb3, 0x81, 0x52, 0xb2, 0xab, 0x31,
	0x03, 0x26, 0x2c, 0xb4, 0x33, 0x3d, 0x24, 0x9d, 0x35, 0x16, 0xc6, 0x2e, 0xf9, 0x68, 0xf9, 0x6f,
	0xa6, 0x7d, 0xd2, 0xae, 0x40, 0x25, 0xa8, 0xd9, 0xa5, 0x6d, 0x7e, 0x52, 0xbd, 0x63, 0xed, 0xa4,
	0x1a, 0xcb, 0x9c, 0x65, 0xb6, 0xda, 0x02, 0x7a, 0x4c, 0x88, 0x7b, 0xab, 0x37, 0x54, 0xc0, 0x12,
	0x57, 0x5b, 0x72, 0x8f, 0x0a, 0x06, 0x3e, 0x69, 0x2f, 0x37, 0x39, 0x4c, 0x13, 0xda, 0x27, 0xbd,
	0xeb, 0x70, 0xb1, 0xb8, 0xb8, 0x0a, 0xf9, 0xf2, 0xee, 0x26, 0xe4, 0xd3, 0x89, 0xf7, 0xf9, 0xfe,
	0xe1, 0x7d, 0x07, 0x00, 0x00, 0xff, 0xff, 0x71, 0x3e, 0x8c, 0x55, 0xe7, 0x01, 0x00, 0x00,
}
