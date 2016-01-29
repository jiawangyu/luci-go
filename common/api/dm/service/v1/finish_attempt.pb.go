// Code generated by protoc-gen-go.
// source: finish_attempt.proto
// DO NOT EDIT!

package dm

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/luci/luci-go/common/proto/google"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// FinishAttemptReq sets the final result of an Attempt.
type FinishAttemptReq struct {
	// required
	Auth       *ExecutionAuth             `protobuf:"bytes,1,opt,name=auth" json:"auth,omitempty"`
	JsonResult string                     `protobuf:"bytes,2,opt,name=json_result" json:"json_result,omitempty"`
	Expiration *google_protobuf.Timestamp `protobuf:"bytes,3,opt,name=expiration" json:"expiration,omitempty"`
}

func (m *FinishAttemptReq) Reset()                    { *m = FinishAttemptReq{} }
func (m *FinishAttemptReq) String() string            { return proto.CompactTextString(m) }
func (*FinishAttemptReq) ProtoMessage()               {}
func (*FinishAttemptReq) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

func (m *FinishAttemptReq) GetAuth() *ExecutionAuth {
	if m != nil {
		return m.Auth
	}
	return nil
}

func (m *FinishAttemptReq) GetExpiration() *google_protobuf.Timestamp {
	if m != nil {
		return m.Expiration
	}
	return nil
}

func init() {
	proto.RegisterType((*FinishAttemptReq)(nil), "dm.FinishAttemptReq")
}

var fileDescriptor4 = []byte{
	// 182 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x34, 0x8d, 0xc1, 0x8e, 0x82, 0x30,
	0x10, 0x86, 0x03, 0xbb, 0xd9, 0x64, 0xcb, 0x65, 0xb7, 0x7a, 0x20, 0x5c, 0x30, 0x9e, 0x3c, 0x0d,
	0x89, 0x3e, 0x01, 0x07, 0x7d, 0x00, 0xe3, 0x9d, 0x14, 0x19, 0xa0, 0x86, 0xd2, 0x4a, 0xa7, 0x09,
	0xbe, 0xbd, 0xd8, 0x86, 0xe3, 0xfc, 0xf3, 0xfd, 0xff, 0xc7, 0xb6, 0xad, 0x1c, 0xa5, 0xed, 0x2b,
	0x41, 0x84, 0xca, 0x10, 0x98, 0x49, 0x93, 0xe6, 0x71, 0xa3, 0xb2, 0x84, 0x5e, 0x06, 0x6d, 0x08,
	0xb2, 0xbc, 0xd3, 0xba, 0x1b, 0xb0, 0xf0, 0x57, 0xed, 0xda, 0x82, 0xa4, 0x42, 0x4b, 0x42, 0x99,
	0x00, 0xec, 0x67, 0xf6, 0x77, 0xf1, 0x4b, 0x65, 0x18, 0xba, 0xe2, 0x93, 0xe7, 0xec, 0x5b, 0x38,
	0xea, 0xd3, 0x68, 0x17, 0x1d, 0x92, 0xe3, 0x3f, 0x34, 0x0a, 0xce, 0x33, 0xde, 0x1d, 0x49, 0x3d,
	0x96, 0xcb, 0x83, 0x6f, 0x58, 0xf2, 0xb0, 0x7a, 0xac, 0x26, 0xb4, 0x6e, 0xa0, 0x34, 0x5e, 0xb8,
	0x5f, 0x0e, 0x8c, 0xe1, 0x6c, 0xe4, 0x24, 0x3e, 0x58, 0xfa, 0xe5, 0xbb, 0x19, 0x04, 0x3f, 0xac,
	0x7e, 0xb8, 0xad, 0xfe, 0xfa, 0xc7, 0x67, 0xa7, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0xa4, 0x49,
	0x36, 0x71, 0xca, 0x00, 0x00, 0x00,
}
