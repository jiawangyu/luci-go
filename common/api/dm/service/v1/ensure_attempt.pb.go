// Code generated by protoc-gen-go.
// source: ensure_attempt.proto
// DO NOT EDIT!

package dm

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// EnsureAttemptReq allows you to ensure that a given Attempt exists.
//
// This is used to start a new Attempt. Retrying this on failure is idempotent
// (since the Attempt's ID is explicitly provided in the request).
//
// Programmatic addition of new Attempts from an Execution must use the AddDeps
// api instead.
//
// Any quests referred to by this method must already have been established with
// DM by calling EnsureQuests.
type EnsureAttemptReq struct {
	ToEnsure *AttemptID `protobuf:"bytes,1,opt,name=to_ensure" json:"to_ensure,omitempty"`
}

func (m *EnsureAttemptReq) Reset()                    { *m = EnsureAttemptReq{} }
func (m *EnsureAttemptReq) String() string            { return proto.CompactTextString(m) }
func (*EnsureAttemptReq) ProtoMessage()               {}
func (*EnsureAttemptReq) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *EnsureAttemptReq) GetToEnsure() *AttemptID {
	if m != nil {
		return m.ToEnsure
	}
	return nil
}

func init() {
	proto.RegisterType((*EnsureAttemptReq)(nil), "dm.EnsureAttemptReq")
}

var fileDescriptor2 = []byte{
	// 102 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x49, 0xcd, 0x2b, 0x2e,
	0x2d, 0x4a, 0x8d, 0x4f, 0x2c, 0x29, 0x49, 0xcd, 0x2d, 0x28, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x4a, 0xc9, 0x95, 0xe2, 0x2e, 0xa9, 0x2c, 0x48, 0x2d, 0x86, 0x08, 0x28, 0x99, 0x70,
	0x09, 0xb8, 0x82, 0x15, 0x3a, 0x42, 0xd4, 0x05, 0xa5, 0x16, 0x0a, 0x29, 0x70, 0x71, 0x96, 0xe4,
	0xc7, 0x43, 0xf4, 0x4b, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x1b, 0xf1, 0xea, 0xa5, 0xe4, 0xea, 0x41,
	0x95, 0x78, 0xba, 0x24, 0xb1, 0x81, 0x35, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xc2, 0x93,
	0x45, 0xe1, 0x65, 0x00, 0x00, 0x00,
}
