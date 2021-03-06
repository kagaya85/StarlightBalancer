// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: services/transcode/v1/transcode.proto

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

// TranscodeServiceClient is the client API for TranscodeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TranscodeServiceClient interface {
	Transcode(ctx context.Context, in *TranscodeRequest, opts ...grpc.CallOption) (*TranscodeResponse, error)
}

type transcodeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTranscodeServiceClient(cc grpc.ClientConnInterface) TranscodeServiceClient {
	return &transcodeServiceClient{cc}
}

func (c *transcodeServiceClient) Transcode(ctx context.Context, in *TranscodeRequest, opts ...grpc.CallOption) (*TranscodeResponse, error) {
	out := new(TranscodeResponse)
	err := c.cc.Invoke(ctx, "/transcode.v1.TranscodeService/Transcode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TranscodeServiceServer is the server API for TranscodeService service.
// All implementations must embed UnimplementedTranscodeServiceServer
// for forward compatibility
type TranscodeServiceServer interface {
	Transcode(context.Context, *TranscodeRequest) (*TranscodeResponse, error)
	mustEmbedUnimplementedTranscodeServiceServer()
}

// UnimplementedTranscodeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTranscodeServiceServer struct {
}

func (UnimplementedTranscodeServiceServer) Transcode(context.Context, *TranscodeRequest) (*TranscodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transcode not implemented")
}
func (UnimplementedTranscodeServiceServer) mustEmbedUnimplementedTranscodeServiceServer() {}

// UnsafeTranscodeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TranscodeServiceServer will
// result in compilation errors.
type UnsafeTranscodeServiceServer interface {
	mustEmbedUnimplementedTranscodeServiceServer()
}

func RegisterTranscodeServiceServer(s grpc.ServiceRegistrar, srv TranscodeServiceServer) {
	s.RegisterService(&TranscodeService_ServiceDesc, srv)
}

func _TranscodeService_Transcode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TranscodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TranscodeServiceServer).Transcode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transcode.v1.TranscodeService/Transcode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TranscodeServiceServer).Transcode(ctx, req.(*TranscodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TranscodeService_ServiceDesc is the grpc.ServiceDesc for TranscodeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TranscodeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "transcode.v1.TranscodeService",
	HandlerType: (*TranscodeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Transcode",
			Handler:    _TranscodeService_Transcode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/transcode/v1/transcode.proto",
}
