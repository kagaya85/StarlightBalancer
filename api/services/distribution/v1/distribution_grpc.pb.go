// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: services/distribution/v1/distribution.proto

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

// DistributionServiceClient is the client API for DistributionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DistributionServiceClient interface {
	Distribute(ctx context.Context, in *DistributeRequest, opts ...grpc.CallOption) (*DistributeResponse, error)
}

type distributionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDistributionServiceClient(cc grpc.ClientConnInterface) DistributionServiceClient {
	return &distributionServiceClient{cc}
}

func (c *distributionServiceClient) Distribute(ctx context.Context, in *DistributeRequest, opts ...grpc.CallOption) (*DistributeResponse, error) {
	out := new(DistributeResponse)
	err := c.cc.Invoke(ctx, "/distribution.v1.DistributionService/Distribute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DistributionServiceServer is the server API for DistributionService service.
// All implementations must embed UnimplementedDistributionServiceServer
// for forward compatibility
type DistributionServiceServer interface {
	Distribute(context.Context, *DistributeRequest) (*DistributeResponse, error)
	mustEmbedUnimplementedDistributionServiceServer()
}

// UnimplementedDistributionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDistributionServiceServer struct {
}

func (UnimplementedDistributionServiceServer) Distribute(context.Context, *DistributeRequest) (*DistributeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Distribute not implemented")
}
func (UnimplementedDistributionServiceServer) mustEmbedUnimplementedDistributionServiceServer() {}

// UnsafeDistributionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DistributionServiceServer will
// result in compilation errors.
type UnsafeDistributionServiceServer interface {
	mustEmbedUnimplementedDistributionServiceServer()
}

func RegisterDistributionServiceServer(s grpc.ServiceRegistrar, srv DistributionServiceServer) {
	s.RegisterService(&DistributionService_ServiceDesc, srv)
}

func _DistributionService_Distribute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DistributeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DistributionServiceServer).Distribute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distribution.v1.DistributionService/Distribute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DistributionServiceServer).Distribute(ctx, req.(*DistributeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DistributionService_ServiceDesc is the grpc.ServiceDesc for DistributionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DistributionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "distribution.v1.DistributionService",
	HandlerType: (*DistributionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Distribute",
			Handler:    _DistributionService_Distribute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/distribution/v1/distribution.proto",
}
