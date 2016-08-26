// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/grpc/discovery/internal/testservices/helloworld.proto
// DO NOT EDIT!

/*
Package testservices is a generated protocol buffer package.

It is generated from these files:
	github.com/luci/luci-go/grpc/discovery/internal/testservices/helloworld.proto

It has these top-level messages:
	HelloRequest
	HelloReply
	MultiplyRequest
	MultiplyResponse
*/
package testservices

import prpc "github.com/luci/luci-go/grpc/prpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the user's name.
type HelloRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *HelloRequest) Reset()                    { *m = HelloRequest{} }
func (m *HelloRequest) String() string            { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()               {}
func (*HelloRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// The response message containing the greetings
type HelloReply struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *HelloReply) Reset()                    { *m = HelloReply{} }
func (m *HelloReply) String() string            { return proto.CompactTextString(m) }
func (*HelloReply) ProtoMessage()               {}
func (*HelloReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type MultiplyRequest struct {
	X int32 `protobuf:"varint,1,opt,name=x" json:"x,omitempty"`
	Y int32 `protobuf:"varint,2,opt,name=y" json:"y,omitempty"`
}

func (m *MultiplyRequest) Reset()                    { *m = MultiplyRequest{} }
func (m *MultiplyRequest) String() string            { return proto.CompactTextString(m) }
func (*MultiplyRequest) ProtoMessage()               {}
func (*MultiplyRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type MultiplyResponse struct {
	Z int32 `protobuf:"varint,1,opt,name=z" json:"z,omitempty"`
}

func (m *MultiplyResponse) Reset()                    { *m = MultiplyResponse{} }
func (m *MultiplyResponse) String() string            { return proto.CompactTextString(m) }
func (*MultiplyResponse) ProtoMessage()               {}
func (*MultiplyResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*HelloRequest)(nil), "testservices.HelloRequest")
	proto.RegisterType((*HelloReply)(nil), "testservices.HelloReply")
	proto.RegisterType((*MultiplyRequest)(nil), "testservices.MultiplyRequest")
	proto.RegisterType((*MultiplyResponse)(nil), "testservices.MultiplyResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Greeter service

type GreeterClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}
type greeterPRPCClient struct {
	client *prpc.Client
}

func NewGreeterPRPCClient(client *prpc.Client) GreeterClient {
	return &greeterPRPCClient{client}
}

func (c *greeterPRPCClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.client.Call(ctx, "testservices.Greeter", "SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := grpc.Invoke(ctx, "/testservices.Greeter/SayHello", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Greeter service

type GreeterServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterServer(s prpc.Registrar, srv GreeterServer) {
	s.RegisterService(&_Greeter_serviceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/testservices.Greeter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "testservices.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

// Client API for Calc service

type CalcClient interface {
	Multiply(ctx context.Context, in *MultiplyRequest, opts ...grpc.CallOption) (*MultiplyResponse, error)
}
type calcPRPCClient struct {
	client *prpc.Client
}

func NewCalcPRPCClient(client *prpc.Client) CalcClient {
	return &calcPRPCClient{client}
}

func (c *calcPRPCClient) Multiply(ctx context.Context, in *MultiplyRequest, opts ...grpc.CallOption) (*MultiplyResponse, error) {
	out := new(MultiplyResponse)
	err := c.client.Call(ctx, "testservices.Calc", "Multiply", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type calcClient struct {
	cc *grpc.ClientConn
}

func NewCalcClient(cc *grpc.ClientConn) CalcClient {
	return &calcClient{cc}
}

func (c *calcClient) Multiply(ctx context.Context, in *MultiplyRequest, opts ...grpc.CallOption) (*MultiplyResponse, error) {
	out := new(MultiplyResponse)
	err := grpc.Invoke(ctx, "/testservices.Calc/Multiply", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Calc service

type CalcServer interface {
	Multiply(context.Context, *MultiplyRequest) (*MultiplyResponse, error)
}

func RegisterCalcServer(s prpc.Registrar, srv CalcServer) {
	s.RegisterService(&_Calc_serviceDesc, srv)
}

func _Calc_Multiply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MultiplyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalcServer).Multiply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/testservices.Calc/Multiply",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalcServer).Multiply(ctx, req.(*MultiplyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Calc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "testservices.Calc",
	HandlerType: (*CalcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Multiply",
			Handler:    _Calc_Multiply_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/grpc/discovery/internal/testservices/helloworld.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 264 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x90, 0x51, 0x4b, 0xc3, 0x30,
	0x14, 0x85, 0x57, 0x99, 0x6e, 0x5e, 0x0a, 0x4a, 0x9e, 0x4a, 0x41, 0x19, 0x79, 0x10, 0x5f, 0xd6,
	0xc0, 0xfc, 0x07, 0xfa, 0xa0, 0x20, 0x7d, 0xe9, 0x7e, 0x41, 0x97, 0x5d, 0xba, 0x40, 0xda, 0xd4,
	0x24, 0x9d, 0xcb, 0x7e, 0xbd, 0x6d, 0x66, 0xb0, 0x8a, 0x7b, 0x09, 0x39, 0x39, 0x1f, 0xb9, 0xe7,
	0x5c, 0xc8, 0x2b, 0x61, 0x77, 0xdd, 0x26, 0xe3, 0xaa, 0x66, 0xb2, 0xe3, 0xc2, 0x1f, 0xcb, 0x4a,
	0xb1, 0x4a, 0xb7, 0x9c, 0x6d, 0x85, 0xe1, 0x6a, 0x8f, 0xda, 0x31, 0xd1, 0x58, 0xd4, 0x4d, 0x29,
	0x99, 0x45, 0x63, 0x0d, 0xea, 0xbd, 0xe0, 0x68, 0xd8, 0x0e, 0xa5, 0x54, 0x9f, 0x4a, 0xcb, 0x6d,
	0xd6, 0x6a, 0x65, 0x15, 0x89, 0xc7, 0x36, 0xa5, 0x10, 0xbf, 0x0d, 0x44, 0x81, 0x1f, 0x5d, 0xff,
	0x4e, 0x08, 0x4c, 0x9b, 0xb2, 0xc6, 0x24, 0x5a, 0x44, 0x8f, 0xd7, 0x85, 0xbf, 0xd3, 0x07, 0x80,
	0x6f, 0xa6, 0x95, 0x8e, 0x24, 0x30, 0xab, 0xd1, 0x98, 0xb2, 0x0a, 0x50, 0x90, 0x74, 0x09, 0x37,
	0x79, 0x27, 0xad, 0xe8, 0xa9, 0xf0, 0x5d, 0x0c, 0xd1, 0xc1, 0x63, 0x97, 0x45, 0x74, 0x18, 0x94,
	0x4b, 0x2e, 0x4e, 0xca, 0xd1, 0x05, 0xdc, 0xfe, 0xe0, 0xa6, 0x55, 0x8d, 0xc1, 0x81, 0x38, 0x06,
	0xfe, 0xb8, 0xca, 0x61, 0xf6, 0xaa, 0x11, 0xfb, 0x5a, 0xe4, 0x19, 0xe6, 0xeb, 0xd2, 0xf9, 0x18,
	0x24, 0xcd, 0xc6, 0x15, 0xb2, 0x71, 0xfe, 0x34, 0xf9, 0xd7, 0xeb, 0x47, 0xd0, 0xc9, 0x6a, 0x0d,
	0xd3, 0x97, 0x52, 0x72, 0xf2, 0x0e, 0xf3, 0x30, 0x98, 0xdc, 0xfd, 0xe6, 0xff, 0xe4, 0x4f, 0xef,
	0xcf, 0xd9, 0xa7, 0xbc, 0x74, 0xb2, 0xb9, 0xf2, 0x5b, 0x7d, 0xfa, 0x0a, 0x00, 0x00, 0xff, 0xff,
	0x4b, 0x23, 0xbb, 0x6f, 0xa6, 0x01, 0x00, 0x00,
}