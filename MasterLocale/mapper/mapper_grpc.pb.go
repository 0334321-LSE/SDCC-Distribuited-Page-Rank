// mapper.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: mapper.proto

package mapper

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
	Mapper_Map_FullMethodName     = "/Mapper/Map"
	Mapper_CleanUp_FullMethodName = "/Mapper/CleanUp"
)

// MapperClient is the client API for Mapper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MapperClient interface {
	Map(ctx context.Context, in *MapperInput, opts ...grpc.CallOption) (*MapperOutput, error)
	CleanUp(ctx context.Context, in *CleanUpInput, opts ...grpc.CallOption) (*CleanUpOutput, error)
}

type mapperClient struct {
	cc grpc.ClientConnInterface
}

func NewMapperClient(cc grpc.ClientConnInterface) MapperClient {
	return &mapperClient{cc}
}

func (c *mapperClient) Map(ctx context.Context, in *MapperInput, opts ...grpc.CallOption) (*MapperOutput, error) {
	out := new(MapperOutput)
	err := c.cc.Invoke(ctx, Mapper_Map_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mapperClient) CleanUp(ctx context.Context, in *CleanUpInput, opts ...grpc.CallOption) (*CleanUpOutput, error) {
	out := new(CleanUpOutput)
	err := c.cc.Invoke(ctx, Mapper_CleanUp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MapperServer is the server API for Mapper service.
// All implementations must embed UnimplementedMapperServer
// for forward compatibility
type MapperServer interface {
	Map(context.Context, *MapperInput) (*MapperOutput, error)
	CleanUp(context.Context, *CleanUpInput) (*CleanUpOutput, error)
	mustEmbedUnimplementedMapperServer()
}

// UnimplementedMapperServer must be embedded to have forward compatible implementations.
type UnimplementedMapperServer struct {
}

func (UnimplementedMapperServer) Map(context.Context, *MapperInput) (*MapperOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Map not implemented")
}
func (UnimplementedMapperServer) CleanUp(context.Context, *CleanUpInput) (*CleanUpOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CleanUp not implemented")
}
func (UnimplementedMapperServer) mustEmbedUnimplementedMapperServer() {}

// UnsafeMapperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MapperServer will
// result in compilation errors.
type UnsafeMapperServer interface {
	mustEmbedUnimplementedMapperServer()
}

func RegisterMapperServer(s grpc.ServiceRegistrar, srv MapperServer) {
	s.RegisterService(&Mapper_ServiceDesc, srv)
}

func _Mapper_Map_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MapperInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapperServer).Map(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Mapper_Map_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapperServer).Map(ctx, req.(*MapperInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mapper_CleanUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CleanUpInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapperServer).CleanUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Mapper_CleanUp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapperServer).CleanUp(ctx, req.(*CleanUpInput))
	}
	return interceptor(ctx, in, info, handler)
}

// Mapper_ServiceDesc is the grpc.ServiceDesc for Mapper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Mapper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Mapper",
	HandlerType: (*MapperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Map",
			Handler:    _Mapper_Map_Handler,
		},
		{
			MethodName: "CleanUp",
			Handler:    _Mapper_CleanUp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mapper.proto",
}

const (
	MapperHeartbeat_Ping_FullMethodName = "/MapperHeartbeat/Ping"
)

// MapperHeartbeatClient is the client API for MapperHeartbeat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MapperHeartbeatClient interface {
	Ping(ctx context.Context, in *MapperHeartbeatRequest, opts ...grpc.CallOption) (*MapperHeartbeatResponse, error)
}

type mapperHeartbeatClient struct {
	cc grpc.ClientConnInterface
}

func NewMapperHeartbeatClient(cc grpc.ClientConnInterface) MapperHeartbeatClient {
	return &mapperHeartbeatClient{cc}
}

func (c *mapperHeartbeatClient) Ping(ctx context.Context, in *MapperHeartbeatRequest, opts ...grpc.CallOption) (*MapperHeartbeatResponse, error) {
	out := new(MapperHeartbeatResponse)
	err := c.cc.Invoke(ctx, MapperHeartbeat_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MapperHeartbeatServer is the server API for MapperHeartbeat service.
// All implementations must embed UnimplementedMapperHeartbeatServer
// for forward compatibility
type MapperHeartbeatServer interface {
	Ping(context.Context, *MapperHeartbeatRequest) (*MapperHeartbeatResponse, error)
	mustEmbedUnimplementedMapperHeartbeatServer()
}

// UnimplementedMapperHeartbeatServer must be embedded to have forward compatible implementations.
type UnimplementedMapperHeartbeatServer struct {
}

func (UnimplementedMapperHeartbeatServer) Ping(context.Context, *MapperHeartbeatRequest) (*MapperHeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedMapperHeartbeatServer) mustEmbedUnimplementedMapperHeartbeatServer() {}

// UnsafeMapperHeartbeatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MapperHeartbeatServer will
// result in compilation errors.
type UnsafeMapperHeartbeatServer interface {
	mustEmbedUnimplementedMapperHeartbeatServer()
}

func RegisterMapperHeartbeatServer(s grpc.ServiceRegistrar, srv MapperHeartbeatServer) {
	s.RegisterService(&MapperHeartbeat_ServiceDesc, srv)
}

func _MapperHeartbeat_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MapperHeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapperHeartbeatServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MapperHeartbeat_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapperHeartbeatServer).Ping(ctx, req.(*MapperHeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MapperHeartbeat_ServiceDesc is the grpc.ServiceDesc for MapperHeartbeat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MapperHeartbeat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MapperHeartbeat",
	HandlerType: (*MapperHeartbeatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _MapperHeartbeat_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mapper.proto",
}