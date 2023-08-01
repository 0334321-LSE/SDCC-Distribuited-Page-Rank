//reducer-proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: reducer.proto

package proto

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
	Reducer_Reduce_FullMethodName = "/Reducer/Reduce"
)

// ReducerClient is the client API for Reducer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReducerClient interface {
	Reduce(ctx context.Context, in *ReducerInput, opts ...grpc.CallOption) (*ReducerOutput, error)
}

type reducerClient struct {
	cc grpc.ClientConnInterface
}

func NewReducerClient(cc grpc.ClientConnInterface) ReducerClient {
	return &reducerClient{cc}
}

func (c *reducerClient) Reduce(ctx context.Context, in *ReducerInput, opts ...grpc.CallOption) (*ReducerOutput, error) {
	out := new(ReducerOutput)
	err := c.cc.Invoke(ctx, Reducer_Reduce_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReducerServer is the server API for Reducer service.
// All implementations must embed UnimplementedReducerServer
// for forward compatibility
type ReducerServer interface {
	Reduce(context.Context, *ReducerInput) (*ReducerOutput, error)
	mustEmbedUnimplementedReducerServer()
}

// UnimplementedReducerServer must be embedded to have forward compatible implementations.
type UnimplementedReducerServer struct {
}

func (UnimplementedReducerServer) Reduce(context.Context, *ReducerInput) (*ReducerOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reduce not implemented")
}
func (UnimplementedReducerServer) mustEmbedUnimplementedReducerServer() {}

// UnsafeReducerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReducerServer will
// result in compilation errors.
type UnsafeReducerServer interface {
	mustEmbedUnimplementedReducerServer()
}

func RegisterReducerServer(s grpc.ServiceRegistrar, srv ReducerServer) {
	s.RegisterService(&Reducer_ServiceDesc, srv)
}

func _Reducer_Reduce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReducerInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReducerServer).Reduce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Reducer_Reduce_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReducerServer).Reduce(ctx, req.(*ReducerInput))
	}
	return interceptor(ctx, in, info, handler)
}

// Reducer_ServiceDesc is the grpc.ServiceDesc for Reducer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Reducer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Reducer",
	HandlerType: (*ReducerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Reduce",
			Handler:    _Reducer_Reduce_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "reducer.proto",
}
