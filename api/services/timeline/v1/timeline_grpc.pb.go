// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: services/timeline/v1/timeline.proto

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

// TimelineServiceClient is the client API for TimelineService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TimelineServiceClient interface {
	PushTimeline(ctx context.Context, in *PushTimelineRequest, opts ...grpc.CallOption) (*PushTimelineResponse, error)
}

type timelineServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTimelineServiceClient(cc grpc.ClientConnInterface) TimelineServiceClient {
	return &timelineServiceClient{cc}
}

func (c *timelineServiceClient) PushTimeline(ctx context.Context, in *PushTimelineRequest, opts ...grpc.CallOption) (*PushTimelineResponse, error) {
	out := new(PushTimelineResponse)
	err := c.cc.Invoke(ctx, "/timeline.v1.TimelineService/PushTimeline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TimelineServiceServer is the server API for TimelineService service.
// All implementations must embed UnimplementedTimelineServiceServer
// for forward compatibility
type TimelineServiceServer interface {
	PushTimeline(context.Context, *PushTimelineRequest) (*PushTimelineResponse, error)
	mustEmbedUnimplementedTimelineServiceServer()
}

// UnimplementedTimelineServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTimelineServiceServer struct {
}

func (UnimplementedTimelineServiceServer) PushTimeline(context.Context, *PushTimelineRequest) (*PushTimelineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushTimeline not implemented")
}
func (UnimplementedTimelineServiceServer) mustEmbedUnimplementedTimelineServiceServer() {}

// UnsafeTimelineServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TimelineServiceServer will
// result in compilation errors.
type UnsafeTimelineServiceServer interface {
	mustEmbedUnimplementedTimelineServiceServer()
}

func RegisterTimelineServiceServer(s grpc.ServiceRegistrar, srv TimelineServiceServer) {
	s.RegisterService(&TimelineService_ServiceDesc, srv)
}

func _TimelineService_PushTimeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushTimelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimelineServiceServer).PushTimeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/timeline.v1.TimelineService/PushTimeline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimelineServiceServer).PushTimeline(ctx, req.(*PushTimelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TimelineService_ServiceDesc is the grpc.ServiceDesc for TimelineService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TimelineService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "timeline.v1.TimelineService",
	HandlerType: (*TimelineServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushTimeline",
			Handler:    _TimelineService_PushTimeline_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/timeline/v1/timeline.proto",
}
