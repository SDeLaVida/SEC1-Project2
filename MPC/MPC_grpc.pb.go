// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: MPC/MPC.proto

package MPC

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MPC_SendChunk_FullMethodName  = "/MPC.MPC/sendChunk"
	MPC_SendResult_FullMethodName = "/MPC.MPC/sendResult"
)

// MPCClient is the client API for MPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MPCClient interface {
	SendChunk(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error)
	SendResult(ctx context.Context, in *Result, opts ...grpc.CallOption) (*Reply, error)
}

type mPCClient struct {
	cc grpc.ClientConnInterface
}

func NewMPCClient(cc grpc.ClientConnInterface) MPCClient {
	return &mPCClient{cc}
}

func (c *mPCClient) SendChunk(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, MPC_SendChunk_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mPCClient) SendResult(ctx context.Context, in *Result, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, MPC_SendResult_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MPCServer is the server API for MPC service.
// All implementations must embed UnimplementedMPCServer
// for forward compatibility
type MPCServer interface {
	SendChunk(context.Context, *Request) (*Reply, error)
	SendResult(context.Context, *Result) (*Reply, error)
	mustEmbedUnimplementedMPCServer()
}

// UnimplementedMPCServer must be embedded to have forward compatible implementations.
type UnimplementedMPCServer struct {
}

func (UnimplementedMPCServer) SendChunk(context.Context, *Request) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendChunk not implemented")
}
func (UnimplementedMPCServer) SendResult(context.Context, *Result) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendResult not implemented")
}
func (UnimplementedMPCServer) mustEmbedUnimplementedMPCServer() {}

// UnsafeMPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MPCServer will
// result in compilation errors.
type UnsafeMPCServer interface {
	mustEmbedUnimplementedMPCServer()
}

func RegisterMPCServer(s grpc.ServiceRegistrar, srv MPCServer) {
	s.RegisterService(&MPC_ServiceDesc, srv)
}

func _MPC_SendChunk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MPCServer).SendChunk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MPC_SendChunk_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MPCServer).SendChunk(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _MPC_SendResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Result)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MPCServer).SendResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MPC_SendResult_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MPCServer).SendResult(ctx, req.(*Result))
	}
	return interceptor(ctx, in, info, handler)
}

// MPC_ServiceDesc is the grpc.ServiceDesc for MPC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MPC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MPC.MPC",
	HandlerType: (*MPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "sendChunk",
			Handler:    _MPC_SendChunk_Handler,
		},
		{
			MethodName: "sendResult",
			Handler:    _MPC_SendResult_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "MPC/MPC.proto",
}
