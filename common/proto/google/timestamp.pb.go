// Code generated by protoc-gen-go.
// source: timestamp.proto
// DO NOT EDIT!

package google

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// A Timestamp represents a point in time independent of any time zone
// or calendar, represented as seconds and fractions of seconds at
// nanosecond resolution in UTC Epoch time. It is encoded using the
// Proleptic Gregorian Calendar which extends the Gregorian calendar
// backwards to year one. It is encoded assuming all minutes are 60
// seconds long, i.e. leap seconds are "smeared" so that no leap second
// table is needed for interpretation. Range is from
// 0001-01-01T00:00:00Z to 9999-12-31T23:59:59.999999999Z.
// By restricting to that range, we ensure that we can convert to
// and from  RFC 3339 date strings.
// See [https://www.ietf.org/rfc/rfc3339.txt](https://www.ietf.org/rfc/rfc3339.txt).
//
// Example 1: Compute Timestamp from POSIX `time()`.
//
//     Timestamp timestamp;
//     timestamp.set_seconds(time(NULL));
//     timestamp.set_nanos(0);
//
// Example 2: Compute Timestamp from POSIX `gettimeofday()`.
//
//     struct timeval tv;
//     gettimeofday(&tv, NULL);
//
//     Timestamp timestamp;
//     timestamp.set_seconds(tv.tv_sec);
//     timestamp.set_nanos(tv.tv_usec * 1000);
//
// Example 3: Compute Timestamp from Win32 `GetSystemTimeAsFileTime()`.
//
//     FILETIME ft;
//     GetSystemTimeAsFileTime(&ft);
//     UINT64 ticks = (((UINT64)ft.dwHighDateTime) << 32) | ft.dwLowDateTime;
//
//     // A Windows tick is 100 nanoseconds. Windows epoch 1601-01-01T00:00:00Z
//     // is 11644473600 seconds before Unix epoch 1970-01-01T00:00:00Z.
//     Timestamp timestamp;
//     timestamp.set_seconds((INT64) ((ticks / 10000000) - 11644473600LL));
//     timestamp.set_nanos((INT32) ((ticks % 10000000) * 100));
//
// Example 4: Compute Timestamp from Java `System.currentTimeMillis()`.
//
//     long millis = System.currentTimeMillis();
//
//     Timestamp timestamp = Timestamp.newBuilder().setSeconds(millis / 1000)
//         .setNanos((int) ((millis % 1000) * 1000000)).build();
//
//
// Example 5: Compute Timestamp from current time in Python.
//
//     now = time.time()
//     seconds = int(now)
//     nanos = int((now - seconds) * 10**9)
//     timestamp = Timestamp(seconds=seconds, nanos=nanos)
//
type Timestamp struct {
	// Represents seconds of UTC time since Unix epoch
	// 1970-01-01T00:00:00Z. Must be from from 0001-01-01T00:00:00Z to
	// 9999-12-31T23:59:59Z inclusive.
	Seconds int64 `protobuf:"varint,1,opt,name=seconds" json:"seconds,omitempty"`
	// Non-negative fractions of a second at nanosecond resolution. Negative
	// second values with fractions must still have non-negative nanos values
	// that count forward in time. Must be from 0 to 999,999,999
	// inclusive.
	Nanos int32 `protobuf:"varint,2,opt,name=nanos" json:"nanos,omitempty"`
}

func (m *Timestamp) Reset()                    { *m = Timestamp{} }
func (m *Timestamp) String() string            { return proto.CompactTextString(m) }
func (*Timestamp) ProtoMessage()               {}
func (*Timestamp) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func init() {
	proto.RegisterType((*Timestamp)(nil), "google.protobuf.Timestamp")
}

var fileDescriptor2 = []byte{
	// 149 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0xc9, 0xcc, 0x4d,
	0x2d, 0x2e, 0x49, 0xcc, 0x2d, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4f, 0xcf, 0xcf,
	0x4f, 0xcf, 0x49, 0x85, 0xf0, 0x92, 0x4a, 0xd3, 0x94, 0xb4, 0xb9, 0x38, 0x43, 0x60, 0x6a, 0x84,
	0xf8, 0xb9, 0xd8, 0x8b, 0x53, 0x93, 0xf3, 0xf3, 0x52, 0x8a, 0x25, 0x18, 0x15, 0x18, 0x35, 0x98,
	0x85, 0x78, 0xb9, 0x58, 0xf3, 0x12, 0xf3, 0xf2, 0x8b, 0x25, 0x98, 0x80, 0x5c, 0x56, 0xa7, 0x40,
	0x2e, 0xe1, 0xe4, 0xfc, 0x5c, 0x3d, 0x34, 0x33, 0x9c, 0xf8, 0xe0, 0x26, 0x04, 0x80, 0x84, 0x02,
	0x18, 0x17, 0x30, 0x32, 0x2e, 0x62, 0x62, 0x76, 0x0f, 0x70, 0x5a, 0xc5, 0x24, 0xe7, 0x0e, 0x51,
	0x1c, 0x00, 0x55, 0xac, 0x17, 0x9e, 0x9a, 0x93, 0xe3, 0x9d, 0x97, 0x5f, 0x9e, 0x17, 0x52, 0x59,
	0x90, 0x5a, 0x9c, 0xc4, 0x06, 0x36, 0xc5, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xd2, 0xa8, 0xe1,
	0x71, 0xaa, 0x00, 0x00, 0x00,
}
