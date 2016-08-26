// Code generated by protoc-gen-go.
// source: github.com/luci/luci-go/dm/api/service/v1/service.proto
// DO NOT EDIT!

package dm

import prpc "github.com/luci/luci-go/grpc/prpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf2 "github.com/luci/luci-go/common/proto/google"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Deps service

type DepsClient interface {
	// allows you to add additional data to the current dependency graph.
	EnsureGraphData(ctx context.Context, in *EnsureGraphDataReq, opts ...grpc.CallOption) (*EnsureGraphDataRsp, error)
	// is called by Execution clients to activate themselves with DM.
	ActivateExecution(ctx context.Context, in *ActivateExecutionReq, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
	// is called by Execution clients to indicate that an Attempt is finished.
	FinishAttempt(ctx context.Context, in *FinishAttemptReq, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
	// runs queries, and walks along the dependency graph from the query results.
	WalkGraph(ctx context.Context, in *WalkGraphReq, opts ...grpc.CallOption) (*GraphData, error)
}
type depsPRPCClient struct {
	client *prpc.Client
}

func NewDepsPRPCClient(client *prpc.Client) DepsClient {
	return &depsPRPCClient{client}
}

func (c *depsPRPCClient) EnsureGraphData(ctx context.Context, in *EnsureGraphDataReq, opts ...grpc.CallOption) (*EnsureGraphDataRsp, error) {
	out := new(EnsureGraphDataRsp)
	err := c.client.Call(ctx, "dm.Deps", "EnsureGraphData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsPRPCClient) ActivateExecution(ctx context.Context, in *ActivateExecutionReq, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := c.client.Call(ctx, "dm.Deps", "ActivateExecution", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsPRPCClient) FinishAttempt(ctx context.Context, in *FinishAttemptReq, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := c.client.Call(ctx, "dm.Deps", "FinishAttempt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsPRPCClient) WalkGraph(ctx context.Context, in *WalkGraphReq, opts ...grpc.CallOption) (*GraphData, error) {
	out := new(GraphData)
	err := c.client.Call(ctx, "dm.Deps", "WalkGraph", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type depsClient struct {
	cc *grpc.ClientConn
}

func NewDepsClient(cc *grpc.ClientConn) DepsClient {
	return &depsClient{cc}
}

func (c *depsClient) EnsureGraphData(ctx context.Context, in *EnsureGraphDataReq, opts ...grpc.CallOption) (*EnsureGraphDataRsp, error) {
	out := new(EnsureGraphDataRsp)
	err := grpc.Invoke(ctx, "/dm.Deps/EnsureGraphData", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsClient) ActivateExecution(ctx context.Context, in *ActivateExecutionReq, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/dm.Deps/ActivateExecution", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsClient) FinishAttempt(ctx context.Context, in *FinishAttemptReq, opts ...grpc.CallOption) (*google_protobuf2.Empty, error) {
	out := new(google_protobuf2.Empty)
	err := grpc.Invoke(ctx, "/dm.Deps/FinishAttempt", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depsClient) WalkGraph(ctx context.Context, in *WalkGraphReq, opts ...grpc.CallOption) (*GraphData, error) {
	out := new(GraphData)
	err := grpc.Invoke(ctx, "/dm.Deps/WalkGraph", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Deps service

type DepsServer interface {
	// allows you to add additional data to the current dependency graph.
	EnsureGraphData(context.Context, *EnsureGraphDataReq) (*EnsureGraphDataRsp, error)
	// is called by Execution clients to activate themselves with DM.
	ActivateExecution(context.Context, *ActivateExecutionReq) (*google_protobuf2.Empty, error)
	// is called by Execution clients to indicate that an Attempt is finished.
	FinishAttempt(context.Context, *FinishAttemptReq) (*google_protobuf2.Empty, error)
	// runs queries, and walks along the dependency graph from the query results.
	WalkGraph(context.Context, *WalkGraphReq) (*GraphData, error)
}

func RegisterDepsServer(s prpc.Registrar, srv DepsServer) {
	s.RegisterService(&_Deps_serviceDesc, srv)
}

func _Deps_EnsureGraphData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnsureGraphDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepsServer).EnsureGraphData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dm.Deps/EnsureGraphData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepsServer).EnsureGraphData(ctx, req.(*EnsureGraphDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deps_ActivateExecution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivateExecutionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepsServer).ActivateExecution(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dm.Deps/ActivateExecution",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepsServer).ActivateExecution(ctx, req.(*ActivateExecutionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deps_FinishAttempt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FinishAttemptReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepsServer).FinishAttempt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dm.Deps/FinishAttempt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepsServer).FinishAttempt(ctx, req.(*FinishAttemptReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deps_WalkGraph_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WalkGraphReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepsServer).WalkGraph(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dm.Deps/WalkGraph",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepsServer).WalkGraph(ctx, req.(*WalkGraphReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Deps_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dm.Deps",
	HandlerType: (*DepsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EnsureGraphData",
			Handler:    _Deps_EnsureGraphData_Handler,
		},
		{
			MethodName: "ActivateExecution",
			Handler:    _Deps_ActivateExecution_Handler,
		},
		{
			MethodName: "FinishAttempt",
			Handler:    _Deps_FinishAttempt_Handler,
		},
		{
			MethodName: "WalkGraph",
			Handler:    _Deps_WalkGraph_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor5,
}

func init() {
	proto.RegisterFile("github.com/luci/luci-go/dm/api/service/v1/service.proto", fileDescriptor5)
}

var fileDescriptor5 = []byte{
	// 279 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x91, 0x41, 0x4b, 0xc4, 0x30,
	0x10, 0x85, 0x51, 0x44, 0x30, 0xb0, 0xa8, 0x41, 0x44, 0xd6, 0xbf, 0xa0, 0x09, 0xea, 0x41, 0x10,
	0x14, 0x2a, 0x5b, 0xbd, 0x7b, 0xf1, 0x58, 0xd2, 0x74, 0x36, 0x0d, 0x36, 0x4d, 0x6c, 0x92, 0xaa,
	0x3f, 0x5e, 0x70, 0x93, 0xb4, 0x0b, 0xba, 0x22, 0xf5, 0x12, 0x92, 0x97, 0x79, 0x8f, 0xf9, 0x66,
	0xd0, 0xb5, 0x90, 0xae, 0xf6, 0x25, 0xe1, 0x5a, 0xd1, 0xc6, 0x73, 0x19, 0x8f, 0x73, 0xa1, 0x69,
	0xa5, 0x28, 0x33, 0x92, 0x5a, 0xe8, 0x7a, 0xc9, 0x81, 0xf6, 0x17, 0xe3, 0x95, 0x98, 0x4e, 0x3b,
	0x8d, 0xb7, 0x2b, 0x35, 0x3f, 0x15, 0x5a, 0x8b, 0x06, 0x68, 0x54, 0x4a, 0xbf, 0xa4, 0xa0, 0x8c,
	0xfb, 0x48, 0x05, 0xf3, 0x9b, 0xe9, 0xc9, 0xa2, 0x63, 0xa6, 0x2e, 0x2a, 0xe6, 0xd8, 0xe0, 0xcd,
	0xa6, 0x7b, 0xa1, 0xb5, 0xbe, 0x83, 0x62, 0x23, 0xe2, 0x7e, 0x7a, 0x04, 0xe3, 0x4e, 0xf6, 0xcc,
	0x41, 0x01, 0xef, 0xc0, 0xbd, 0x93, 0xba, 0x1d, 0x32, 0xee, 0xa6, 0x67, 0x2c, 0x65, 0x2b, 0x6d,
	0x5d, 0x30, 0xe7, 0xc2, 0x0c, 0xfe, 0x3f, 0x82, 0x37, 0xd6, 0xbc, 0x24, 0x88, 0xe4, 0xbd, 0xfc,
	0xdc, 0x42, 0x3b, 0x0b, 0x30, 0x16, 0x67, 0x68, 0x3f, 0x8f, 0x8c, 0x8f, 0xe1, 0x77, 0xb1, 0x22,
	0xc4, 0xc7, 0xa4, 0x52, 0xe4, 0x87, 0xf8, 0x04, 0xaf, 0xf3, 0x5f, 0x75, 0x6b, 0x70, 0x8e, 0x0e,
	0xb3, 0x81, 0x31, 0x1f, 0x11, 0xf1, 0x49, 0x28, 0xde, 0x90, 0x53, 0x4c, 0xda, 0x2b, 0x19, 0xf7,
	0x4a, 0xf2, 0xb0, 0x57, 0x7c, 0x8b, 0x66, 0x0f, 0x11, 0x33, 0x4b, 0x94, 0xf8, 0x28, 0x44, 0x7c,
	0x93, 0xfe, 0xb2, 0x9f, 0xa1, 0xbd, 0xe7, 0x15, 0x65, 0xec, 0x0c, 0x1f, 0x04, 0xeb, 0xfa, 0x19,
	0x6c, 0xb3, 0xa0, 0xac, 0xdb, 0x2e, 0x77, 0xa3, 0xfb, 0xea, 0x2b, 0x00, 0x00, 0xff, 0xff, 0xef,
	0xb9, 0x47, 0x29, 0xa1, 0x02, 0x00, 0x00,
}