// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: balancer/v1/update.proto

package v1

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

// WeightUpdaterClient is the client API for WeightUpdater service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WeightUpdaterClient interface {
	Update(ctx context.Context, in *UpdateRequeset, opts ...grpc.CallOption) (WeightUpdater_UpdateClient, error)
}

type weightUpdaterClient struct {
	cc grpc.ClientConnInterface
}

func NewWeightUpdaterClient(cc grpc.ClientConnInterface) WeightUpdaterClient {
	return &weightUpdaterClient{cc}
}

func (c *weightUpdaterClient) Update(ctx context.Context, in *UpdateRequeset, opts ...grpc.CallOption) (WeightUpdater_UpdateClient, error) {
	stream, err := c.cc.NewStream(ctx, &WeightUpdater_ServiceDesc.Streams[0], "/balancer.v1.WeightUpdater/Update", opts...)
	if err != nil {
		return nil, err
	}
	x := &weightUpdaterUpdateClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type WeightUpdater_UpdateClient interface {
	Recv() (*UpdateReply, error)
	grpc.ClientStream
}

type weightUpdaterUpdateClient struct {
	grpc.ClientStream
}

func (x *weightUpdaterUpdateClient) Recv() (*UpdateReply, error) {
	m := new(UpdateReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// WeightUpdaterServer is the server API for WeightUpdater service.
// All implementations must embed UnimplementedWeightUpdaterServer
// for forward compatibility
type WeightUpdaterServer interface {
	Update(*UpdateRequeset, WeightUpdater_UpdateServer) error
	mustEmbedUnimplementedWeightUpdaterServer()
}

// UnimplementedWeightUpdaterServer must be embedded to have forward compatible implementations.
type UnimplementedWeightUpdaterServer struct {
}

func (UnimplementedWeightUpdaterServer) Update(*UpdateRequeset, WeightUpdater_UpdateServer) error {
	return status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedWeightUpdaterServer) mustEmbedUnimplementedWeightUpdaterServer() {}

// UnsafeWeightUpdaterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WeightUpdaterServer will
// result in compilation errors.
type UnsafeWeightUpdaterServer interface {
	mustEmbedUnimplementedWeightUpdaterServer()
}

func RegisterWeightUpdaterServer(s grpc.ServiceRegistrar, srv WeightUpdaterServer) {
	s.RegisterService(&WeightUpdater_ServiceDesc, srv)
}

func _WeightUpdater_Update_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(UpdateRequeset)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WeightUpdaterServer).Update(m, &weightUpdaterUpdateServer{stream})
}

type WeightUpdater_UpdateServer interface {
	Send(*UpdateReply) error
	grpc.ServerStream
}

type weightUpdaterUpdateServer struct {
	grpc.ServerStream
}

func (x *weightUpdaterUpdateServer) Send(m *UpdateReply) error {
	return x.ServerStream.SendMsg(m)
}

// WeightUpdater_ServiceDesc is the grpc.ServiceDesc for WeightUpdater service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WeightUpdater_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "balancer.v1.WeightUpdater",
	HandlerType: (*WeightUpdaterServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Update",
			Handler:       _WeightUpdater_Update_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "balancer/v1/update.proto",
}